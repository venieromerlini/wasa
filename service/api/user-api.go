package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasa/service/model"
)

func (rt *_router) findAllUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	users := rt.memdb.FindAllUsers()
	body, err1 := json.Marshal(users)
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

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	log.Info("invoked ", r.URL.Path)
	username := ps.ByName("username")
	homepage := rt.memdb.FindUserHomePageByUsername(username)
	body, err1 := json.Marshal(homepage)
	if err1 != nil {
		log.Error("error: ", err1)
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	_, err2 := w.Write(body)
	if err2 != nil {
		log.Error("error: ", err2)
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	log.Info("invoked ", r.URL.Path)
	username := ps.ByName("username")
	userProfile := rt.memdb.FindUserProfileByUsername(username)
	body, err1 := json.Marshal(userProfile)
	if err1 != nil {
		log.Error("error: ", err1)
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	_, err2 := w.Write(body)
	if err2 != nil {
		log.Error("error: ", err2)
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	requestorUser := r.Header.Get("X-User-Session-Identifier")
	oldUsername := ps.ByName("username")

	var newUser model.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Error("error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := rt.memdb.UpdateUsername(requestorUser, oldUsername, newUser.Username)
	w.Header().Set("content-type", "application/json")
	body, err1 := json.Marshal(user)
	if err1 != nil {
		log.Error("error: ", err1)
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	_, err2 := w.Write(body)
	if err2 != nil {
		log.Error("error: ", err2)
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
}
