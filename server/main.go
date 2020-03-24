package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type comment struct {
	Name    string    `json:"name"`
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
}

type post struct {
	ID          int       `json:"ID"`
	Date        time.Time `json:"date"`
	Likes       int       `json:"likes"`
	Title       string    `json:"title"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Image       string    `json:"image"`
	Message     string    `json:"message"`
	VardagFM    bool      `json:"vardag_fm"`
	VardagEM    bool      `json:"vardag_em"`
	VardagKvall bool      `json:"vardag_kvall"`
	Helg        bool      `json:"helg"`
	Age06       bool      `json:"0-6"`
	Age712      bool      `json:"7-12"`
	Age1318     bool      `json:"13-18"`
	Comments    []comment `json:"comment"`
}

var posts = []post{
	{
		ID:       1,
		Date:     time.Now(),
		Likes:    0,
		Title:    "Hjälp!",
		Name:     "Daniel Karlsson",
		Email:    "mail@danielk.se",
		Image:    "https://placedog.net/400/600",
		Message:  "Aenean et erat ac justo vehicula volutpat. Ut diam lorem, suscipit ut ex quis, accumsan auctor erat. Etiam iaculis purus vitae dolor convallis, ac consequat nisl tincidunt. Nullam imperdiet iaculis mattis. Fusce tincidunt dui lectus, nec consectetur tellus ornare id. Phasellus ac elementum risus. Sed consectetur viverra purus, sed eleifend augue aliquet in. Morbi ultrices eleifend felis, vel luctus ex euismod eget. Mauris massa ligula, ornare ut ipsum sed, tempus auctor quam. Donec ut quam scelerisque, varius nunc eu, congue erat. Pellentesque hendrerit ac ex ut feugiat",
		VardagFM: false,
		Comments: []comment{},
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func listPosts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var newPost post
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newPost)
	newPost.Date = time.Now()
	newPost.Comments = []comment{}
	posts = append(posts, newPost)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/posts", listPosts).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedMethods([]string{"POST", "GET"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}))(router)))
}
