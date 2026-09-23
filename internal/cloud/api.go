package cloud

import (
    "net/http"
)

func validateAPIKey(r *http.Request) bool {
    key := r.Header.Get("X-API-Key")
    if key == "" {
        return false
    }

    stored, ok := projStore.GetSetting("agent_api_key")
    if !ok {
        return false
    }

    return key == stored
}

func RequireAPIKey(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !validateAPIKey(r) {
            http.Error(w, "invalid api key", 401)
            return
        }
        next(w, r)
    }
}
