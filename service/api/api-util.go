package api

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Util interface {
	GetAuthToken(r *http.Request) (string, error)
	WriteResponse(w http.ResponseWriter, body []byte) error
	WriteError500(w http.ResponseWriter, err error)
	WriteError404(w http.ResponseWriter, err error)
	WriteError401(w http.ResponseWriter, err error)
	WriteError400(w http.ResponseWriter, err error)
}

type utilImpl struct {
	log logrus.FieldLogger
}

func (u utilImpl) WriteError400(w http.ResponseWriter, err error) {
	u.log.Error("error: ", err)
	http.Error(w, "", http.StatusBadRequest)
}

func (u utilImpl) GetAuthToken(r *http.Request) (string, error) {
	requestorUser := r.Header.Get("X-User-Session-Identifier")
	if requestorUser != "" {
		return requestorUser, nil
	}
	return "", errors.New("auth is required")
}

func (u utilImpl) WriteError500(w http.ResponseWriter, err error) {
	u.log.Error("error: ", err)
	http.Error(w, "", http.StatusInternalServerError)
}

func (u utilImpl) WriteError404(w http.ResponseWriter, err error) {
	u.log.Error("error: ", err)
	http.Error(w, "", http.StatusNotFound)
}

func (u utilImpl) WriteError401(w http.ResponseWriter, err error) {
	u.log.Error("error: ", err)
	http.Error(w, "", http.StatusUnauthorized)
}

func (u utilImpl) WriteResponse(w http.ResponseWriter, body []byte) error {

	w.Header().Set("content-type", "application/json")
	_, err1 := w.Write(body)
	if err1 != nil {
		u.log.Error("error: ", err1)
		http.Error(w, "", http.StatusInternalServerError)
		return err1
	}
	return nil
}

func NewUtil(log logrus.FieldLogger) Util {
	return &utilImpl{
		log: log,
	}
}
