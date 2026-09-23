package api

import (
    "encoding/json"
    "net/http"

    "schemagit/internal/cloud/store"
)

func TokenCreateHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid session token", 401)
        return
    }

    newToken, err := store.CreateToken(email)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "token": newToken,
    })
}

func TokenListHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid session token", 401)
        return
    }

    tokens, _ := store.ListTokens(email)
    json.NewEncoder(w).Encode(tokens)
}

func TokenDeleteHandler(w http.ResponseWriter, req *http.Request) {
    var body struct {
        Token string `json:"token"`
    }
    if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
        http.Error(w, "invalid request", 400)
        return
    }

    if err := store.DeleteToken(body.Token); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "status": "deleted",
    })
}
