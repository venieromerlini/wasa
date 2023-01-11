package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasa/service/model"
)

func (rt *_router) findAllUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	log.Info("invoked ", r.URL.Path)

	users := rt.memdb.FindAllUsers()
	body, errParse := json.Marshal(users)
	if errParse != nil {
		rt.util.WriteError500(w, errParse)
		return
	}

	errResponse := rt.util.WriteResponse(w, body)
	if errResponse != nil {
		rt.util.WriteError500(w, errResponse)
	}
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	log.Info("invoked ", r.URL.Path)
	_, errAuth := rt.util.GetAuthToken(r)
	if errAuth != nil {
		rt.util.WriteError401(w, errAuth)
		return
	}
	username := ps.ByName("username")
	homepage, errDb := rt.memdb.FindUserHomePageByUsername(username)
	if errDb != nil {
		rt.util.WriteError404(w, errDb)
		return
	}

	body, errParse := json.Marshal(homepage)
	if errParse != nil {
		rt.util.WriteError500(w, errParse)
		return
	}

	errResponse := rt.util.WriteResponse(w, body)
	if errResponse != nil {
		rt.util.WriteError500(w, errResponse)
	}
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	log.Info("invoked ", r.URL.Path)
	_, errAuth := rt.util.GetAuthToken(r)
	if errAuth != nil {
		rt.util.WriteError401(w, errAuth)
		return
	}
	username := ps.ByName("username")
	userProfile, errDb := rt.memdb.FindUserProfileByUsername(username)
	if errDb != nil {
		rt.util.WriteError404(w, errDb)
		return
	}

	body, errParse := json.Marshal(userProfile)
	if errParse != nil {
		rt.util.WriteError500(w, errParse)
		return
	}

	errResponse := rt.util.WriteResponse(w, body)
	if errResponse != nil {
		rt.util.WriteError500(w, errResponse)
	}
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	log.Info("invoked ", r.URL.Path)
	requestorUser, err := rt.util.GetAuthToken(r)
	if err != nil {
		rt.util.WriteError401(w, err)
	}
	oldUsername := ps.ByName("username")
	var newUser model.User
	errDecode := json.NewDecoder(r.Body).Decode(&newUser)
	if errDecode != nil {
		rt.util.WriteError400(w, errDecode)
		return
	}

	user, errUpdateUsername := rt.memdb.UpdateUsername(requestorUser, oldUsername, newUser.Username)
	if errUpdateUsername != nil {
		rt.util.WriteError409(w, errUpdateUsername)
		return
	}

	body, errParse := json.Marshal(user)

	if errParse != nil {
		rt.util.WriteError500(w, errParse)
		return
	}

	errResponse := rt.util.WriteResponse(w, body)
	if errResponse != nil {
		rt.util.WriteError500(w, errResponse)
	}

}
