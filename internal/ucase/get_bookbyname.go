package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/repositories"

	"github.com/gorilla/mux"
)

type getBookByName struct {
	repo repositories.Booking
}

func NewGetBookByName(repo repositories.Booking) *getBookByName {
	return &getBookByName{repo: repo}
}

func (u *getBookByName) Serve(data *appctx.Data) appctx.Response {

	bookedby := mux.Vars(data.Request)["bookedBy"]
	books, err := u.repo.GetBookByName(data.Request.Context(), bookedby)

	if err != nil {
		err := fmt.Errorf("getting books: %w", err)
		fmt.Println(err)

		return *appctx.NewResponse().WithCode(http.StatusInternalServerError).WithStatus("ERROR").WithMessage("Internal Server Error")
	}
	return *appctx.NewResponse().WithCode(http.StatusOK).WithStatus("SUCCESS").WithEntity("Bookings").WithState("GettingBookingSuccess").WithMessage("Getting Bookings Success").WithData(books)
}
