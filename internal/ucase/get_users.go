package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/repositories"
)

type getUsers struct {
	repo repositories.User
}

func NewGetUsers(repo repositories.User) *getUsers {
	return &getUsers{repo: repo}
}

func (u *getUsers) Serve(data *appctx.Data) appctx.Response {

	users, err := u.repo.GetAllUser(data.Request.Context())

	if err != nil {
		err := fmt.Errorf("getting rooms: %w", err)
		fmt.Println(err)

		return *appctx.NewResponse().WithCode(http.StatusInternalServerError).WithStatus("ERROR").WithMessage("Internal Server Error")
	}
	return *appctx.NewResponse().WithCode(http.StatusOK).WithStatus("SUCCESS").WithEntity("Users").WithState("GetAllUserSuccess").WithMessage("GetAllUserSuccess").WithData(users)
}
