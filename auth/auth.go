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
	"github.com/markbates/goth/providers/github"
    "github.com/gorilla/sessions"
)

const COOKIE_MAX_AGE = 86400 * 30

func init() {
    err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load!")
	}

    sessionKey := os.Getenv("SESSION_KEY")

    if sessionKey == "" {
        log.Fatal("SESSION KEY missing")
    }

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleCallbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

	if googleClientID == "" || googleClientSecret == "" || googleCallbackURL == "" {
		log.Fatal("Environment variables (GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_CALLBACK_URL) are required")
	}

    githubClientID := os.Getenv("GITHUB_CLIENT_ID")
    githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
    githubCallbackURL := os.Getenv("GITHUB_CALLBACK_URL")

	if githubClientID == "" || githubClientSecret == "" || githubCallbackURL == "" {
		log.Fatal("Environment variables (GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, GITHUB_CALLBACK_URL) are required")
	}

    store := sessions.NewCookieStore([]byte(sessionKey))
    store.MaxAge(COOKIE_MAX_AGE)
    store.Options.Path = "/"
    store.Options.HttpOnly = true
    store.Options.Secure = true

    gothic.Store = store

	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, googleCallbackURL, "email", "profile"),
        github.New(githubClientID, githubClientSecret, githubCallbackURL, "user"),
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

    if g_user.Email == "" {
        log.Printf("The user email couldn't be fetched from %s, please try another login method", provider)
    }

    sessionToken := generateToken(32)
    csrfToken := generateToken(32)

    usr := users.User{-1, g_user.NickName, g_user.Email, g_user.AvatarURL, sessionToken, csrfToken, provider, 0, 0}
    err = usr.StoreUser()

    if err != nil {
        log.Println(err)
    }

    c.SetCookie("user_id", strconv.Itoa(usr.Id), COOKIE_MAX_AGE, "/", "localhost", true, true)
    c.SetCookie("session_token", sessionToken, COOKIE_MAX_AGE, "/", "localhost", true, true)
    c.SetCookie("csrf_token", csrfToken, COOKIE_MAX_AGE, "/", "localhost", true, true)

	c.Redirect(http.StatusTemporaryRedirect, "/profile")
}

func ServeProfile(c *gin.Context) {
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

