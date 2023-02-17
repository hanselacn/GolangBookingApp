package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/repositories"
)

type deleteAllBooking struct {
	repo repositories.Booking
}

func NewDeleteAllBooking(repo repositories.Booking) *deleteAllBooking {
	return &deleteAllBooking{repo: repo}
}

func (u *deleteAllBooking) Serve(data *appctx.Data) appctx.Response {
	err := u.repo.DeleteAll(data.Request.Context())

	if err != nil {
		err := fmt.Errorf("deleting bookings: %w", err)
		fmt.Println(err)

		return *appctx.NewResponse().WithCode(http.StatusNotFound).WithStatus("FAIL").WithEntity("Booking").WithState("DeleteAllBookingFail").WithMessage("Delete All Booking Fail").WithError(err.Error())
	}

	return *appctx.NewResponse().WithCode(http.StatusOK).WithStatus("SUCCESS").WithEntity("Booking").WithState("DeleteAllBookingSuccess").WithMessage("Delete All Booking Success")
}
