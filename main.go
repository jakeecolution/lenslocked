package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jakeecolution/lenslocked/controllers"
	"github.com/jakeecolution/lenslocked/migrations"
	"github.com/jakeecolution/lenslocked/models"
	msql "github.com/jakeecolution/lenslocked/models/sql"
	"github.com/jakeecolution/lenslocked/templates"
	"github.com/jakeecolution/lenslocked/views"

	"github.com/joho/godotenv"
)

type Hotel struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
	Price   int    `json:"price"`
}

var hotels []Hotel = make([]Hotel, 0)

func addHotel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var hotel Hotel
	err := decoder.Decode(&hotel)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hotels = append(hotels, hotel)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"msg\": \"Hotel added successfully\", \"status\": \"ok\"}")

}

type WebAdmin struct {
	Name     string
	Username string
	Password string
	Email    string
	Address  string
	Phone    string
	Age      int
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// putConversion is a middleware component that checks for the combination of a
// POST method with a form field named _method having a value of PUT. This is the
// convention for how to establish a PUT method for form submission in a restful API.
// It should run early in the middleware stack.
func putConversion(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				r.ParseForm()
				if r.Form["_method"] != nil && r.FormValue("_method") == "PUT" {
					r.Method = "PUT"
				} else if r.Form["_method"] != nil && r.FormValue("_method") == "DELETE" {
					r.Method = "DELETE"
				}
			}
			next.ServeHTTP(w, r)
		})
}

func main() {
	err := godotenv.Load(".credentials.env")
	CheckErr(err)
	err = godotenv.Load(".env")
	CheckErr(err)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(putConversion)
	//r.Mount("/api/v1", v1.ApiRouter())
	homeTpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	contactTpl := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	myAdmin := WebAdmin{
		Name:     "Jake",
		Username: "admin",
		Password: "password",
		Email:    "jake@ecojake.dev",
		Address:  "123 Main St.",
		Phone:    "123-456-7890",
		Age:      25,
	}
	faqTpl := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	aboutTpl := views.Must(views.ParseFS(templates.FS, "about.gohtml", "tailwind.gohtml"))

	// Users Section
	cfg := msql.PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
		SSLMode:  false,
	}
	db, err := msql.Open(cfg)
	CheckErr(err)
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	var usersC controllers.Users
	usersC.UserService = &models.UserService{
		DB: db,
	}
	usersC.SessionService = &models.SessionService{
		DB: db,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	usersC.Templates.CurrentUser = views.Must(views.ParseFS(templates.FS, "currentuser.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(homeTpl, nil))

	r.Get("/contact", controllers.StaticHandler(contactTpl, myAdmin))
	r.Get("/faq", controllers.FAQ(faqTpl))
	r.Get("/about", controllers.StaticHandler(aboutTpl, myAdmin))
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/users/me", usersC.CurrentUser)
	r.Put("/users/update", usersC.ProcessUserUpdate)
	// r.Post("/hotels", addHotel)
	// r.Get("/hotels", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	jsonBytes, err := json.Marshal(hotels)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	fmt.Println(string(jsonBytes))
	// 	fmt.Fprintf(w, string(jsonBytes))
	// })

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	csrfKey := os.Getenv("CSRF_KEY")
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false),
	)
	log.Println("Starting server at http://localhost:8081...")
	http.ListenAndServe(":8081", csrfMw(r))
}
