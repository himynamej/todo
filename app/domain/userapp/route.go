package userapp

import (
	"net/http"

	"github.com/himynamej/todo/app/sdk/auth"
	"github.com/himynamej/todo/app/sdk/authclient"
	"github.com/himynamej/todo/app/sdk/mid"
	"github.com/himynamej/todo/business/domain/userbus"
	"github.com/himynamej/todo/foundation/logger"
	"github.com/himynamej/todo/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	UserBus    *userbus.Business
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)
	ruleAuthorizeUser := mid.AuthorizeUser(cfg.AuthClient, cfg.UserBus, auth.RuleAdminOrSubject)
	ruleAuthorizeAdmin := mid.AuthorizeUser(cfg.AuthClient, cfg.UserBus, auth.RuleAdminOnly)

	api := newApp(cfg.UserBus)

	app.HandlerFunc(http.MethodGet, version, "/users", api.query, authen, ruleAdmin)
	app.HandlerFunc(http.MethodGet, version, "/users/{user_id}", api.queryByID, authen, ruleAuthorizeUser)
	app.HandlerFunc(http.MethodPost, version, "/users", api.create, authen, ruleAdmin)
	app.HandlerFunc(http.MethodPut, version, "/users/role/{user_id}", api.updateRole, authen, ruleAuthorizeAdmin)
	app.HandlerFunc(http.MethodPut, version, "/users/{user_id}", api.update, authen, ruleAuthorizeUser)
	app.HandlerFunc(http.MethodDelete, version, "/users/{user_id}", api.delete, authen, ruleAuthorizeUser)
}
