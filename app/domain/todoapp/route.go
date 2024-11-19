package todoapp

import (
	"net/http"

	"github.com/himynamej/todo/app/sdk/authclient"
	"github.com/himynamej/todo/app/sdk/mid"
	"github.com/himynamej/todo/business/domain/todobus"
	"github.com/himynamej/todo/foundation/logger"
	"github.com/himynamej/todo/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	TodoBus    *todobus.Business
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	//	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.TodoBus)
	app.HandlerFunc(http.MethodPost, version, "/todo", api.CreateTodoItem, authen)
	app.HandlerFunc(http.MethodPost, version, "/upload", api.UploadFile, authen)
	app.HandlerFunc(http.MethodGet, version, "/download/{file_id}", api.DownloadFile, authen)
}
