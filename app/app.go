package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gookit/validate"
	"github.com/gorilla/mux"
)

// App contain router and middleware
type App struct {
	Router *mux.Router
}

type route map[string]http.HandlerFunc

type shortenReq struct {
	URL                 string `json:"url" validate:"required"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min:0"`
}

type shortenResp struct {
	Shortlink string `json:"shortlink"`
}

// Initialize init app
func (a *App) Initialize() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// initializeRoutes init route
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/shorten", a.createShortlink).Methods("POST")
	a.Router.HandleFunc("/api/info", a.getShortlinkInfo).Methods("GET")
	a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]+}", a.redirect).Methods("GET")
}

func (a *App) createShortlink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	defer r.Body.Close()
	v := validate.New(&req)
	if v.Validate() {
		fmt.Println(req)
	} else {
		fmt.Println(v.Errors)
	}
	fmt.Println("createShortlink")
}

func (a *App) getShortlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")

	fmt.Println("getShortlinkInfo", s)
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("getShortlinkInfo", vars["shortlink"])
}

// Run run app
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
