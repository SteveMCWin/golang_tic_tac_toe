package main

import (
	// "fmt"
	"html/template"
	"log"
	"net/http"
	"os"
    "crypto/rand"
    "encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
    "github.com/gorilla/sessions"
)

var user goth.User
var SessionToken string
var CSRFToken string

// type Login struct {
//     user goth.User
//     SessionToken string
//     CSRFToken string
// }

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load!")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

    store := sessions.NewCookieStore([]byte("random_str1found50m3wh3r3"))
    store.MaxAge(90)
    store.Options.Path = "/"
    store.Options.HttpOnly = true
    store.Options.Secure = true

    gothic.Store = store

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL, "email", "profile"),
    )

	r.LoadHTMLGlob("templates/*")
	
	r.GET("/", home)
	r.GET("/auth/:provider", signInWithProvider)
	r.GET("/auth/:provider/callback/", callbackHandler)
    r.GET("/logout/:provider/", logoutHandler)
	r.GET("/success", Success)

	r.RunTLS(":5000", "./testdata/server.pem", "./testdata/server.key")
}

func home(c *gin.Context) {
    // cookie, err := c.Cookie("session_token")
    // if err != nil {
    //     fmt.Println("NO COOKIEEES")
    // } else {
    //     fmt.Println("I FOUND DA COOKIEEEEE YIPPPIIIEE")
    //     fmt.Println("%v", cookie)
    // }

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(c.Writer, gin.H{})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func signInWithProvider(c *gin.Context) {
    var err error
    user, err = gothic.CompleteUserAuth(c.Writer, c.Request)
    if err != nil {
        provider := c.Param("provider")
        q := c.Request.URL.Query()
        q.Add("provider", provider)
        c.Request.URL.RawQuery = q.Encode()

        gothic.BeginAuthHandler(c.Writer, c.Request)
    } else {
        c.Redirect(http.StatusTemporaryRedirect, "/success")
    }
}

func callbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

    var err error
	user, err = gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    sessionToken := generateToken(32)
    csrfToken := generateToken(32)

    c.SetCookie("session_token", sessionToken, 60, "/", "localhost", true, true)
    c.SetCookie("csrf_token", csrfToken, 60, "/", "localhost", true, false)
	
	c.Redirect(http.StatusTemporaryRedirect, "/success")
}

func Success(c *gin.Context) {
    tmpl, err := template.ParseFiles("templates/success.html")
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(c.Writer, user)
    if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
}

func logoutHandler(c *gin.Context) {
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

