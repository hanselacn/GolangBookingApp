package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/repositories"
)

type getDays struct {
	repo repositories.Day
}

func NewGetDays(repo repositories.Day) *getDays {
	return &getDays{repo: repo}
}

func (u *getDays) Serve(data *appctx.Data) appctx.Response {

	days, err := u.repo.GetDays(data.Request.Context())

	if err != nil {
		err := fmt.Errorf("getting days: %w", err)
		fmt.Println(err)

		return *appctx.NewResponse().WithCode(http.StatusInternalServerError).WithStatus("ERROR").WithMessage("Internal Server Error")
	}
	return *appctx.NewResponse().WithCode(http.StatusOK).WithStatus("SUCCESS").WithEntity("Days").WithState("GettingDaySuccess").WithMessage("Getting Days Success").WithData(days)
}
