package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/entity"
	"GolangBookingApp/internal/repositories"

	"github.com/google/uuid"
)

type createRoom struct {
	roomRepository repositories.Room
	urepo          repositories.User
}

func NewCreateRoom(roomRepository repositories.Room, urepo repositories.User) *createRoom {
	return &createRoom{
		roomRepository: roomRepository,
		urepo:          urepo,
	}
}

func (cr createRoom) Serve(data *appctx.Data) appctx.Response {
	a := &entity.Room{}

	err := data.Cast(&a)
	fmt.Println(a)
	if err != nil {
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Room").WithState("AddRoomFailed").WithMessage("Adding Room Failed").WithError(err.Error())
	}

	user := entity.User{}
	user, _ = cr.urepo.VerifySubmit(data.Request.Context(), a.CreatedBy)

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

	if user.Username != a.CreatedBy {
		err = fmt.Errorf("user not exist")
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("UserNotExist").WithMessage("User Not Exist").WithError(err.Error())
	}

	user, err = cr.roomRepository.VerifyRoom(data.Request.Context(), a.CreatedBy)

	if user.Username != a.CreatedBy {
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Room").WithState("AddRoomFailed").WithMessage("Adding Room Failed").WithError(err.Error())
	}

	room := entity.Room{}

	room, _ = cr.roomRepository.VerifyRoomEntity(data.Request.Context(), a.RoomName)

	if room.RoomName == a.RoomName {
		err = fmt.Errorf("room already exist")
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Room").WithState("AddRoomFailed").WithMessage("Adding Room Failed").WithError(err.Error())
	}

	id, err := cr.roomRepository.Create(data.Request.Context(), a)

	if err != nil {
		fmt.Println(err)
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Room").WithState("AddRoomFailed").WithMessage("Adding Room Failed").WithError(err.Error())
	}

	type resp struct {
		ID uuid.UUID `json:"id"`
	}

	var res resp
	res.ID = *id

	return *appctx.NewResponse().WithCode(http.StatusCreated).WithStatus("SUCCESS").WithEntity("Room").WithState("AddingRoomSuccess").WithMessage("Adding Room Success").WithData(res)
}
