package rpc

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vitelabs/vite-portal/relayer/internal/app"
	"github.com/vitelabs/vite-portal/relayer/version"
)

func Name(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	WriteResponse(w, app.AppName, ContentTypeTextPlain)
}

func Version(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	WriteResponse(w, version.PROJECT_BUILD_VERSION, ContentTypeTextPlain)
}
