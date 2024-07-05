package userRoute

import (
	userHandler "github.com/ahsansandiah/dpo-test/api/user/delivery/handler"
	"github.com/ahsansandiah/dpo-test/packages/manager"
	"github.com/gorilla/mux"
)

func NewUserRoute(mgr manager.Manager, route *mux.Router) {
	userHandler := userHandler.NewUserHandler(mgr)

	// user
	route.Handle("/users", userHandler.Detail()).Methods("GET")
	route.Handle("/users", userHandler.Create()).Methods("POST")

	// authentication
	route.Handle("/auth/login", userHandler.Login()).Methods("POST")
}
