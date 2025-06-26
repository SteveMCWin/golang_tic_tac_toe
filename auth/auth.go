package auth

import (
    "os"
    "log"
    "strconv"
    "net/http"
    "crypto/rand"
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
    err := godotenv.Load()  // loads data like client ids and secrets from the .env file in the project root dir
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

    gothic.Store = store    // required for goth to work but not used because I couldn't really figure it out ig

	goth.UseProviders(  // specify with which platform you want the users to be able to log in
		google.New(googleClientID, googleClientSecret, googleCallbackURL, "email", "profile"),  // the string params at the end represent the scope of data you want access to
        github.New(githubClientID, githubClientSecret, githubCallbackURL, "user"),
    )

}

func SignInWithProvider(c *gin.Context) {

    provider := c.Param("provider") // we need to pass the provider into the gin context url manually for some reason, I forgor why :skull::skull:
    q := c.Request.URL.Query()
    q.Add("provider", provider)
    c.Request.URL.RawQuery = q.Encode()

    gothic.BeginAuthHandler(c.Writer, c.Request)    // this redirects to the providers log in page
}

func CallbackHandler(c *gin.Context) {  // this gets called once the user logs in with a provider
	provider := c.Param("provider") // need to pass this again to gin because it's used in the completeion of the user auth
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

    g_user, err := gothic.CompleteUserAuth(c.Writer, c.Request) // returns the user data
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

    if g_user.Email == "" { // if the email cannot be obtained, don't log the user in
        log.Printf("The user email couldn't be fetched from %s, please try another login method", provider)
        return
    }

    // the tokens are used for security and stuff hihi
    sessionToken := generateToken(32)
    csrfToken := generateToken(32)

    // the magic numbers will all be replaced in the call to StoreUser so no worries
    usr := users.User{-1, g_user.NickName, g_user.Email, g_user.AvatarURL, sessionToken, csrfToken, provider, 0, 0, 800}
    err = usr.StoreUser()

    if err != nil {
        log.Println(err)
    }

    c.SetCookie("user_id", strconv.Itoa(usr.Id), COOKIE_MAX_AGE, "/", "localhost", true, true)  // storing cookies in the browser
    c.SetCookie("session_token", sessionToken, COOKIE_MAX_AGE, "/", "localhost", true, true)    // this is used to identify the user next time they wisid the web app
    c.SetCookie("csrf_token", csrfToken, COOKIE_MAX_AGE, "/", "localhost", true, true)

	c.Redirect(http.StatusTemporaryRedirect, "/profile")
}

func ServeProfile(c *gin.Context) { // displays users profile page

    this_user, err := users.LoadUserData(c)
    if err != nil {
        log.Println("Couldn't load user, error: ", err)
        c.Redirect(http.StatusTemporaryRedirect, "/")
		return
    }

    c.HTML(http.StatusOK, "profile.html", this_user)
}

func LogoutHandler(c *gin.Context) {    // logging out erases the cookies that are used for remembering the user and telling
    c.SetCookie("user_id", "", -1, "/", "localhost", true, true)    // the cookies are erased by setting maxAge to a negative number, telling the browser they expired
    c.SetCookie("session_token", "", -1, "/", "localhost", true, true)
    c.SetCookie("csrf_token", "", -1, "/", "localhost", true, true)
    gothic.Logout(c.Writer, c.Request)  // also notify gothic we logged out ig
    c.Redirect(http.StatusTemporaryRedirect, "/")
}

func generateToken(length int) string { // returns a base64 encoded string
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        log.Fatalf("Failed to generate token: %v", err)
    }
    return base64.URLEncoding.EncodeToString(bytes)
}

