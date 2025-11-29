package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Router struct {
	Mux *http.ServeMux
}

func (router *Router) handlers(apiCfg *apiConfig, staticDir string, headerMethod map[string]string, endPoints map[string]string) {
	// App Endpoints ------------------------------------------------------------------------------------------------------
	router.Mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(staticDir)))))

	// Api Endpoints --------------------------------------------------------------------------------------------------------------------------------
	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["GET"], endPoints["api"], "/healthz"), func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		io.WriteString(w, "OK")
	})
	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["POST"], endPoints["api"], "/validate_chirp"), func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		chirpBody := struct {
			Body string `json:"body"`
		}{}
		chirpClean := struct {
			Body string `json:"cleaned_body"`
		}{}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			respondWithError(w, 400, "could not read request")
			return
		}

		err = json.Unmarshal(body, &chirpBody)
		if err != nil {
			respondWithError(w, 400, "could not unmarshal request")
			return
		}

		if len(chirpBody.Body) > 140 {
			respondWithError(w, 400, "chirp is too long")
			return
		}

		var createWord string
		valid := true
		wordSlice := strings.Split(chirpBody.Body, " ")
		for index := range wordSlice {
			switch strings.ToLower(wordSlice[index]) {
			case "kerfuffle":
				valid = false
				wordSlice[index] = "****"
			case "sharbert":
				valid = false
				wordSlice[index] = "****"
			case "fornax":
				valid = false
				wordSlice[index] = "****"
			}
		}
		for index, val := range wordSlice {
			if index == 0 {
				createWord = val
				continue
			}
			createWord = createWord + " " + val
		}
		if valid == false {
			chirpClean.Body = createWord
			respondWithJSON(w, 200, chirpClean)
			return
		}
		chirpClean.Body = chirpBody.Body
		respondWithJSON(w, 200, chirpClean)

	})
	// Admin Endpoints --------------------------------------------------------------------------------------------------------------------------------
	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["GET"], endPoints["admin"], "/metrics"), func(w http.ResponseWriter, req *http.Request) {
		metricHTML := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %v times!</p></body></html>", apiCfg.metrics())
		w.WriteHeader(200)
		io.WriteString(w, metricHTML)
	})
	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["POST"], endPoints["admin"], "/reset"), func(w http.ResponseWriter, req *http.Request) {
		apiCfg.resetMetric()
		w.WriteHeader(200)
	})
}
