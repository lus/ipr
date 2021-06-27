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
	group.POST("/machines/{name}/report", app.middlewareInjectMachine(app.middlewareMachineAuthorization(app.endpointReportMachineAddress)))
	group.GET("/machines", app.authorized(app.endpointGetMachines))
	group.POST("/machines", app.authorized(app.endpointCreateMachine))
	group.DELETE("/machines/{name}", app.authorized(app.middlewareInjectMachine(app.endpointDeleteMachine)))
}
