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

	"github.com/google/uuid"
)

type createUser struct {
	userRepository repositories.User
}

func NewCreateUser(userRepository repositories.User) *createUser {
	return &createUser{
		userRepository: userRepository,
	}
}
func (u *createUser) Serve(data *appctx.Data) appctx.Response {
	p := &entity.User{}

	err := data.Cast(&p)

	if err != nil {
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("User").WithState("UserAddFailed").WithMessage("Adding User Fail").WithError(err.Error())
	}

	p.Password, err = encrypt.HashPassword(p.Password)

	if err != nil {
		logger.Error(err)
		return *appctx.NewResponse().WithCode(consts.CodeBadRequest).WithStatus("FAIL").WithEntity("User").WithState("UserAddFailed").WithMessage("Adding User Fail").WithError(err.Error())
	}

	id, err := u.userRepository.Create(data.Request.Context(), p)

	if err != nil {
		fmt.Println(err)
		return *appctx.NewResponse().WithCode(http.StatusBadRequest).WithStatus("FAIL").WithEntity("Users").WithState("UserAddFailed").WithMessage("Adding User Fail").WithError(err.Error())
	}

	type resp struct {
		ID uuid.UUID `json:"id"`
	}

	var res resp
	res.ID = *id

	return *appctx.NewResponse().WithCode(http.StatusCreated).WithStatus("SUCCESS").WithEntity("Users").WithState("AddingUserSuccess").WithMessage("Adding User Success").WithData(res)
}

// func (uc createUser) Serve(data *appctx.Data) appctx.Response {
// 	var reqbody struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 		Name     string `json:"name"`
// 		Role     string `json:"role"`
// 	}

// 	err := json.NewDecoder(data.Request.Body).Decode(&reqbody)
// 	if err != nil {
// 		err = fmt.Errorf("decoding json body : %w", err)
// 		return *appctx.NewResponse().
// 			WithCode(http.StatusInternalServerError).
// 			WithError(err)
// 	}

// 	h := sha256.New()
// 	h.Write([]byte(reqbody.Password))
// 	hmac64 := base64.StdEncoding.EncodeToString(h.Sum(nil))
// 	reqbody.Password = hmac64

// 	user := entity.User{
// 		Username: reqbody.Username,
// 		Password: reqbody.Password,
// 		Name:     reqbody.Name,
// 		Role:     reqbody.Role,
// 	}

// 	id, err = uc.userRepository.Create(data.Request.Context(), &user)

// 	type resp struct {
// 		ID uuid.UUID `json:"id"`
// 	}

// 	var res resp
// 	res.ID = *id

// 	if err != nil {
// 		err = fmt.Errorf("creating user : %w", err)
// 		return *appctx.NewResponse().
// 			WithCode(http.StatusInternalServerError).
// 			WithError(err).
// 			WithData(id)
// 	}

// 	return *appctx.NewResponse().
// 		WithCode(http.StatusCreated).
// 		WithData(user)
// }
