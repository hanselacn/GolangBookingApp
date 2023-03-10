// Package example
package example

import (
	"github.com/gorilla/mux"
	"github.com/spf13/cast"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/consts"
	"GolangBookingApp/internal/repositories"
	"GolangBookingApp/internal/ucase/contract"

	"GolangBookingApp/pkg/logger"
)

type exampleDelete struct {
	repo repositories.Example
}

func NewExampleDelete(repo repositories.Example) contract.UseCase {
	return &exampleDelete{repo: repo}
}

// Serve partner list data
func (u *exampleDelete) Serve(data *appctx.Data) appctx.Response {

	id := mux.Vars(data.Request)["id"]

	err := u.repo.Delete(data.Request.Context(), cast.ToUint64(id))

	if err != nil {
		logger.Error(logger.MessageFormat("[example-delete] %v", err))

		return *appctx.NewResponse().WithCode(consts.CodeInternalServerError)
	}

	return *appctx.NewResponse().WithCode(consts.CodeSuccess)
}
