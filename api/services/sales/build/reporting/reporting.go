// Package reporting binds the reporting domain set of routes into the specified app.
package reporting

import (
	"github.com/himynamej/todo/app/domain/checkapp"
	"github.com/himynamej/todo/app/sdk/mux"
	"github.com/himynamej/todo/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {

	// Construct the business domain packages we need here so we are using the
	// sames instances for the different set of domain apis.

	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

}
