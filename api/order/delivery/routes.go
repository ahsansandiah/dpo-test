package customerRoutes

import (
	orderRoute "github.com/ahsansandiah/dpo-test/api/order/delivery/route"
	"github.com/ahsansandiah/dpo-test/packages/manager"
	"github.com/gorilla/mux"
)

func NewRoutes(r *mux.Router, mgr manager.Manager) {
	apiAuth := r.PathPrefix("").Subrouter()
	apiAuth.Use(mgr.GetMiddleware().CheckToken)

	orderRoute.NewOrderRoute(mgr, apiAuth)
}
