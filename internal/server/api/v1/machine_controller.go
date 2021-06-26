package v1

import "github.com/valyala/fasthttp"

// endpointGetMachines handles the 'GET /api/v1/machines' endpoint
func (app *App) endpointGetMachines(ctx *fasthttp.RequestCtx) {
	// Look up all stored machines
	machines, err := app.MachineRepository.All()
	if err != nil {
		app.error(ctx, -1, err)
		return
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
