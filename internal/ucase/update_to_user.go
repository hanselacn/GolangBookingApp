package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/entity"
	"GolangBookingApp/internal/repositories"

	"github.com/gorilla/mux"
)

type updateUser struct {
	repo repositories.User
}

func NewUpdateUser(repo repositories.User) *updateUser {
	return &updateUser{repo: repo}
}

func (u *updateUser) Serve(data *appctx.Data) appctx.Response {
	id := mux.Vars(data.Request)["userID"]
	uname := mux.Vars(data.Request)["username"]

	user := entity.User{}
	user, err := u.repo.VerifySubmit(data.Request.Context(), uname)

	fmt.Println(user)

	// const ruser int8 = 0
	// const admin int8 = 1
	const sup int8 = 2
	role := user.Role

	if role != sup {
		err = fmt.Errorf("access not permitted")
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("AccessNotPermitted").WithMessage("Access Not Permited").WithError(err.Error())
	}

	if !user.Logged {
		err = fmt.Errorf("user not logged in")
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("UserNotLoggedIn").WithMessage("User Not Logged In").WithError(err.Error())
	}

	granted := u.repo.DemoteToUser(data.Request.Context(), id)

	if granted != nil {
		err := fmt.Errorf("demote to user: %w", err)
		fmt.Println(err)
		return *appctx.NewResponse().WithCode(http.StatusNotFound).WithStatus("FAIL").WithEntity("User").WithState("UpdateUserFailed").WithMessage("Update User Fail").WithError(err.Error())
	}

	return *appctx.NewResponse().WithCode(http.StatusOK).WithStatus("SUCCESS").WithEntity("User").WithState("updateUserSuccess").WithMessage("Update User Success")
}
