package main

import (
	// "fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
    "github.com/gorilla/sessions"
)

var user goth.User

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

    store := sessions.NewCookieStore([]byte("random_string1found50m3wh3r3"))
    store.MaxAge(86400 * 30)
    store.Options.Path = "/"
    store.Options.HttpOnly = true
    store.Options.Secure = true

    gothic.Store = store

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL))

	r.LoadHTMLGlob("templates/*")
	
	r.GET("/", home)
	r.GET("/auth/:provider", signInWithProvider)
	r.GET("/auth/:provider/callback/", callbackHandler)
	r.GET("/success", Success)

	r.RunTLS(":5000", "./testdata/server.pem", "./testdata/server.key")
}

func home(c *gin.Context) {
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
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
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
