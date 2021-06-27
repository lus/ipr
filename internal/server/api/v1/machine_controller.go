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
	hash, err := tkn.Hash()
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
	copy.Token = tkn.Raw()
	if err := app.json(ctx, fasthttp.StatusCreated, copy); err != nil {
		app.error(ctx, -1, err)
	}
}
