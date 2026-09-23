package api

import (
    "encoding/json"
    "net/http"
    "schemagit/internal/cloud/store"
    "schemagit/internal/cloud/stripe"
)

func StripeCheckoutHandler(w http.ResponseWriter, req *http.Request) {
    token := req.Header.Get("Authorization")
    email, err := store.ValidateToken(token)
    if err != nil {
        http.Error(w, "invalid token", 401)
        return
    }

    stripe.Init()
    session, err := stripe.CreateCheckoutSession(email)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "checkout_url": session.URL,
    })
}
