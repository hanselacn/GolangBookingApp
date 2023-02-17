package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/repositories"

	"github.com/gorilla/mux"
)

type deleteBooking struct {
	repo repositories.Booking
}

func NewDeleteBooking(repo repositories.Booking) *deleteBooking {
	return &deleteBooking{repo: repo}
}

func (u *deleteBooking) Serve(data *appctx.Data) appctx.Response {
	id := mux.Vars(data.Request)["bookingID"]
	err := u.repo.Delete(data.Request.Context(), id)

	if err != nil {
		err := fmt.Errorf("deleting booking: %w", err)
		fmt.Println(err)

		return *appctx.NewResponse().WithCode(http.StatusNotFound).WithStatus("FAIL").WithEntity("Booking").WithState("DeleteBookingFail").WithMessage("Delete Booking Fail").WithError(err.Error())
	}

	return *appctx.NewResponse().WithCode(http.StatusOK).WithStatus("SUCCESS").WithEntity("Booking").WithState("DeleteBookingSuccess").WithMessage("Delete Booking Success")
}
