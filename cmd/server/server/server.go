package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"

	handlers "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	rand.Seed(42)

	r := mux.NewRouter()

	// POST /api/register - registering a new user and getting a JWT token
	r.HandleFunc("/api/register", registerUser).Methods("POST")
	// POST /api/login - log in as an existing user and get a JWT token
	r.HandleFunc("/api/login", loginUser).Methods("POST")
	// POST /api/posts/ - adding a post with url or text
	r.HandleFunc("/api/posts", isAuthorized(createPost)).Methods("POST")
	// GET /api/posts/ - get all posts
	r.HandleFunc("/api/posts", getAllPost).Methods("GET")
	// GET /api/user/{USER_LOGIN} - get all post by user
	r.HandleFunc("/api/user/{user_login}", getAllPostByUser).Methods("GET")
	// GET /api/post/{POST_ID} - get post by id
	r.HandleFunc("/api/posts/{post_id}", getPostByID).Methods("GET")
	// GET /a/funny/{CATEGORY_NAME} - get all post by category
	r.HandleFunc("/a/funny/{category_name}", getAllPostByCategory).Methods("GET")
	//r.HandleFunc("/api/posts/{category_name}", getAllPostByCategory).Methods("GET")
	// DELETE /api/post/{POST_ID} - delete post by id
	r.HandleFunc("/api/posts/{post_id}", isAuthorized(deletePostByID)).Methods("DELETE")
	// GET /api/post/{POST_ID}/upvote - rating post up vote by id
	r.HandleFunc("/api/posts/{post_id}/upvote", isAuthorized(upvotePostByID)).Methods("GET")
	// GET /api/post/{POST_ID}/downvote - rating post down vote by id
	r.HandleFunc("/api/posts/{post_id}/downvote", isAuthorized(downvotePostByID)).Methods("GET")
	// POST /api/post/{POST_ID} - add comment to exists post by id
	r.HandleFunc("/api/posts/{post_id}", isAuthorized(addCommentToPostByID)).Methods("POST")
	// DELETE /api/post/{POST_ID}/{COMMENT_ID} - delete comment from exists post by id
	r.HandleFunc("/api/posts/{post_id}/{comment_id}", isAuthorized(deleteCommentToPostByID)).Methods("DELETE")

	fmt.Println("starting server at :8080")

	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}

var notImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	resp, _ := json.Marshal(map[string]interface{}{
		"status": http.StatusInternalServerError,
		"error":  "not implemented",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(resp)
})

func jsonError(w http.ResponseWriter, status int, msg string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}

func jsonErrorRegister(w http.ResponseWriter, status int, user *User, msg string) {
	var errs = Errors{[]Error{
		{Location: "body",
			Param: "username",
			Value: user.Username,
			Msg:   msg},
	}}
	resp, _ := json.Marshal(errs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}

func jsonMessage(w http.ResponseWriter, status int, msg string) {
	var errs = map[string]interface{}{
		"message": msg,
	}
	resp, _ := json.Marshal(errs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}
