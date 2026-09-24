package api

import (
    "encoding/json"
    "log"
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

    if orgs, err := store.UserOrgs(body.Email); err == nil && len(orgs) > 0 {
        log.Printf("login success email=%s org=%s", body.Email, orgs[0].ID)
        store.AddAudit(orgs[0].ID, "login", body.Email)
    }

    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
    })
}
