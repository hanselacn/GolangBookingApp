package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/entity"
	"GolangBookingApp/internal/repositories"

	"github.com/google/uuid"
)

type createDay struct {
	dayRepository repositories.Day
	urepo         repositories.User
}

func NewCreateDay(dayRepository repositories.Day, urepo repositories.User) *createDay {
	return &createDay{
		dayRepository: dayRepository,
		urepo:         urepo,
	}
}

func (cd createDay) Serve(data *appctx.Data) appctx.Response {
	a := &entity.Day{}

	err := data.Cast(&a)

	if err != nil {
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Room").WithState("AddRoomFailed").WithMessage("Adding Room Failed").WithError(err.Error())
	}

	user := entity.User{}

	user, _ = cd.urepo.VerifySubmit(data.Request.Context(), a.CreatedBy)

	fmt.Println(user)

	// const ruser int8 = 0
	// const admin int8 = 1
	// const sup int8 = 2
	// role := user.Role

	// if role != ruser{
	// 	err = fmt.Errorf("access not permitted")
	// 	return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("AccessNotPermitted").WithMessage("Access Not Permited").WithError(err.Error())
	// }

	if !user.Logged {
		err = fmt.Errorf("user not logged in")
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("UserNotLoggedIn").WithMessage("User Not Logged In").WithError(err.Error())
	}

	if user.Username != a.CreatedBy {
		err = fmt.Errorf("user not exist")
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("UserNotExist").WithMessage("User Not Exist").WithError(err.Error())
	}

	user, err = cd.dayRepository.VerifyDay(data.Request.Context(), a.CreatedBy)

	if user.Username != a.CreatedBy {
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Room").WithState("AddRoomFailed").WithMessage("Adding Room Failed").WithError(err.Error())
	}

	id, err := cd.dayRepository.NewCreateDay(data.Request.Context(), a)

	if err != nil {
		fmt.Println(err)
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Room").WithState("AddRoomFailed").WithMessage("Adding Room Failed").WithError(err.Error())
	}

	type resp struct {
		ID uuid.UUID `json:"id"`
	}

	var res resp
	res.ID = *id

	return *appctx.NewResponse().WithCode(http.StatusCreated).WithStatus("SUCCESS").WithEntity("Users").WithState("AddingUserSuccess").WithMessage("Adding User Success").WithData(res)
}
