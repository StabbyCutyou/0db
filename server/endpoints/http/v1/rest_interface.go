package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/server/node"
	"github.com/gorilla/mux"
	"net/http"
)

func Listen(port int, zdb *node.ZeroDBNode) {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/v1").Subrouter()
	// Write the data to ZDB
	subrouter.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {
		logRequestBegin(r.Method, r.RequestURI)
		key := mux.Vars(r)["key"]
		var buf bytes.Buffer
		buf.ReadFrom(r.Body)
		data := buf.String()
		err := zdb.Write(key, data, true)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			logRequestEnd(r.Method, r.RequestURI, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		} else {
			logRequestEnd(r.Method, r.RequestURI, http.StatusOK)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"key": key, "value": data})
		}

	}).Methods("PUT")

	// Read the data from ZDB
	subrouter.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {
		logRequestBegin(r.Method, r.RequestURI)
		key := mux.Vars(r)["key"]
		data, err := zdb.Read(key)

		if err != nil {
			logRequestEnd(r.Method, r.RequestURI, http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		} else {
			logRequestEnd(r.Method, r.RequestURI, http.StatusOK)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"key": key, "value": data})
		}

	}).Methods("GET")
	logrus.Infof("Starting HTTP Endpoint on port %d", port)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), subrouter))
}

func logRequestBegin(verb string, route string) {
	logrus.Infof("Beginning Request: %s - %s", verb, route)
}

func logRequestEnd(verb string, route string, status int) {
	logrus.Infof("Completed Request: %d %s - %s", status, verb, route)
}
