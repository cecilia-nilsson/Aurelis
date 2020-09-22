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
	"golang.org/x/crypto/bcrypt"
)

type config struct {
	DbUser string `json:"db_user"`
	DbPass string `json:"db_pass"`
	DbName string `json:"db_name"`
	DbHost string `json:"db_host"`
	DbPort string `json:"db_port"`
}

type comment struct {
	PostID  int       `json:"post_id"`
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
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
	Age06       bool      `json:"age_0_6"`
	Age712      bool      `json:"age_7_12"`
	Age1318     bool      `json:"age_13_18"`
	Comments    []comment `json:"comment"`
}

type search struct {
	FreeText    string `json:"free_text"`
	VardagFM    bool   `json:"vardag_fm"`
	VardagEM    bool   `json:"vardag_em"`
	VardagKvall bool   `json:"vardag_kvall"`
	Helg        bool   `json:"helg"`
	Age06       bool   `json:"age_0_6"`
	Age712      bool   `json:"age_7_12"`
	Age1318     bool   `json:"age_13_18"`
}

type user struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var db *sql.DB

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
		p.Comments = listComments(p.ID)
		posts = append(posts, p)
	}

	json.NewEncoder(w).Encode(posts)
}

func listComments(PostID int) []comment {
	comments := []comment{}

	rows, err := db.Query(`
	SELECT
		id, name, message, created
	FROM
		comments
	WHERE
		post_id = $1
	ORDER BY
		created DESC`,
		PostID)
	if err != nil {
		fmt.Println(err)
		return comments
	}

	defer rows.Close()

	for rows.Next() {
		var c comment
		// Read all comments
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Message,
			&c.Created,
		)
		if err != nil {
			log.Fatal(err)
		}
		comments = append(comments, c)
	}

	return comments
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

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data")
		return
	}

	json.Unmarshal(reqBody, &newComment)
	newComment.Created = time.Now()

	_, err = db.Exec(`
	INSERT INTO 
		comments(name, message, created, post_id)
	VALUES($1, $2, $3, $4)`,
		newComment.Name,
		newComment.Message,
		newComment.Created,
		newComment.PostID,
	)
	if err != nil {
		fmt.Println(err, newComment)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newComment)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	var updatePost post
	vars := mux.Vars(r)
	id := vars["id"]

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
		return
	}

	json.Unmarshal(reqBody, &updatePost)

	_, err = db.Exec(`
	UPDATE
		posts
	SET
		title = $1,
		name = $2,
		email = $3,
		image = $4,
		message = $5,
		vardagfm = $6,
		vardagem = $7,
		vardagkvall = $8,
		helg = $9,
		age06 = $10,
		age712 = $11,
		age1318 = $12
	WHERE
		id=$13`,
		updatePost.Title,
		updatePost.Name,
		updatePost.Email,
		updatePost.Image,
		updatePost.Message,
		updatePost.VardagFM,
		updatePost.VardagEM,
		updatePost.VardagKvall,
		updatePost.Helg,
		updatePost.Age06,
		updatePost.Age712,
		updatePost.Age1318,
		id,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	// The following code row is just to give a response to the browser, to avoid complaint.
	json.NewEncoder(w).Encode(map[string]int{"foo": 1, "bar": 2})
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec(`
		DELETE FROM
			posts
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

func makeSQLQuery(searchPost search) string {
	var sqlQuery string

	if searchPost.FreeText != "" {
		sqlQuery = sqlQuery + `(
			title ILIKE '%' || $1 || '%'
		OR
			name ILIKE '%' || $1 || '%'
		OR
			email ILIKE '%' || $1 || '%'
		OR
			message ILIKE '%' || $1 || '%'
		)`
	}

	if searchPost.VardagFM || searchPost.VardagEM || searchPost.VardagKvall || searchPost.Helg {
		var sqlTimeQuery string

		if searchPost.VardagFM {
			sqlTimeQuery = `vardagfm = true `
		}
		if searchPost.VardagEM {
			if sqlTimeQuery != "" {
				sqlTimeQuery = sqlTimeQuery + `OR vardagem = true `
			} else {
				sqlTimeQuery = `vardagem = true `
			}
		}
		if searchPost.VardagKvall {
			if sqlTimeQuery != "" {
				sqlTimeQuery = sqlTimeQuery + `OR vardagkvall = true `
			} else {
				sqlTimeQuery = `vardagkvall = true `
			}
		}
		if searchPost.Helg {
			if sqlTimeQuery != "" {
				sqlTimeQuery = sqlTimeQuery + `OR helg = true `
			} else {
				sqlTimeQuery = `helg = true `
			}
		}

		if sqlQuery != "" && sqlTimeQuery != "" {
			sqlQuery = sqlQuery + ` AND (` + sqlTimeQuery + `)`
		} else if sqlQuery == "" && sqlTimeQuery != "" {
			sqlQuery = `(` + sqlTimeQuery + `)`
		}
	}

	if searchPost.Age06 || searchPost.Age712 || searchPost.Age1318 {
		var sqlAgeQuery string

		if searchPost.Age06 {
			sqlAgeQuery = `age06 = true `
		}
		if searchPost.Age712 {
			if sqlAgeQuery != "" {
				sqlAgeQuery = sqlAgeQuery + `OR age712 = true `
			} else {
				sqlAgeQuery = `age712 = true `
			}
		}
		if searchPost.Age1318 {
			if sqlAgeQuery != "" {
				sqlAgeQuery = sqlAgeQuery + `OR age1318 = true `
			} else {
				sqlAgeQuery = `age1318 = true `
			}
		}

		if sqlQuery != "" && sqlAgeQuery != "" {
			sqlQuery = sqlQuery + ` AND (` + sqlAgeQuery + `)`
		} else if sqlQuery == "" && sqlAgeQuery != "" {
			sqlQuery = `(` + sqlAgeQuery + `)`
		}
	}

	if sqlQuery != "" {
		sqlQuery = `SELECT * FROM posts WHERE ` + sqlQuery + `ORDER BY created DESC`
	} else if sqlQuery == "" {
		sqlQuery = `SELECT * FROM posts ORDER BY created DESC`
	} else {
		fmt.Println("Error in SQL query.")
	}

	return sqlQuery
}

func searchPosts(w http.ResponseWriter, r *http.Request) {
	var searchPost search
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
		return
	}

	json.Unmarshal(reqBody, &searchPost)

	sqlQuery := makeSQLQuery(searchPost)

	var rows *sql.Rows
	// var err error

	if searchPost.FreeText != "" {
		rows, err = db.Query(
			sqlQuery,
			searchPost.FreeText,
		)
	} else {
		rows, err = db.Query(
			sqlQuery,
		)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	posts := []post{}

	defer rows.Close()

	for rows.Next() {
		var p post
		// Read all comments
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
		p.Comments = listComments(p.ID)
		posts = append(posts, p)
	}

	json.NewEncoder(w).Encode(posts)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with your username and password.")
		return
	}

	var newUser user
	json.Unmarshal(reqBody, &newUser)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	INSERT INTO
		users(email, password)
	VALUES
		($1, $2)`,
		newUser.Email,
		hashedPassword,
	)

	if err != nil {
		fmt.Println(err)
	}
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter your email and password.")
		return
	}

	var existingUser user

	json.Unmarshal(reqBody, &existingUser)

	row := db.QueryRow(`
		SELECT
			password
		FROM
			users
		WHERE
			email LIKE $1`,
		existingUser.Email,
	)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		fmt.Println(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(existingUser.Password))
	if err != nil {
		fmt.Println("Problem logging in.")
		return
	}

	fmt.Println("Password matches.")
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
	router.HandleFunc("/searchPosts", searchPosts).Methods("POST")
	router.HandleFunc("/posts/{id:[0-9]+}", deletePost).Methods("DELETE")
	// router.HandleFunc("/posts/{id:[0-9]+}/update", updatePost).Methods("POST")
	router.HandleFunc("/login", loginUser).Methods("POST")
	router.HandleFunc("/registerUser", registerUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedMethods([]string{"POST", "GET", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}))(router)))
}
