package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"github.com/oklays/golang-restapi/config"
	"github.com/oklays/golang-restapi/src/middleware"
	"github.com/oklays/golang-restapi/src/modules/user/repository"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Author Struct
type ResponseData struct {
	Status string `json:"status"`
	Data   string `json:"data"`
	Msg    string `json:"msg"`
}

// Init books var as a slice Book struct
var books []Book
var res ResponseData
var mySigningKey = []byte("HelloBrads")

// Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicatoin/json")
	log.SetPrefix("LOG : ")
	log.Println("Success get all books list")
	json.NewEncoder(w).Encode(books)
}

// Get Single Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicatoin/json")
	params := mux.Vars(r) // Get any params
	// Loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicatoin/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) // Mock ID - not safe
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicatoin/json")
	params := mux.Vars(r) // Get any params
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = strconv.Itoa(rand.Intn(10000000)) // Mock ID - not safe
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

// Delete Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicatoin/json")
	params := mux.Vars(r) // Get any params
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func getAllUserData(w http.ResponseWriter, r *http.Request) {
	log.Println("do process get users")
	db, err := config.GetPostgresDB()

	if err != nil {
		log.Fatal(err)
	}

	userRepositoryPostgres := repository.NewUserRespositoryPostgres(db)

	users, err := userRepositoryPostgres.FindAll()

	e, err := json.Marshal(users)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("JSON data : ", string(e))

	json.NewEncoder(w).Encode(&users)
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// var res Response
		if r.Header["Authorization"] != nil {
			res = ResponseData{Status: "SUCCESS", Data: "null", Msg: "Success Auth"}

			// log token :
			fmt.Println("token is : ", r.Header["Authorization"][0])

			// Split token :
			re := regexp.MustCompile(" ")
			bearerToken := r.Header["Authorization"][0]

			// Split the token  with whitespace
			split := re.Split(bearerToken, -1)
			set := []string{}

			// Looping and append into array
			for i := range split {
				set = append(set, split[i])
			}
			resultToken := set[1]

			token, err := jwt.Parse(resultToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There Was Error on token")
				}
				return mySigningKey, nil
			})

			if err != nil {
				res = ResponseData{Status: "ERROR", Data: "null", Msg: "Token is not valid"}
				fmt.Println(w, err.Error())
				json.NewEncoder(w).Encode(res)
			}

			if token.Valid {
				endpoint(w, r)
			}

			// json.NewEncoder(w).Encode(res)
		} else {
			res = ResponseData{Status: "ERROR", Data: "null", Msg: "Not Authorized User"}
			fmt.Println(w, "Not Authorized User")
			json.NewEncoder(w).Encode(res)

		}
	})
}

func loadRoutes() {
	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	// User Routes
	r.Handle("/api/users", isAuthorized(getAllUserData)).Methods("GET")

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))

}

func main() {
	// Mock Data - @todo - implement DB
	books = append(books, Book{ID: "1", Isbn: "123456", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "777645", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})

	// Generate JWT
	tokenString, err := middleware.GenerateJWT()

	if err != nil {
		log.Fatal("Something wrong when generating JWT token", err.Error())
	}

	log.Println(tokenString)

	loadRoutes()
}
