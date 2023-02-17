package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/repositories"
)

type getBookByDay struct {
	repo repositories.Booking
}

func NewGetBookByDay(repo repositories.Booking) *getBookByDay {
	return &getBookByDay{repo: repo}
}

func (u *getBookByDay) Serve(data *appctx.Data) appctx.Response {

	bookedday := data.Request.URL.Query().Get("booked_day")
	bookedroom := data.Request.URL.Query().Get("booked_room")
	books, err := u.repo.GetBookByDayRoom(data.Request.Context(), bookedday, bookedroom)

	if err != nil {
		err := fmt.Errorf("getting books: %w", err)
		fmt.Println(err)

		return *appctx.NewResponse().WithCode(http.StatusInternalServerError).WithStatus("ERROR").WithMessage("Internal Server Error")
	}
	return *appctx.NewResponse().WithCode(http.StatusOK).WithStatus("SUCCESS").WithEntity("Bookings").WithState("GettingBookingSuccess").WithMessage("Getting Bookings Success").WithData(books)
}
