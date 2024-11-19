// Package all binds all the routes into the specified app.
package all

import (
	"time"

	"github.com/himynamej/todo/app/domain/checkapp"
	"github.com/himynamej/todo/app/domain/rawapp"
	"github.com/himynamej/todo/app/domain/todoapp"
	"github.com/himynamej/todo/app/domain/userapp"
	"github.com/himynamej/todo/app/sdk/mux"
	"github.com/himynamej/todo/business/domain/todobus"
	"github.com/himynamej/todo/business/domain/todobus/stores/itemdb"
	"github.com/himynamej/todo/business/domain/userbus"
	"github.com/himynamej/todo/business/domain/userbus/stores/usercache"
	"github.com/himynamej/todo/business/domain/userbus/stores/userdb"
	"github.com/himynamej/todo/business/sdk/delegate"
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
	delegate := delegate.New(cfg.Log)
	userBus := userbus.NewBusiness(cfg.Log, delegate, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB), time.Minute))
	todoBus := todobus.NewBusiness(cfg.Log, itemdb.NewStore(cfg.Log, cfg.DB), cfg.SQSClient, cfg.S3Client)
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	rawapp.Routes(app)

	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    userBus,
		AuthClient: cfg.AuthClient,
	})
	todoapp.Routes(app, todoapp.Config{
		Log:        cfg.Log,
		TodoBus:    todoBus,
		AuthClient: cfg.AuthClient,
	})

}
