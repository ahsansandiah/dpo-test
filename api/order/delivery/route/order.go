package orderRoute

import (
	orderHandler "github.com/ahsansandiah/dpo-test/api/order/delivery/handler"
	"github.com/ahsansandiah/dpo-test/packages/manager"
	"github.com/gorilla/mux"
)

func NewOrderRoute(mgr manager.Manager, route *mux.Router) {
	orderHandler := orderHandler.NewOrderHandler(mgr)

	route.Handle("/orders", orderHandler.GetAll()).Methods("GET")
	route.Handle("/orders/{id}", orderHandler.Delete()).Methods("DELETE")
	route.Handle("/orders/{id}", orderHandler.GetByID()).Methods("GET")
	route.Handle("/orders/{id}", orderHandler.Update()).Methods("PUT")
	route.Handle("/orders", orderHandler.Create()).Methods("POST")
}
