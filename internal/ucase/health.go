package ucase

import (
	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/consts"
	"GolangBookingApp/internal/ucase/contract"
)

type healthCheck struct {
}

func NewHealthCheck() contract.UseCase {
	return &healthCheck{}
}

func (u *healthCheck) Serve(*appctx.Data) appctx.Response {
	return *appctx.NewResponse().WithCode(consts.CodeSuccess).WithMessage("ok")
}
