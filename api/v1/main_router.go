package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func GetPings(w http.ResponseWriter, r *http.Request) {
	times := chi.URLParam(r, "times")
	if times == "" {
		times = "1"
	}
	timed, err := strconv.Atoi(times)
	if err != nil {
		http.Error(w, "Invalid times type. Must be an integer value.", http.StatusBadRequest)
		return
	}
	if timed < 0 {
		timed *= -1
	}
	timed = timed % 100
	times = strings.Repeat("pong ", timed)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"msg\": \"pong\", \"times\": \"%s\"}", strings.TrimSuffix(times, " "))
}

func ApiRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/ping/{times}", GetPings)
	return r
}
