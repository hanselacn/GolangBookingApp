// Package router
package router

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/bootstrap"
	"GolangBookingApp/internal/consts"
	"GolangBookingApp/internal/handler"
	"GolangBookingApp/internal/middleware"
	"GolangBookingApp/internal/repositories"
	"GolangBookingApp/internal/ucase"
	"GolangBookingApp/pkg/logger"
	"GolangBookingApp/pkg/msgx"
	"GolangBookingApp/pkg/routerkit"

	//"GolangBookingApp/pkg/mariadb"
	//"GolangBookingApp/internal/repositories"
	//"GolangBookingApp/internal/ucase/example"

	ucaseContract "GolangBookingApp/internal/ucase/contract"
)

type router struct {
	config *appctx.Config
	router *routerkit.Router
}

// NewRouter initialize new router wil return Router Interface
func NewRouter(cfg *appctx.Config) Router {
	bootstrap.RegistryMessage()
	bootstrap.RegistryLogger(cfg)

	return &router{
		config: cfg,
		router: routerkit.NewRouter(routerkit.WithServiceName(cfg.App.AppName)),
	}
}

func (rtr *router) handle(hfn httpHandlerFunc, svc ucaseContract.UseCase, mdws ...middleware.MiddlewareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get(consts.HeaderLanguageKey)
		if !msgx.HaveLang(consts.RespOK, lang) {
			lang = rtr.config.App.DefaultLang
			r.Header.Set(consts.HeaderLanguageKey, lang)
		}

		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
				w.WriteHeader(http.StatusInternalServerError)
				res := appctx.Response{
					Code: consts.CodeInternalServerError,
				}

				res.WithLang(lang)
				logger.Error(logger.MessageFormat("error %v", string(debug.Stack())))
				json.NewEncoder(w).Encode(res.Byte())

				return
			}
		}()

		ctx := context.WithValue(r.Context(), "access", map[string]interface{}{
			"path":      r.URL.Path,
			"remote_ip": r.RemoteAddr,
			"method":    r.Method,
		})

		req := r.WithContext(ctx)

		// validate middleware
		if !middleware.FilterFunc(w, req, rtr.config, mdws) {
			return
		}

		resp := hfn(req, svc, rtr.config)
		resp.WithLang(lang)
		rtr.response(w, resp)
	}
}

// response prints as a json and formatted string for DGP legacy
func (rtr *router) response(w http.ResponseWriter, resp appctx.Response) {
	w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
	resp.Generate()
	w.WriteHeader(resp.Code)
	w.Write(resp.Byte())
	return
}

// postgresql
// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "postgres"
// 	dbname   = "intern-privy-simple-orders"
// )

// Route preparing http router and will return mux router object
func (rtr *router) Route() *routerkit.Router {

	root := rtr.router.PathPrefix("/").Subrouter()
	//in := root.PathPrefix("/in/").Subrouter()
	liveness := root.PathPrefix("/").Subrouter()
	//inV1 := in.PathPrefix("/v1/").Subrouter()

	// open tracer setup
	bootstrap.RegistryOpenTracing(rtr.config)

	//db := bootstrap.RegistryMariaMasterSlave(rtr.config.WriteDB, rtr.config.ReadDB, rtr.config.App.Timezone))
	//db := bootstrap.RegistryMariaDB(rtr.config.WriteDB, rtr.config.App.Timezone)

	// use case
	healthy := ucase.NewHealthCheck()

	// healthy
	liveness.HandleFunc("/liveness", rtr.handle(
		handler.HttpRequest,
		healthy,
	)).Methods(http.MethodGet)

	// //CORS HEADER

	// db := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// result :=

	// db := bootstrap.RegistryPostgreSQLMasterSlave(rtr.config.ReadDB, rtr.config.WriteDB, rtr.config.App.Timezone)

	// this is use case for example purpose, please delete
	//repoExample := repositories.NewExample(db)
	//el := example.NewExampleList(repoExample)
	//ec := example.NewPartnerCreate(repoExample
	//ed := example.NewExampleDelete(repoExample)

	db := bootstrap.RegistryPostgreSQLMasterSlave(rtr.config.WriteDB, rtr.config.ReadDB, "Asia/Jakarta")

	userRepository := repositories.NewUserImplementation(db)
	bookingRepository := repositories.NewBookingImplementation(db)
	dayRepository := repositories.NewDayImplementation(db)
	roomRepository := repositories.NewRoomImplementation(db)
	// loginRepository := repositories.NewLoginImplementation(db)
	// logoutRepository := repositories.NewLogoutImplementation(db)

	//USER
	createUser := ucase.NewCreateUser(userRepository)
	root.HandleFunc("/user", rtr.handle(
		handler.HttpRequest,
		createUser,
	)).Methods(http.MethodPost)

	// logoutUser := ucase.NewLogout(logoutRepository)
	// root.HandleFunc("/logout", rtr.handle(
	// 	handler.HttpRequest,
	// 	logoutUser,
	// )).Methods(http.MethodPut)

	// loginUser := ucase.NewLogin(loginRepository)
	// root.HandleFunc("/login", rtr.handle(
	// 	handler.HttpRequest,
	// 	loginUser,
	// )).Methods(http.MethodPut)

	grantAdmin := ucase.NewUpdateAdmin(userRepository)
	root.HandleFunc("/user/grant/{userID}/{username}", rtr.handle(
		handler.HttpRequest,
		grantAdmin,
	)).Methods(http.MethodPut)

	demoteUser := ucase.NewUpdateUser(userRepository)
	root.HandleFunc("/user/demote/{userID}/{username}", rtr.handle(
		handler.HttpRequest,
		demoteUser,
	)).Methods(http.MethodPut)

	getAllUser := ucase.NewGetUsers(userRepository)
	root.HandleFunc("/users", rtr.handle(
		handler.HttpRequest,
		getAllUser,
	)).Methods(http.MethodGet)

	loginUser := ucase.NewLogin(userRepository)
	root.HandleFunc("/users/login", rtr.handle(
		handler.HttpRequest,
		loginUser,
	)).Methods(http.MethodPut)

	logoutUser := ucase.NewLogout(userRepository)
	root.HandleFunc("/users/logout", rtr.handle(
		handler.HttpRequest,
		logoutUser,
	)).Methods(http.MethodPut)

	//BOOKING

	createBooking := ucase.NewCreateBooking(bookingRepository, userRepository)
	root.HandleFunc("/booking", rtr.handle(
		handler.HttpRequest,
		createBooking,
	)).Methods(http.MethodPost)

	deleteBooking := ucase.NewDeleteBooking(bookingRepository)
	root.HandleFunc("/booking/{bookingID}", rtr.handle(
		handler.HttpRequest,
		deleteBooking,
	)).Methods(http.MethodDelete)

	deleteAllBooking := ucase.NewDeleteAllBooking(bookingRepository)
	root.HandleFunc("/bookingreset", rtr.handle(
		handler.HttpRequest,
		deleteAllBooking,
	)).Methods(http.MethodDelete)

	getBooking := ucase.NewGetBookByDay(bookingRepository)
	root.HandleFunc("/booking", rtr.handle(
		handler.HttpRequest,
		getBooking,
	)).Methods(http.MethodGet)

	getBookingby := ucase.NewGetBookByName(bookingRepository)
	root.HandleFunc("/booking/{bookedBy}", rtr.handle(
		handler.HttpRequest,
		getBookingby,
	)).Methods(http.MethodGet)

	//ROOM

	createRoom := ucase.NewCreateRoom(roomRepository, userRepository)
	root.HandleFunc("/room", rtr.handle(
		handler.HttpRequest,
		createRoom,
	)).Methods(http.MethodPost)

	getRooms := ucase.NewGetRooms(roomRepository)
	root.HandleFunc("/room", rtr.handle(
		handler.HttpRequest,
		getRooms,
	)).Methods(http.MethodGet)

	deleteRoom := ucase.NewDeleteRoom(roomRepository, userRepository)
	root.HandleFunc("/deleteroom/{username}/{roomID}", rtr.handle(
		handler.HttpRequest,
		deleteRoom,
	)).Methods(http.MethodDelete)

	//DAY
	createDay := ucase.NewCreateDay(dayRepository, userRepository)
	root.HandleFunc("/day", rtr.handle(
		handler.HttpRequest,
		createDay,
	)).Methods(http.MethodPost)

	getDay := ucase.NewGetDays(dayRepository)
	root.HandleFunc("/day", rtr.handle(
		handler.HttpRequest,
		getDay,
	)).Methods(http.MethodGet)

	// TODO: create your route here

	// this route for example rest, please delete
	// example list
	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    el,
	//)).Methods(http.MethodGet)

	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    ec,
	//)).Methods(http.MethodPost)

	//inV1.HandleFunc("/example/{id:[0-9]+}", rtr.handle(
	//    handler.HttpRequest,
	//    ed,
	//)).Methods(http.MethodDelete)

	return rtr.router

}
