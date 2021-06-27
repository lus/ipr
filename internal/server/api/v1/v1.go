package v1

import (
	routing "github.com/fasthttp/router"
	"github.com/lus/ipr/internal/shared"
)

// App represents the V1 API app
type App struct {
	AuthToken         string
	MachineRepository shared.MachineRepository
}

// Route routes the v1 API endpoints
func (app *App) Route(group *routing.Group) {
	group.GET("/machines", app.authorized(app.endpointGetMachines))
	group.POST("/machines", app.authorized(app.endpointCreateMachine))
}
