package v1

import (
	"encoding/json"
	"strings"

	"github.com/valyala/fasthttp"
)

// authorized wraps a fasthttp request handler so that it requires auth token authentication
func (app *App) authorized(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := string(ctx.Request.Header.Peek("Authorization"))
		authHeaderSplit := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderSplit) != 2 || authHeaderSplit[0] != "Bearer" || authHeaderSplit[1] != app.AuthToken {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("unauthorized")
			return
		}

		handler(ctx)
	}
}

// error sends a basic error response
func (app *App) error(ctx *fasthttp.RequestCtx, code int, err error) {
	ctx.SetStatusCode(code)
	if code < 0 {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
	ctx.SetBodyString(err.Error())
}

// json sends a basic json response
func (app *App) json(ctx *fasthttp.RequestCtx, code int, data interface{}) error {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx.SetStatusCode(code)
	ctx.SetBody(marshalled)

	return nil
}
