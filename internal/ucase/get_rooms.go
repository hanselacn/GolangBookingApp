package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/repositories"
)

type getRooms struct {
	repo repositories.Room
}

func NewGetRooms(repo repositories.Room) *getRooms {
	return &getRooms{repo: repo}
}

func (u *getRooms) Serve(data *appctx.Data) appctx.Response {

	rooms, err := u.repo.GetAllRoom(data.Request.Context())

	if err != nil {
		err := fmt.Errorf("getting rooms: %w", err)
		fmt.Println(err)

		return *appctx.NewResponse().WithCode(http.StatusInternalServerError).WithStatus("ERROR").WithMessage("Internal Server Error")
	}
	return *appctx.NewResponse().WithCode(http.StatusOK).WithStatus("SUCCESS").WithEntity("Rooms").WithState("GettingRoomsSuccess").WithMessage("Getting Rooms Success").WithData(rooms)
}
