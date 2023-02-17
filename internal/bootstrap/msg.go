// Package bootstrap
package bootstrap

import (
	"GolangBookingApp/internal/consts"
	"GolangBookingApp/pkg/logger"
	"GolangBookingApp/pkg/msgx"
)

func RegistryMessage() {
	err := msgx.Setup("msg.yaml", consts.ConfigPath)
	if err != nil {
		logger.Fatal(logger.MessageFormat("file message multi language load error %s", err.Error()))
	}

}
