// Package router
package server

import (
	"context"
	"net/http"

	"GolangBookingApp/internal/appctx"
	ucase "GolangBookingApp/internal/ucase/contract"
)

// httpHandlerFunc abstraction for http handler
type httpHandlerFunc func(request *http.Request, svc ucase.UseCase, conf *appctx.Config) appctx.Response

// Server contract
type Server interface {
	Run(context.Context) error
	Done()
	Config() *appctx.Config
}
