package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasa/service/model"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	log.Info("invoked ", r.URL.Path)

	var sessionBody model.SessionBody
	err := json.NewDecoder(r.Body).Decode(&sessionBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Info("session object :", sessionBody)
	user := rt.memdb.SaveUser(sessionBody.Name)
	successfulLogin := new(model.SuccessfulLogin)
	successfulLogin.Identifier = user.Username
	body, err1 := json.Marshal(successfulLogin)
	if err1 != nil {
		log.Error("error: ", err)
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("response: ", string(body))
	w.Header().Set("X-user-session-identifier", successfulLogin.Identifier)
	_, err2 := w.Write(body)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
}
