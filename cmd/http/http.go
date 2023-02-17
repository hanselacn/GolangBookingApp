package http

import (
	"context"

	"GolangBookingApp/internal/consts"
	"GolangBookingApp/internal/server"
	"GolangBookingApp/pkg/logger"
)

// Start function handler starting http listener
func Start(ctx context.Context) {

	serve := server.NewHTTPServer()
	defer serve.Done()
	logger.Info(logger.MessageFormat("starting go_graphql services... %d", serve.Config().App.Port), logger.EventName(consts.LogEventNameServiceStarting))

	if err := serve.Run(ctx); err != nil {
		logger.Warn(logger.MessageFormat("service stopped, err:%s", err.Error()), logger.EventName(consts.LogEventNameServiceStarting))
	}

	return
}
