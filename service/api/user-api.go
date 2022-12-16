package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
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

	rt.memdb.SaveUser(ps.ByName("username"))
	//name, err := rt.db.GetName()
	user := rt.memdb.FindAllUsers()
	body, err := json.Marshal(user)
	log.Error("error: ", err)
	log.Info("response: ", string(body))
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
