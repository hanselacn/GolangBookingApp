package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/consts"
	"GolangBookingApp/internal/entity"
	"GolangBookingApp/internal/repositories"
	"GolangBookingApp/pkg/logger"
)

type logout struct {
	userRepository repositories.User
}

func NewLogout(userRepository repositories.User) *logout {
	return &logout{
		userRepository: userRepository,
	}
}
func (u *logout) Serve(data *appctx.Data) appctx.Response {
	p := &entity.User{}

	err := data.Cast(&p)

	if err != nil {
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("User").WithState("LoginFailed").WithMessage("Logout Fail").WithError(err.Error())
	}

	fmt.Println(p.Username)

	user, _ := u.userRepository.VerifyLogout(data.Request.Context(), p.Username)

	if err != nil {
		logger.Error(err)
		return *appctx.NewResponse().WithCode(consts.CodeBadRequest).WithStatus("FAIL").WithEntity("User").WithState("LogoutFailed").WithMessage("Logout failed").WithError(err.Error())
	}

	return *appctx.NewResponse().WithCode(http.StatusCreated).WithStatus("SUCCESS").WithEntity("Users").WithState("LogoutSuccess").WithMessage("Logout Success").WithData(user)
}
