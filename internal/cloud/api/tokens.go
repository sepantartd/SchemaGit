package api

import (
    "encoding/json"
    "log"
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
    if orgs, err := store.UserOrgs(email); err == nil && len(orgs) > 0 {
        log.Printf("token created org=%s email=%s", orgs[0].ID, email)
        store.AddAudit(orgs[0].ID, "token.create", newToken)
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
    json.NewDecoder(req.Body).Decode(&body)

    store.DeleteToken(body.Token)

    json.NewEncoder(w).Encode(map[string]string{
        "status": "deleted",
    })
}
