package api

import (
    "encoding/json"
    "net/http"
    "schemagit/internal/cloud/store"
)

func BillingGetHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", 401)
        return
    }

    plan, _ := store.GetBillingPlan(email)
    json.NewEncoder(w).Encode(plan)
}

func BillingSetHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", 401)
        return
    }

    var body struct {
        Plan string `json:"plan"`
    }
    json.NewDecoder(req.Body).Decode(&body)

    err = store.SetBillingPlan(email, body.Plan)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "status": "ok",
    })
}
