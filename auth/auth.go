package auth

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)
//
// const (
//     key = "somerandomstring"
//     MaxAge = 86400 * 30
//     IsProd = false
// )
//
// func NewAuth() {
//     err := godotenv.Load()
//     if err != nil {
//         log.Fatal("Error loading .env file")
//     }
//
    // googleClientId := os.GetEnv("GOOGLE_CLIENT_ID")
    // googleClientSecret := os.GetEnv("GOOGLE_CLIENT_SECRET")
//
    // store := sessions.NewCookieStore([]byte(key))
    // store.MaxAge(MaxAge)
    //
    // store.Options.Path = "/"
    // store.Options.HttpOnly = true
    // store.Options.Secure = IsProd
    //
    // gothic.Store = store
//
//     goth.UseProviders(googleClientId, googleClientSecret, "http://localhost:8080/callback-gl")
// }
