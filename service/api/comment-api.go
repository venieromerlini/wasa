package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa/service/model"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	var commentRequest model.CommentRequest
	err := json.NewDecoder(r.Body).Decode(&commentRequest)
	if err != nil {
		log.Error("error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comment := rt.memdb.SaveComment(commentRequest)
	body, err1 := json.Marshal(comment)
	if err1 != nil {
		log.Error("error: ", err1)
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err2 := w.Write(body)
	if err2 != nil {
		log.Error("error: ", err2)
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

}

func (rt *_router) findAllComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	// requestorUser := r.Header.Get("X-User-Session-Identifier")
	photoId := r.URL.Query().Get("photoId")
	photoIdInt, err0 := strconv.ParseInt(photoId, 10, 64)
	if err0 != nil {
		log.Error("error: ", err0)
		http.Error(w, err0.Error(), http.StatusInternalServerError)
		return
	}
	comments := rt.memdb.FindAllComments(photoIdInt)
	body, err := json.Marshal(comments)
	if err != nil {
		log.Error("error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err1 := w.Write(body)
	if err1 != nil {
		log.Error("error: ", err1)
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	username := r.Header.Get("X-User-Session-Identifier")
	id := ps.ByName("commentId")
	inInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rt.memdb.DeleteComment(username, inInt64)

}
