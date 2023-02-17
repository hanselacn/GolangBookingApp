package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/entity"
	"GolangBookingApp/internal/repositories"

	"github.com/google/uuid"
)

type createBooking struct {
	bookingRepository repositories.Booking
	userRepository    repositories.User
}

func NewCreateBooking(bookingRepository repositories.Booking, userRepository repositories.User) *createBooking {
	return &createBooking{
		bookingRepository: bookingRepository,
		userRepository:    userRepository,
	}
}

func (cb createBooking) Serve(data *appctx.Data) appctx.Response {
	a := &entity.Booking{}

	err := data.Cast(&a)

	if err != nil {
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("AddBookingFailed").WithMessage("Adding Book Failed").WithError(err.Error())
	}

	user := entity.User{}

	user, _ = cb.userRepository.VerifySubmit(data.Request.Context(), a.BookedBy)

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

	if user.Username != a.BookedBy {
		err = fmt.Errorf("user not exist")
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("UserNotExist").WithMessage("User Not Exist").WithError(err.Error())
	}

	day := entity.Day{}

	day, err = cb.bookingRepository.VerifyBookingDay(data.Request.Context(), a.BookedDay)

	if day.DayName != a.BookedDay {
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("AddBookingFailed").WithMessage("Adding Book Failed").WithError(err.Error())
	}

	b := entity.Booking{}

	b, _ = cb.bookingRepository.VerifyBookingEntity(data.Request.Context(), a)

	fmt.Println("a : ", a)
	fmt.Println("b : ", b)
	fmt.Println("a : ", a.BookedRoom)
	fmt.Println("b : ", b.BookedRoom)

	if b.BookedRoom == a.BookedRoom {
		err = fmt.Errorf("session already booked")
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("AddBookingFailed").WithMessage("Adding Book Failed").WithError(err.Error())
	}

	id, err := cb.bookingRepository.Create(data.Request.Context(), a)

	if err != nil {
		fmt.Println(err)
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Booking").WithState("AddBookingFailed").WithMessage("Adding Book Failed").WithError(err.Error())
	}

	type resp struct {
		ID uuid.UUID `json:"id"`
	}

	var res resp
	res.ID = *id

	return *appctx.NewResponse().WithCode(http.StatusCreated).WithStatus("SUCCESS").WithEntity("Booking").WithState("AddBookingSuccess").WithMessage("Adding Booking Success").WithData(res)
}
