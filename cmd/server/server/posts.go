package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	mux "github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

//Post type
type Post struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Category         string        `json:"category" bson:"category"`
	Type             string        `json:"type" bson:"type"`
	Title            string        `json:"title" bson:"title"`
	URL              string        `json:"url,omitempty" bson:"url"`
	Text             string        `json:"text,omitempty" bson:"text"`
	Author           Author        `json:"author" bson:"author"`
	Comments         []Comment     `json:"comments" bson:"comments"`
	Created          string        `json:"created" bson:"create"`
	Score            int           `json:"scope" bson:"scope"`
	Views            int           `json:"views" bson:"views"`
	UpvotePercentage int           `json:"upvotePercentage" bson:"upvotePercentage"`
	Votes            []Vote        `json:"votes" bson:"votes"`
}

//Author type
type Author struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username string        `json:"username" bson:"username"`
}

//Vote type
type Vote struct {
	User bson.ObjectId `json:"user" bson:"_id,omitempty"`
	Vote int           `json:"vote" bson:"vote"`
}

//Comment type
type Comment struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Author  Author        `json:"author" bson:"author"`
	Created string        `json:"created" bson:"create"`
	Body    string        `json:"body" bson:"body"`
	Comment string        `json:"comment,omitempty"`
}

func (a *Author) fullFromJWTToken(token *jwt.Token) error {
	if a == nil {
		return errors.New("author type not initialized")
	}

	var (
		user []byte
		err  error
	)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("none claims")
	}

	if user, err = json.Marshal(claims["user"]); err != nil {
		return err
	}

	if err = json.Unmarshal(user, a); err != nil {
		return err
	}

	return nil
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var (
		body  []byte
		resp  []byte
		db    *mgo.Database
		token *jwt.Token
		err   error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	post := &Post{}
	if err = json.Unmarshal(body, post); err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	if token, err = getToken(r); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	author := &Author{}
	if err := author.fullFromJWTToken(token); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if status, err := postModel.create(post, author); status != http.StatusOK || err != nil {
		jsonError(w, status, err.Error())
		return
	}

	if resp, err = json.Marshal(post); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func deletePostByID(w http.ResponseWriter, r *http.Request) {
	var (
		db    *mgo.Database
		token *jwt.Token
		err   error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	vars := mux.Vars(r)
	postID := vars["post_id"]
	if postID == "" {
		jsonError(w, http.StatusBadRequest, "none post id param")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	if token, err = getToken(r); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	author := &Author{}
	if err := author.fullFromJWTToken(token); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	post := &Post{ID: bson.ObjectId(postID)}
	if status, err := postModel.deleteByID(post, author); status != http.StatusOK || err != nil {
		jsonMessage(w, status, err.Error())
		return
	}

	jsonMessage(w, http.StatusOK, "success")
}

func getAllPost(w http.ResponseWriter, r *http.Request) {
	var (
		resp []byte
		db   *mgo.Database
		err  error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	posts, status, err := postModel.getAll()
	if status != http.StatusOK || err != nil {
		jsonError(w, status, err.Error())
		return
	}

	if resp, err = json.Marshal(posts); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusNotModified)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func getPostByID(w http.ResponseWriter, r *http.Request) {
	var (
		resp []byte
		db   *mgo.Database
		err  error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	vars := mux.Vars(r)
	postID := vars["post_id"]
	if postID == "" {
		jsonError(w, http.StatusBadRequest, "none post id param")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	post := &Post{ID: bson.ObjectId(postID)}
	if status, err := postModel.getByID(post); http.StatusOK != http.StatusOK || err != nil {
		jsonMessage(w, status, err.Error())
		return
	}

	if resp, err = json.Marshal(post); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func getAllPostByUser(w http.ResponseWriter, r *http.Request) {
	var (
		resp []byte
		db   *mgo.Database
		err  error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	vars := mux.Vars(r)
	username := vars["user_login"]
	if username == "" {
		jsonError(w, http.StatusBadRequest, "none user name param")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	author := &Author{Username: username}
	posts, status, err := postModel.getAllByUser(author)
	if status != http.StatusOK || err != nil {
		jsonError(w, status, err.Error())
		return
	}

	if resp, err = json.Marshal(posts); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func getAllPostByCategory(w http.ResponseWriter, r *http.Request) {
	var (
		resp []byte
		db   *mgo.Database
		err  error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	vars := mux.Vars(r)
	categoryname := vars["category_name"]
	if categoryname == "" {
		jsonError(w, http.StatusBadRequest, "none category name param")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	posts, status, err := postModel.getAllByCategory(categoryname)
	if status != http.StatusOK || err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if resp, err = json.Marshal(posts); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func upvotePostByID(w http.ResponseWriter, r *http.Request) {
	var (
		db    *mgo.Database
		resp  []byte
		token *jwt.Token
		err   error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	vars := mux.Vars(r)
	postID := vars["post_id"]
	if postID == "" {
		jsonError(w, http.StatusBadRequest, "none post id param")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	if token, err = getToken(r); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	author := &Author{}
	if err := author.fullFromJWTToken(token); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	post := &Post{ID: bson.ObjectId(postID)}
	if status, err := postModel.upvoteByID(post, author); status != http.StatusOK || err != nil {
		jsonMessage(w, status, err.Error())
		return
	}

	if resp, err = json.Marshal(post); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func downvotePostByID(w http.ResponseWriter, r *http.Request) {
	var (
		db    *mgo.Database
		resp  []byte
		token *jwt.Token
		err   error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	vars := mux.Vars(r)
	postID := vars["post_id"]
	if postID == "" {
		jsonError(w, http.StatusBadRequest, "none post id param")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	if token, err = getToken(r); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	author := &Author{}
	if err := author.fullFromJWTToken(token); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	post := &Post{ID: bson.ObjectId(postID)}
	if status, err := postModel.downvoteByID(post, author); status != http.StatusOK || err != nil {
		jsonMessage(w, status, err.Error())
		return
	}

	if resp, err = json.Marshal(post); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func addCommentToPostByID(w http.ResponseWriter, r *http.Request) {
	var (
		body  []byte
		resp  []byte
		db    *mgo.Database
		token *jwt.Token
		err   error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	comment := &Comment{}
	if err = json.Unmarshal(body, comment); err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	vars := mux.Vars(r)
	postID := vars["post_id"]
	if postID == "" {
		jsonError(w, http.StatusBadRequest, "none post id param")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	if token, err = getToken(r); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	author := &Author{}
	if err := author.fullFromJWTToken(token); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	post := &Post{ID: bson.ObjectId(postID)}
	if status, err := postModel.addCommentByID(post, author, comment); status != http.StatusOK || err != nil {
		jsonMessage(w, status, err.Error())
		return
	}

	if resp, err = json.Marshal(post); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func deleteCommentToPostByID(w http.ResponseWriter, r *http.Request) {
	var (
		resp  []byte
		db    *mgo.Database
		token *jwt.Token
		err   error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	vars := mux.Vars(r)
	postID := vars["post_id"]
	if postID == "" {
		jsonError(w, http.StatusBadRequest, "none post id param")
		return
	}

	commentID := vars["comment_id"]
	if commentID == "" {
		jsonError(w, http.StatusBadRequest, "none comment id param")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	if token, err = getToken(r); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	author := &Author{}
	if err := author.fullFromJWTToken(token); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	post := &Post{ID: bson.ObjectId(postID)}
	comment := &Comment{ID: bson.ObjectId(commentID)}
	if status, err := postModel.deleteCommentByID(post, author, comment); status != http.StatusOK || err != nil {
		jsonMessage(w, status, err.Error())
		return
	}

	if resp, err = json.Marshal(post); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	w.Write([]byte("\n\n"))
}
