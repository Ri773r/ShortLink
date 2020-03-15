package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/validator.v2"
)

// App contain router and middleware
type App struct {
	Router      *mux.Router
	Middlewares *Middleware
	config      *Env
}

type route map[string]http.HandlerFunc

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=1"`
}

type shortenResp struct {
	Shortlink string `json:"shortlink"`
}

// Initialize init app
func (a *App) Initialize(e *Env) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	a.Middlewares = &Middleware{}
	a.config = e
	a.initializeRoutes()
}

// initializeRoutes init route
func (a *App) initializeRoutes() {
	m := alice.New(a.Middlewares.LoggingHandler, a.Middlewares.RecoverHandler)
	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortlink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.getShortlinkInfo)).Methods("GET")
	a.Router.Handle("/{shortlink:[a-zA-Z0-9]+}", m.ThenFunc(a.redirect)).Methods("GET")
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

	if err := validator.Validate(req); err != nil {
		respondWithError(w, StatusError{
			http.StatusBadRequest,
			fmt.Errorf("validate parameters failed \n %v", err),
		})
		return
	}

	s, err := a.config.S.Shorten(req.URL, req.ExpirationInMinutes)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJSON(w, http.StatusCreated, shortenResp{Shortlink: s})
	}

}

func (a *App) getShortlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")

	d, err := a.config.S.ShortlinkInfo(s)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJSON(w, http.StatusOK, d)
	}
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s := vars["shortlink"]
	u, err := a.config.S.Unshorten(s)
	if err != nil {
		respondWithError(w, err)
	} else {
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
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
