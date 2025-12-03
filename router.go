package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Reece-Reklai/go_serve/internal/auth"
	"github.com/Reece-Reklai/go_serve/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Router struct {
	Mux *http.ServeMux
}
type SingleChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
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
	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["POST"], endPoints["api"], "/login"), func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		authenticateUser := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		userJSON := struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
		}{}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			respondWithError(w, 400, "failed to read request body")
			return
		}
		err = json.Unmarshal(body, &authenticateUser)
		if err != nil {
			respondWithError(w, 400, "failed to unmarshal request")
			return
		}
		user, err := apiCfg.databaseQuery.GetUserByEmail(req.Context(), authenticateUser.Email)
		if err != nil {
			error := fmt.Sprintf("failed: %v", err)
			respondWithError(w, 500, error)
			return
		}
		match, err := auth.CheckPasswordHash(authenticateUser.Password, user.Password)
		if err != nil {
			respondWithError(w, 500, "failed to check password match")
			return
		}
		if match == true {
			userJSON.ID = user.ID
			userJSON.CreatedAt = user.CreatedAt
			userJSON.UpdatedAt = user.UpdatedAt
			userJSON.Email = user.Email
			respondWithJSON(w, 200, userJSON)
			return
		} else {
			respondWithError(w, 500, "authentication failed")
			return
		}
	})
	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["POST"], endPoints["api"], "/users"), func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		getUSER := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		userJSON := struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
		}{}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			respondWithError(w, 400, "failed to read request body")
			return
		}
		err = json.Unmarshal(body, &getUSER)
		if err != nil {
			respondWithError(w, 400, "failed to unmarshal request")
			return
		}
		hashPassword, err := auth.HashPassword(getUSER.Password)
		if err != nil {
			respondWithError(w, 400, "failed to hash")
			return
		}
		user, err := apiCfg.databaseQuery.CreateUser(req.Context(), database.CreateUserParams{ID: uuid.New(), Email: getUSER.Email, Password: hashPassword})
		if err != nil {
			error := fmt.Sprintf("failed: %v", err)
			respondWithError(w, 500, error)
			return
		}
		userJSON.ID = user.ID
		userJSON.CreatedAt = user.CreatedAt
		userJSON.UpdatedAt = user.UpdatedAt
		userJSON.Email = user.Email
		response, err := json.Marshal(userJSON)
		if err != nil {
			respondWithError(w, 400, "failed to marshal json into bytes")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(response)
	})
	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["GET"], endPoints["api"], "/chirps/{chirpID}"), func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var responseChirp SingleChirp
		chirpID := req.PathValue("chirpID")
		if chirpID == "" {
			respondWithError(w, 404, "no match found in request")
			return
		}
		stringToUUID, err := uuid.ParseBytes([]byte(chirpID))
		if err != nil {
			respondWithError(w, 400, "failed to parse request")
			return
		}
		chirp, err := apiCfg.databaseQuery.GetChirpById(req.Context(), stringToUUID)
		if err != nil {
			error := fmt.Sprintf("failed: %v", err)
			respondWithError(w, 500, error)
			return
		}
		responseChirp.ID = chirp.ID
		responseChirp.CreatedAt = chirp.CreatedAt
		responseChirp.UpdatedAt = chirp.UpdatedAt
		responseChirp.Body = chirp.Body
		responseChirp.UserID = chirp.UserID
		respondWithJSON(w, 200, responseChirp)
	})

	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["GET"], endPoints["api"], "/chirps"), func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var responseChirp SingleChirp
		chirps, err := apiCfg.databaseQuery.GetAllChirps(req.Context())
		if err != nil {
			error := fmt.Sprintf("failed: %v", err)
			respondWithError(w, 500, error)
			return
		}
		allChirps := make([]SingleChirp, 0)
		for _, value := range chirps {
			responseChirp.ID = value.ID
			responseChirp.CreatedAt = value.CreatedAt
			responseChirp.UpdatedAt = value.UpdatedAt
			responseChirp.Body = value.Body
			responseChirp.UserID = value.UserID
			allChirps = append(allChirps, responseChirp)
		}
		respondWithJSON(w, 200, allChirps)
	})
	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["POST"], endPoints["api"], "/chirps"), func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		chirpProfanity := struct {
			ID     uuid.UUID `json:"id"`
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}{}
		chirpClean := struct {
			ID     uuid.UUID `json:"id"`
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}{}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			respondWithError(w, 400, "failed to read request body")
			return
		}

		err = json.Unmarshal(body, &chirpProfanity)
		if err != nil {
			respondWithError(w, 400, "failed to unmarshal request")
			return
		}

		if len(chirpProfanity.Body) > 100 {
			respondWithError(w, 400, "chirp is too long (100 characters)")
			return
		}

		var createWord string
		valid := true
		wordSlice := strings.Split(chirpProfanity.Body, " ")
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
			chirpClean.UserID = chirpProfanity.UserID
			chirp, err := apiCfg.databaseQuery.CreateChirp(req.Context(), database.CreateChirpParams{Body: chirpClean.Body, UserID: chirpClean.UserID})
			chirpClean.ID = chirp.ID
			if err != nil {
				error := fmt.Sprintf("failed: %v", err)
				respondWithError(w, 500, error)
				return
			}
			respondWithJSON(w, 201, chirpClean)
			return
		} else {
			chirpClean.Body = chirpProfanity.Body
			chirpClean.UserID = chirpProfanity.UserID
			chirp, err := apiCfg.databaseQuery.CreateChirp(req.Context(), database.CreateChirpParams{Body: chirpClean.Body, UserID: chirpClean.UserID})
			chirpClean.ID = chirp.ID
			if err != nil {
				error := fmt.Sprintf("failed: %v", err)
				respondWithError(w, 500, error)
				return
			}
			respondWithJSON(w, 201, chirpClean)
			return
		}
	})

	// Admin Endpoints --------------------------------------------------------------------------------------------------------------------------------

	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["GET"], endPoints["admin"], "/metrics"), func(w http.ResponseWriter, req *http.Request) {
		metricHTML := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %v times!</p></body></html>", apiCfg.metrics())
		w.WriteHeader(200)
		io.WriteString(w, metricHTML)
	})

	router.Mux.HandleFunc(fmt.Sprintf("%s %s%s", headerMethod["POST"], endPoints["admin"], "/reset"), func(w http.ResponseWriter, req *http.Request) {
		if apiCfg.platform != "dev" {
			respondWithError(w, 403, "something went wrong permission rights")
			return
		}
		apiCfg.resetMetric()
		err := apiCfg.databaseQuery.DeleteAllUsers(req.Context())
		if err != nil {
			respondWithError(w, 500, "something went wrong with database")
			return
		}
		w.WriteHeader(200)
		io.WriteString(w, "OK")
	})
}
