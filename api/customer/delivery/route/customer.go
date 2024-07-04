package customerRoute

import (
	customerHandler "github.com/ahsansandiah/dpo-test/api/customer/delivery/handler"
	"github.com/ahsansandiah/dpo-test/packages/manager"
	"github.com/gorilla/mux"
)

func NewCustomerRoute(mgr manager.Manager, route *mux.Router) {
	customerHandler := customerHandler.NewCustomerHandler(mgr)

	route.Handle("/customers", customerHandler.GetAll()).Methods("GET")
	route.Handle("/customers/{id}", customerHandler.Delete()).Methods("DELETE")
	route.Handle("/customers/{id}", customerHandler.GetByID()).Methods("GET")
	route.Handle("/customers/{id}", customerHandler.Update()).Methods("PATCH")
	route.Handle("/customers", customerHandler.Create()).Methods("POST")
}
