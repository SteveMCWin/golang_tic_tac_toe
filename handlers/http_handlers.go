package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"tic_tac_toe.fun/defs"
	"tic_tac_toe.fun/game"
	"tic_tac_toe.fun/models"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var SessionManager *scs.SessionManager

var Domain string

// used by gin to load template funcs
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"rangeN": func(n int) []int {
			out := make([]int, n)
			for i := range n {
				out[i] = i
			}
			return out
		},
	}
}

func SetUpRouter(db *models.DataBase, lb *models.LeaderBoard, hub *game.Hub) http.Handler {

	err := godotenv.Load() // loads data like client secrets from the .env file in the project root dir
	if err != nil {
		log.Fatal(".env file failed to load!")
	}

	domain := os.Getenv("DOMAIN")
	csrf_key := os.Getenv("CSRF_KEY")

	if domain == "" || csrf_key == "" {
		log.Fatal("Missing domain or csrf_key")
	}

	Domain = domain

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
	store.MaxAge(defs.COOKIE_MAX_AGE)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = true

	gothic.Store = store // required for goth to work but not used because I couldn't really figure it out ig

	goth.UseProviders( // specify with which platform you want the users to be able to log in
		google.New(googleClientID, googleClientSecret, googleCallbackURL, "email", "profile"), // the string params at the end represent the scope of data you want access to
		github.New(githubClientID, githubClientSecret, githubCallbackURL, "user"),
	)

	SessionManager = scs.New()
	SessionManager.Lifetime = time.Hour * 24 * 30
	SessionManager.Store = sqlite3store.New(db.Data)
	SessionManager.Cookie.Persist = true
	SessionManager.Cookie.Secure = true

	if gin.Mode() == gin.TestMode {
		log.Println("WARNING: TEST MODE")
		SessionManager.Cookie.Secure = false
		SessionManager.Cookie.SameSite = http.SameSiteDefaultMode
		SessionManager.Cookie.Name = "test_session"
	}

	router := gin.Default()

	router.GET("/", HandleGetHome())
	router.GET("/error-page", HandleGetErrorPage())
	router.GET("/about", HandleGetAbout())
	router.GET("/leaderboard", HandleGetLeaderboard(lb))
	router.GET("/hub", ServeHub())
	router.GET("/play", ServePlay())
	router.GET("/profile/:user_id", HandleGetProfile(db))
	router.GET("/logout/:provider/", LogoutHandler())
	router.GET("/auth/:provider", SignInWithProvider())
	router.GET("/auth/:provider/callback/", CallbackHandler(db))
	router.GET("/ws", HandleWebsocketConnection(hub, db))

	router.SetFuncMap(templateFuncs())
	router.LoadHTMLGlob("templates/*") // loads all templates from the templates directory

	handler := SessionManager.LoadAndSave(router)

	if gin.Mode() != gin.TestMode {
		handler = csrf.Protect(
			[]byte(csrf_key),
			csrf.Secure(true),
		)(handler)
	}

	return handler
}

func GetUserId(c *gin.Context) int {
	if SessionManager.Exists(c.Request.Context(), "user_id") == false {
		return defs.NO_USER_ID
	}
	return SessionManager.GetInt(c.Request.Context(), "user_id")
}

func HandleGetHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		user_id := GetUserId(c)
		c.HTML(http.StatusOK, "index.html", gin.H{ "user_id": user_id }) 
	}
}

func HandleGetAbout() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "https://github.com/SteveMCWin/golang_tic_tac_toe/blob/master/readme.md")
	}
}

func HandleGetErrorPage() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Error Encountered :<") 
	}
}

func HandleGetLeaderboard(lb *models.LeaderBoard) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "leaderboard.html", gin.H{
			"TopPlayers": lb.TopPlayers,
		})
	}
}

func SignInWithProvider() func(c *gin.Context) {
	return func(c *gin.Context) {
		provider := c.Param("provider") // we need to pass the provider into the gin context url manually for some reason, I forgor why :skull::skull:
		q := c.Request.URL.Query()
		q.Add("provider", provider)
		c.Request.URL.RawQuery = q.Encode()

		gothic.BeginAuthHandler(c.Writer, c.Request) // this redirects to the providers log in page
	}
}

func CallbackHandler(db *models.DataBase) func(c *gin.Context) { // this gets called once the user logs in with a provider
	return func(c *gin.Context) {
		provider := c.Param("provider") // need to pass this again to gin because it's used in the completeion of the user auth
		q := c.Request.URL.Query()
		q.Add("provider", provider)
		c.Request.URL.RawQuery = q.Encode()

		g_user, err := gothic.CompleteUserAuth(c.Writer, c.Request) // returns the user data
		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if g_user.Email == "" { // if the email cannot be obtained, don't log the user in
			log.Printf("The user email couldn't be fetched from %s, please try another login method", provider)
			c.Redirect(http.StatusPermanentRedirect, "/error-page")
			return
		}

		usr := &models.User{
			UserName:  g_user.NickName,
			Email:     g_user.Email,
			AvatarURL: g_user.AvatarURL,
			Provider:  provider,
		}

		err = db.StoreUser(usr)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusPermanentRedirect, "/error-page")
			return
		}

		SessionManager.Put(c.Request.Context(), "user_id", usr.Id)

		c.Redirect(http.StatusTemporaryRedirect, "/profile/"+strconv.Itoa(usr.Id))
	}
}

func HandleGetProfile(db *models.DataBase) func(c *gin.Context) { // displays users profile page
	return func(c *gin.Context) {

		requesting_user_id := GetUserId(c)
		if requesting_user_id == defs.NO_USER_ID {
			// c.Redirect(http.StatusPermanentRedirect, "/user/login")
			c.Redirect(http.StatusPermanentRedirect, "/") // TODO: add a login page to redirect to
			return
		}

		user_id_param := c.Param("user_id")
		user_id, err := strconv.Atoi(user_id_param)
		if err != nil {
			log.Println("Couldn't read user id from url", err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		// TODO: limit info seen by other users

		this_user, err := db.ReadUser(user_id)
		if err != nil {
			log.Println("Couldn't load user, error: ", err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		c.HTML(http.StatusOK, "profile.html", this_user)
	}
}

func LogoutHandler() func(c *gin.Context) { // logging out erases the cookies that are used for remembering the user and telling
	return func(c *gin.Context) {
		// NOTE: not logging out with goth since it makes scs panic
		gothic.Logout(c.Writer, c.Request)
		SessionManager.Destroy(c.Request.Context())
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

func HandleWebsocketConnection(hub *game.Hub, db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		user_id := GetUserId(c)
		usr, err := db.ReadUser(user_id)
		if err != nil {
			usr = &models.User{ UserName: "Guest", Elo: defs.STARTING_ELO }
		}

		err = game.ConnectPlayerToSocket(hub, usr, c)

		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusPermanentRedirect, "/error-page")
			return
		}
	}
}

func ServePlay() func(c *gin.Context) {
	return func(c *gin.Context) {
		game_mode := c.Query("game_mode")
		if game_mode == "" {    // safeguard
			game_mode = "0"
		}

		c.HTML(http.StatusOK, "board.html", gin.H{  // this sets up the websocket in html and then the html redirects to /ws which calls MakePlayer
			"game_mode": game_mode,
		})
	}
}

func ServeHub() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "hub.html", gin.H{})
	}
}

