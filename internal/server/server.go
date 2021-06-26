package server

import (
	routing "github.com/fasthttp/router"
	v1 "github.com/lus/ipr/internal/server/api/v1"
	"github.com/lus/ipr/internal/shared"
	"github.com/valyala/fasthttp"
)

// Settings represents the settings for the web server
type Settings struct {
	Address   string
	AuthToken string
}

// Settings represents the repositories for the web server
type Repositories struct {
	MachineRepository shared.MachineRepository
}

// RunBlocking runs the web server
func RunBlocking(settings *Settings, repositories *Repositories) error {
	router := routing.New()

	// Route the v1 API endpoints
	(&v1.App{
		AuthToken:         settings.AuthToken,
		MachineRepository: repositories.MachineRepository,
	}).Route(router.Group("/api/v1"))

	return fasthttp.ListenAndServe(settings.Address, router.Handler)
}
