package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/entity"
	"GolangBookingApp/internal/repositories"

	"github.com/gorilla/mux"
)

type deleteRoom struct {
	repo  repositories.Room
	urepo repositories.User
}

func NewDeleteRoom(repo repositories.Room, urepo repositories.User) *deleteRoom {
	return &deleteRoom{
		repo:  repo,
		urepo: urepo}
}

func (u deleteRoom) Serve(data *appctx.Data) appctx.Response {
	// a := &entity.Room{}
	// err := data.Cast(&a)
	// fmt.Println(a)
	// if err != nil {
	// 	return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Room").WithState("AddRoomFailed").WithMessage("Adding Room Failed").WithError(err.Error())
	// }

	id := mux.Vars(data.Request)["roomID"]
	uname := mux.Vars(data.Request)["username"]

	user := entity.User{}
	user, err := u.urepo.VerifySubmit(data.Request.Context(), uname)

	fmt.Println(user)

	const ruser int8 = 0
	// const admin int8 = 1
	// const sup int8 = 2
	role := user.Role

	if role == ruser {
		err = fmt.Errorf("access not permitted")
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("AccessNotPermitted").WithMessage("Access Not Permited").WithError(err.Error())
	}

	if !user.Logged {
		err = fmt.Errorf("user not logged in")
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("UserNotLoggedIn").WithMessage("User Not Logged In").WithError(err.Error())
	}

	delete := u.repo.DeleteRoom(data.Request.Context(), id)

	if delete != nil {
		err := fmt.Errorf("deleting room: %w", err)
		fmt.Println(err)

		return *appctx.NewResponse().WithCode(http.StatusNotFound).WithStatus("FAIL").WithEntity("Room").WithState("DeleteRoomFail").WithMessage("Delete Room Fail").WithError(err.Error())
	}

	return *appctx.NewResponse().WithCode(http.StatusOK).WithStatus("SUCCESS").WithEntity("Room").WithState("DeleteRoomSuccess").WithMessage("Delete Room Success")
}
