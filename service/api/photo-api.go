package api

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	username := r.Header.Get("X-User-Session-Identifier")
	id := ps.ByName("photoId")
	inInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rt.memdb.DeletePhoto(username, inInt64)

}

func (rt *_router) findUserPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	photoId := ps.ByName("photoId")
	photoIdInt, err := strconv.ParseInt(photoId, 10, 64)
	if err != nil {
		log.Error("error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	photo, err0 := rt.memdb.FindPhoto(photoIdInt)
	if err0 != nil {
		log.Error("error: ", err0)
		http.Error(w, "", http.StatusNotFound)
		return
	}
	body, err1 := json.Marshal(photo)
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

func (rt *_router) findUserPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	//requestorUser := r.Header.Get("X-User-Session-Identifier")
	username := r.URL.Query().Get("username")

	photos := rt.memdb.FindAllPhotos(username)
	body, err := json.Marshal(photos)
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

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := rt.baseLogger
	err0 := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err0 != nil {
		log.Error("error: ", err0)
		http.Error(w, err0.Error(), http.StatusBadRequest)
	}
	username := r.PostFormValue("username")
	_, h, err1 := r.FormFile("file")
	filerc, err2 := h.Open()

	defer func(filerc multipart.File) {
		err := filerc.Close()
		if err != nil {
			log.Error("error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}(filerc)
	if err1 != nil {
		log.Error("error: ", err1)
		http.Error(w, err1.Error(), http.StatusBadRequest)
	}
	if err2 != nil {
		log.Error("error: ", err2)
		http.Error(w, err2.Error(), http.StatusBadRequest)
	}

	buf := new(bytes.Buffer)
	_, err22 := buf.ReadFrom(filerc)
	if err22 != nil {
		log.Error("error: ", err22)
		http.Error(w, err22.Error(), http.StatusInternalServerError)
		return
	}
	//contents := buf.String()
	//encodedPhoto := b64.StdEncoding.EncodeToString([]byte(contents))
	//fmt.Println(encodedPhoto)

	photo := rt.memdb.SavePhoto(username, buf.Bytes())
	body, err3 := json.Marshal(photo)
	if err3 != nil {
		log.Error("error: ", err3)
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}
	_, err4 := w.Write(body)
	if err4 != nil {
		log.Error("error: ", err4)
		http.Error(w, err4.Error(), http.StatusInternalServerError)
		return
	}
	/*
		//Access the name key - First Approach
		fmt.Println(r.Form["username"])
		//Access the name key - Second Approach
		fmt.Println(r.PostForm["username"])
		//Access the name key - Third Approach
		fmt.Println(r.MultipartForm.Value["username"])
		//Access the name key - Fourth Approach
		//
		//Access the name key - First Approach
		fmt.Println(r.Form["file"])
		//Access the name key - Second Approach
		fmt.Println(r.PostForm["file"])
		//Access the name key - Third Approach
		fmt.Println(r.MultipartForm.Value["file"])
		//Access the name key - Fourth Approach
	*/

}
