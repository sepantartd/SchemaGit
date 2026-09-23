package api

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "os"
    "schemagit/internal/cloud/store"

    "github.com/stripe/stripe-go/v74"
    "github.com/stripe/stripe-go/v74/webhook"
)

func StripeWebhookHandler(w http.ResponseWriter, req *http.Request) {
    payload, _ := ioutil.ReadAll(req.Body)
    sig := req.Header.Get("Stripe-Signature")

    event, err := webhook.ConstructEvent(payload, sig, os.Getenv("STRIPE_WEBHOOK_SECRET"))
    if err != nil {
        http.Error(w, "invalid signature", 400)
        return
    }

    if event.Type == "checkout.session.completed" {
        var session stripe.CheckoutSession
        json.Unmarshal(event.Data.Raw, &session)

        email := session.CustomerEmail
        if email != "" {
            store.SetBillingPlan(email, "pro")
        }
    }

    w.WriteHeader(200)
}
