package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"database/sql"

	_ "github.com/lib/pq"
)

type config struct {
	DbUser string `json:"db_user"`
	DbPass string `json:"db_pass"`
	DbName string `json:"db_name"`
	DbHost string `json:"db_host"`
	DbPort string `json:"db_port"`
}

type comment struct {
	Name    string    `json:"name"`
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
}

type post struct {
	ID          int       `json:"id"`
	Created     time.Time `json:"created"`
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

var db *sql.DB

// func homeLink(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome home!")
// }

func listPosts(w http.ResponseWriter, r *http.Request) {
	posts := []post{}

	rows, err := db.Query(`
	SELECT
		id, title, name, email, image, message, created, likes, vardagfm, vardagem, vardagkvall, helg, age06, age712, age1318
	FROM
		posts
	ORDER BY
		created DESC`)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var p post
		// Read all comments
		p.Comments = []comment{}
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Name,
			&p.Email,
			&p.Image,
			&p.Message,
			&p.Created,
			&p.Likes,
			&p.VardagFM,
			&p.VardagEM,
			&p.VardagKvall,
			&p.Helg,
			&p.Age06,
			&p.Age712,
			&p.Age1318,
		)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, p)
	}

	json.NewEncoder(w).Encode(posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var newPost post
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
		return
	}

	json.Unmarshal(reqBody, &newPost)
	newPost.Created = time.Now()
	newPost.Likes = 0
	newPost.Comments = []comment{}

	// fmt.Println(newPost.Name)
	// fmt.Println(newPost.VardagFM)
	// fmt.Println(newPost.VardagEM)
	_, err = db.Exec(`
	INSERT INTO 
		posts(title, name, email, image, message, created, likes, vardagfm, vardagem, vardagkvall, helg, age06, age712, age1318)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		newPost.Title,
		newPost.Name,
		newPost.Email,
		newPost.Image,
		newPost.Message,
		newPost.Created,
		newPost.Likes,
		newPost.VardagFM,
		newPost.VardagEM,
		newPost.VardagKvall,
		newPost.Helg,
		newPost.Age06,
		newPost.Age712,
		newPost.Age1318,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	var newComment comment
	vars := mux.Vars(r)
	id := vars["id"]

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data")
		return
	}

	json.Unmarshal(reqBody, &newComment)
	newComment.Date = time.Now()

	_, err = db.Exec(`
	INSERT INTO 
		comments(name, message, created, post_id)
	VALUES($1, $2, $3, $4)`,
		newComment.Name,
		newComment.Message,
		newComment.Date,
		id,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newComment)
}

func likePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec(`
		UPDATE
			posts
		SET
			likes = likes + 1
		WHERE
			id = $1
	`, id)

	if err != nil {
		fmt.Println(err)
		return
	}

	// The following code row is just to give a response to the browser, to avoid complaint.
	json.NewEncoder(w).Encode(map[string]int{"foo": 1, "bar": 2})
}

func main() {
	// Open config file
	configFile, errConfigFile := os.Open("config.json")
	if errConfigFile != nil {
		fmt.Println(errConfigFile)
	}
	defer configFile.Close()

	// Read and implement config file
	decoder := json.NewDecoder(configFile)
	configuration := config{}
	errDecode := decoder.Decode(&configuration)
	if errDecode != nil {
		fmt.Println("error:", errDecode)
	}

	// Open the database
	var err error
	db, err = sql.Open("postgres", "host="+configuration.DbHost+" port="+configuration.DbPort+" user="+configuration.DbUser+" password="+configuration.DbPass+" dbname="+configuration.DbName+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/", homeLink)
	router.HandleFunc("/posts", listPosts).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/comments", createComment).Methods("POST")
	router.HandleFunc("/posts/{id:[0-9]+}/like", likePost).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedMethods([]string{"POST", "GET"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}))(router)))
}
