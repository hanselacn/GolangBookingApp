package ucase

import (
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/consts"
	"GolangBookingApp/internal/entity"
	"GolangBookingApp/internal/helper/encrypt"
	"GolangBookingApp/internal/repositories"
	"GolangBookingApp/pkg/logger"
)

type login struct {
	userRepository repositories.User
}

func NewLogin(userRepository repositories.User) *login {
	return &login{
		userRepository: userRepository,
	}
}
func (u *login) Serve(data *appctx.Data) appctx.Response {
	p := &entity.User{}

	err := data.Cast(&p)

	if err != nil {
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("User").WithState("UserLoginFailed").WithMessage("Login Fail").WithError(err.Error())
	}

	user := entity.User{}

	user, _ = u.userRepository.VerifyLogin(data.Request.Context(), p.Username)

	success := encrypt.CheckPasswordHash(user.Password, p.Password)

	if !success {
		return *appctx.NewResponse().WithCode(http.StatusUnauthorized).WithStatus("FAIL").WithEntity("User").WithState("AuthenticationFail").WithMessage("Authentication Fail")
	}

	if err != nil {
		logger.Error(err)
		return *appctx.NewResponse().WithCode(consts.CodeBadRequest).WithStatus("FAIL").WithEntity("User").WithState("UserLoginFailed").WithMessage("Login Fail").WithError(err.Error())
	}

	return *appctx.NewResponse().WithCode(http.StatusCreated).WithStatus("SUCCESS").WithEntity("Users").WithState("UserLoginSuccess").WithMessage("Login Success").WithData(user)
}
