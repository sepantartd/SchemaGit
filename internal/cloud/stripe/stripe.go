package stripe

import (
    "os"
    "github.com/stripe/stripe-go/v74"
    "github.com/stripe/stripe-go/v74/checkout/session"
)

func Init() {
    stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

func CreateCheckoutSession(email string) (*stripe.CheckoutSession, error) {
    params := &stripe.CheckoutSessionParams{
        Mode: stripe.String("payment"),
        SuccessURL: stripe.String("https://your-domain.com/billing?success=true"),
        CancelURL:  stripe.String("https://your-domain.com/billing?cancel=true"),
        CustomerEmail: stripe.String(email),
        LineItems: []*stripe.CheckoutSessionLineItemParams{
            {
                Price:    stripe.String(os.Getenv("STRIPE_PRICE_PRO")),
                Quantity: stripe.Int64(1),
            },
        },
    }

    return session.New(params)
}
