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
		respondWithError(w, StatusError{
			http.StatusBadRequest,
			fmt.Errorf("parse parameters failed \n %v", r.Body),
		})
		return
	}
	defer r.Body.Close()
	v := validate.New(&req)
	if v.Validate() {
		fmt.Println(req)
	} else {
		respondWithError(w, StatusError{
			http.StatusBadRequest,
			fmt.Errorf("validate parameters failed \n %v", v.Errors),
		})
		return
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

func respondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d - %s", e.Status(), e.Error())
		respondWithJSON(w, e.Status(), e.Error())
	default:
		respondWithJSON(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

// Run run app
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
