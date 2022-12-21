package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa/service/model"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger

	log.Info("invoked ", r.URL.Path)
	var banRequest model.BanRequest
	err := json.NewDecoder(r.Body).Decode(&banRequest)
	if err != nil {
		log.Error("error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ban := rt.memdb.SaveBan(banRequest)
	body, err1 := json.Marshal(ban)
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

func (rt *_router) findAllBans(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	log.Info("invoked ", r.URL.Path)
	username := r.URL.Query().Get("username")
	bans := rt.memdb.FindAllBans(username)
	body, err := json.Marshal(bans)
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

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	log.Info("invoked ", r.URL.Path)
	username := r.Header.Get("X-User-Session-Identifier")
	id := ps.ByName("banId")
	inInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rt.memdb.DeleteBan(username, inInt64)
}
