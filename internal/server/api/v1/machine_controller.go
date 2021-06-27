package v1

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/lus/ipr/internal/shared"
	"github.com/lus/ipr/internal/token"
	"github.com/valyala/fasthttp"
)

// middlewareMachineAuthorization handles machine token OR auth token authorization
func (app *App) middlewareMachineAuthorization(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := string(ctx.Request.Header.Peek("Authorization"))
		authHeaderSplit := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderSplit) != 2 || authHeaderSplit[0] != "Bearer" {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("unauthorized")
			return
		}

		if authHeaderSplit[1] == app.AuthToken {
			handler(ctx)
			return
		}

		machine := ctx.UserValue("_machine").(*shared.Machine)
		valid, _ := token.Check(machine.Token, authHeaderSplit[1])
		if !valid {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("unauthorized")
			return
		}

		handler(ctx)
	}
}

// middlewareInjectMachine handles machine injection based on the 'name' request parameter
func (app *App) middlewareInjectMachine(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		machineName := ctx.UserValue("name").(string)

		machine, err := app.MachineRepository.Lookup(machineName)
		if err != nil {
			app.error(ctx, -1, err)
			return
		}
		if machine == nil {
			app.error(ctx, fasthttp.StatusNotFound, errors.New("machine not found"))
			return
		}

		ctx.SetUserValue("_machine", machine)
		handler(ctx)
	}
}

// endpointReportMachineAddress handles the 'POST /api/v1/machines/{name}/report' endpoint
func (app *App) endpointReportMachineAddress(ctx *fasthttp.RequestCtx) {
	machine := ctx.UserValue("_machine").(*shared.Machine)
	address := string(ctx.Request.Body())

	machine.Address = address
	machine.Updated = time.Now().Unix()

	if err := app.MachineRepository.Upsert(machine); err != nil {
		app.error(ctx, -1, err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// endpointGetMachines handles the 'GET /api/v1/machines' endpoint
func (app *App) endpointGetMachines(ctx *fasthttp.RequestCtx) {
	// Look up all stored machines
	machines, err := app.MachineRepository.All()
	if err != nil {
		app.error(ctx, -1, err)
		return
	}
	if machines == nil {
		machines = []*shared.Machine{}
	}

	// Remove the token field of every machine
	for _, machine := range machines {
		machine.Token = ""
	}

	// Respond with a JSON representation of the machine list
	if err := app.json(ctx, fasthttp.StatusOK, machines); err != nil {
		app.error(ctx, -1, err)
	}
}

type endpointCreateMachinePayload struct {
	Name string `json:"name"`
}

// endpointGetMachines handles the 'POST /api/v1/machines' endpoint
func (app *App) endpointCreateMachine(ctx *fasthttp.RequestCtx) {
	// Read, unmarshal and validate the request payload
	payload := new(endpointCreateMachinePayload)
	if err := json.Unmarshal(ctx.Request.Body(), payload); err != nil {
		app.error(ctx, fasthttp.StatusBadRequest, err)
		return
	}
	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" {
		app.error(ctx, fasthttp.StatusBadRequest, errors.New("name must not be empty"))
		return
	}

	// Check if a machine with that name already exists
	existing, err := app.MachineRepository.Lookup(payload.Name)
	if err != nil {
		app.error(ctx, -1, err)
		return
	}
	if existing != nil {
		app.error(ctx, fasthttp.StatusConflict, errors.New("machine name taken"))
		return
	}

	// Generate a new machine token
	tkn := token.Generate()
	hash, err := token.Hash(tkn)
	if err != nil {
		app.error(ctx, -1, err)
		return
	}

	// Create and store the machine
	machine := &shared.Machine{
		Name:    payload.Name,
		Token:   hash,
		Address: "<initial>",
		Updated: time.Now().Unix(),
	}
	if err := app.MachineRepository.Upsert(machine); err != nil {
		app.error(ctx, -1, err)
		return
	}

	// Respond with a JSON representation of the created machine including the raw token
	copy := *machine
	copy.Token = tkn
	if err := app.json(ctx, fasthttp.StatusCreated, copy); err != nil {
		app.error(ctx, -1, err)
	}
}

// endpointDeleteMachine handles the 'DELETE /api/v1/machines/{name}' endpoint
func (app *App) endpointDeleteMachine(ctx *fasthttp.RequestCtx) {
	machine := ctx.UserValue("_machine").(*shared.Machine)

	if err := app.MachineRepository.Delete(machine.Name); err != nil {
		app.error(ctx, -1, err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}
