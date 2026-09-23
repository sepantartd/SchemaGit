package api

import (
    "encoding/json"
    "net/http"
    "schemagit/internal/cloud/store"
)

type SignupRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func SignupHandler(w http.ResponseWriter, req *http.Request) {
    var body SignupRequest
    json.NewDecoder(req.Body).Decode(&body)

    err := store.CreateUser(body.Email, body.Password)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "status": "ok",
    })
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {
    var body LoginRequest
    json.NewDecoder(req.Body).Decode(&body)

    token, err := store.Authenticate(body.Email, body.Password)
    if err != nil {
        http.Error(w, "invalid credentials", 401)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
    })
}
