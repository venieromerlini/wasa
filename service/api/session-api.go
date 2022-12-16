package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}
