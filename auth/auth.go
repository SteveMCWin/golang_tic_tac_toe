package auth

import (
    "log"
    "os"
    "net/http"
    "crypto/rand"
	"html/template"
    "encoding/base64"
    "tic_tac_toe.fun/users"

	"github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
    "github.com/gorilla/sessions"
)

func init() {
    err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load!")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")
    sessionKey := os.Getenv("SESSION_KEY")

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" || sessionKey == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

    log.Println("making store")
    store := sessions.NewCookieStore([]byte(sessionKey))
    store.MaxAge(90)
    store.Options.Path = "/"
    store.Options.HttpOnly = true
    store.Options.Secure = true

    gothic.Store = store

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL, "email", "profile"),
    )

}

func SignInWithProvider(c *gin.Context) {
    provider := c.Param("provider")
    q := c.Request.URL.Query()
    q.Add("provider", provider)
    c.Request.URL.RawQuery = q.Encode()

    gothic.BeginAuthHandler(c.Writer, c.Request)
}

func CallbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

    user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    sessionToken := generateToken(32)
    csrfToken := generateToken(32)

    users.AddActiveUser(user.Name, user.Email, sessionToken, csrfToken, provider)

    c.SetCookie("username", user.Name, 90, "/", "localhost", true, true)
    c.SetCookie("session_token", sessionToken, 90, "/", "localhost", true, true)
    c.SetCookie("csrf_token", csrfToken, 90, "/", "localhost", true, true)

    // log.Println("storing stuff in session")
    // // gothic.StoreInSession("session_token", sessionToken, c.Request, c.Writer)
    // // gothic.StoreInSession("csrf_token", csrfToken, c.Request, c.Writer)
    // // gothic.StoreInSession("username", user.Name, c.Request, c.Writer)
    // session, err := gothic.Store.Get(c.Request, gothic.SessionName)
    // if err != nil {
    //     c.AbortWithError(http.StatusInternalServerError, err)
    // }
    // session.Values["username"] = user.Name
    // session.Values["csrf_token"] = csrfToken
    // session.Values["session_token"] = sessionToken
    // log.Println("STORE SESSION VALUES:", session.Values)
    // err = gothic.Store.Save(c.Request, c.Writer, session)
    //
    // if err != nil {
    //     log.Fatal("COULDN'T SAVE SESSION")
    // }

	c.Redirect(http.StatusTemporaryRedirect, "/profile")
}

func ProfilePageHandler(c *gin.Context) {
    tmpl, err := template.ParseFiles("templates/profile.html")
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    user, err := c.Cookie("username")
    if err != nil {
        log.Println("BRUHHHHHHHHHHHHHHH(username)")
    } else {
        log.Println("username", user)
    }

    csrf, err := c.Cookie("csrf_token")
    if err != nil {
        log.Println("BRUHHHHHHHHHHHHHHH(csrf)")
    } else {
        log.Println("csrf_token", csrf)
    }

    sess, err := c.Cookie("session_token")
    if err != nil {
        log.Println("BRUHHHHHHHHHHHHHHH(sess)")
    } else {
        log.Println("session_token", sess)
    }

    this_user, found := users.Users[user]
    if found != true {
        log.Fatal("man...")
    } else {
        log.Println("this_user:", this_user)
    }
    err = tmpl.Execute(c.Writer, this_user)
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
}

func LogoutHandler(c *gin.Context) {
    gothic.Logout(c.Writer, c.Request)
    c.Redirect(http.StatusTemporaryRedirect, "/")
}

func generateToken(length int) string {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        log.Fatalf("Failed to generate token: %v", err)
    }
    return base64.URLEncoding.EncodeToString(bytes)
}

