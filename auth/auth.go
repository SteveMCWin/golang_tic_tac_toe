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

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

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

    users.AddActiveUser(user.Name, user.Email, sessionToken, csrfToken)

    gothic.StoreInSession("username", user.Name, c.Request, c.Writer)
    gothic.StoreInSession("session_token", sessionToken, c.Request, c.Writer)
    gothic.StoreInSession("csrf_token", csrfToken, c.Request, c.Writer)

	c.Redirect(http.StatusTemporaryRedirect, "/profile")
}

func ProfilePageHandler(c *gin.Context) {
    tmpl, err := template.ParseFiles("templates/profile.html")
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    usr, _ := gothic.GetFromSession("username", c.Request)
    // if err != nil {
    //     log.Println("yeah got stuck here :<")
    //     c.AbortWithStatus(http.StatusInternalServerError)
    //     return
    // }
    err = tmpl.Execute(c.Writer, users.Users[usr])
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

