package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasa/service/model"
)

func (rt *_router) findAllUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger

	log.Info("header: ", r.Header)
	log.Info("params: ", ps)
	log.Info("request body: ", r.GetBody)
	user := new(model.User)
	w.Header().Set("content-type", "application/json")
	rt.db.SetName(ps.ByName("username"))

	name, err := rt.db.GetName()
	user.Username = name
	body, err := json.Marshal(user)
	log.Info("name: ", name)
	log.Error("error: ", err)
	log.Info("response: ", string(body))
	_, _ = w.Write(body)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
