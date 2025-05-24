package auth

import (
    "os"
    "log"
    "strconv"
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

    g_user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    sessionToken := generateToken(32)
    csrfToken := generateToken(32)

    usr := users.User{-1, g_user.Name, g_user.Email, sessionToken, csrfToken, provider}
    err = usr.AddUser()

    if err != nil {
        panic(err)
    }

    c.SetCookie("user_id", strconv.Itoa(usr.Id), 86400 * 30, "/", "localhost", true, true)
    c.SetCookie("session_token", sessionToken, 86400 * 30, "/", "localhost", true, true)
    c.SetCookie("csrf_token", csrfToken, 86400 * 30, "/", "localhost", true, true)

	c.Redirect(http.StatusTemporaryRedirect, "/profile")
}

func ProfilePageHandler(c *gin.Context) {
    tmpl, err := template.ParseFiles("templates/profile.html")
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    this_user, err := users.LoadUserData(c)
    if err != nil {
        log.Println("Couldn't load user, error: ", err)
        c.Redirect(http.StatusTemporaryRedirect, "/")
    }

    err = tmpl.Execute(c.Writer, this_user)
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
}

func LogoutHandler(c *gin.Context) {
    // gotta erase the cookies here
    c.SetCookie("user_id", "", -1, "/", "localhost", true, true)
    c.SetCookie("session_token", "", -1, "/", "localhost", true, true)
    c.SetCookie("csrf_token", "", -1, "/", "localhost", true, true)
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

