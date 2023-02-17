package ucase

import (
	"fmt"
	"net/http"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/consts"
	"GolangBookingApp/internal/entity"
	"GolangBookingApp/internal/helper/encrypt"
	"GolangBookingApp/internal/repositories"
	"GolangBookingApp/pkg/logger"
)

type getLogin struct {
	userRepository repositories.User
}

func GetLogin(userRepository repositories.User) *getLogin {
	return &getLogin{
		userRepository: userRepository,
	}
}
func (u *getLogin) Serve(data *appctx.Data) appctx.Response {
	p := &entity.User{}

	err := data.Cast(&p)

	if err != nil {
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("User").WithState("UserAddFailed").WithMessage("Adding User Fail").WithError(err.Error())
	}

	user := entity.User{}

	fmt.Println(p.Username)

	user, _ = u.userRepository.VerifyLogin(data.Request.Context(), p.Username)

	fmt.Println(user)
	fmt.Println(p.Password)
	fmt.Println(user.Password)

	success := encrypt.CheckPasswordHash(user.Password, p.Password)

	if !success {
		return *appctx.NewResponse().WithCode(http.StatusUnauthorized).WithStatus("FAIL").WithEntity("User").WithState("AuthenticationFail").WithMessage("Authentication Fail")
	}

	if err != nil {
		logger.Error(err)
		return *appctx.NewResponse().WithCode(consts.CodeBadRequest).WithStatus("FAIL").WithEntity("User").WithState("UserAddFailed").WithMessage("Adding User Fail").WithError(err.Error())
	}

	return *appctx.NewResponse().WithCode(http.StatusCreated).WithStatus("SUCCESS").WithEntity("Users").WithState("AddingUserSuccess").WithMessage("Adding User Success").WithData(user)
}
