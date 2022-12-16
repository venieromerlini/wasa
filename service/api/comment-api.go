package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (rt *_router) findAllComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}
