package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) findAllLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}
