// Package handlers is used for setting up the router and handling http requests.
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
)

var sessionManager *scs.SessionManager

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
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"div": func(a, b float64) float64 {
			return a / b
		},
		"float64": func(a int) float64 {
			return float64(a)
		},
		"substr": func(s string, start, length int) string {
			if start < 0 || start >= len(s) {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
	}
}

// SetUpRouter does what it says.
// It expects the following .env variables:
// - CSRF_KEY
// - SESSION_KEY
// - GOOGLE_CLIENT_ID
// - GOOGLE_CLIENT_SECRET
// - GOOGLE_CALLBACK_URL
// - GITHUB_CLIENT_ID
// - GITHUB_CLIENT_SECRET
// - GITHUB_CALLBACK_URL
// - DOMAIN
// It sets up the goth package for OAuth 2.0, and also wraps the default gin handler with csrf protection and session management.
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

	// Domain = domain

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

	sessionManager = scs.New()
	sessionManager.Lifetime = time.Hour * 24 * 30
	sessionManager.Store = sqlite3store.New(db.Data)
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = true

	if gin.Mode() == gin.TestMode {
		log.Println("WARNING: TEST MODE")
		sessionManager.Cookie.Secure = false
		sessionManager.Cookie.SameSite = http.SameSiteDefaultMode
		sessionManager.Cookie.Name = "test_session"
	}

	router := gin.Default()

	router.Static("/css", "./css")
	router.Static("/js", "./js")

	router.GET("/", HandleGetHome())
	router.GET("/error-page", HandleGetErrorPage())
	router.GET("/about", HandleGetAbout())
	router.GET("/leaderboard", HandleGetLeaderboard(lb))
	router.GET("/hub", ServeHub())
	router.GET("/play", ServePlay())
	router.GET("/search", HandleGetSearch(db))
	router.GET("/profile/:user_id", MiddlewareNoCache(), HandleGetProfile(db))
	router.GET("/profile/:user_id/delete", MiddlewareNoCache(), HandleGetDeleteProfile())
	router.DELETE("/profile/:user_id", HandleDeleteProfile(db))
	router.GET("/profile/:user_id/games_played", HandleGetUserGames(db))
	router.GET("/replay/:record_id", HandleGetGameReplay(db))
	router.POST("/replay", HandlePostGameReplay())
	router.GET("/logout/:provider/", LogoutHandler())
	router.GET("/auth/:provider", SignInWithProvider())
	router.GET("/auth/:provider/callback/", CallbackHandler(db))
	router.GET("/ws", HandleWebsocketConnection(hub, db))

	router.SetFuncMap(templateFuncs())
	router.LoadHTMLGlob("templates/*") // loads all templates from the templates directory

	handler := sessionManager.LoadAndSave(router)

	if gin.Mode() != gin.TestMode {
		handler = csrf.Protect(
			[]byte(csrf_key),
			csrf.Secure(true),
		)(handler)
	}

	return handler
}

// GetUserId uses the scs session manager to return the user's id from their browser.
// If there is no user id stored (meaning the user is not logged in), the NO_USER_ID value from the defs package is returned.
func GetUserId(c *gin.Context) int {
	if sessionManager.Exists(c.Request.Context(), "user_id") == false {
		return defs.NO_USER_ID
	}
	return sessionManager.GetInt(c.Request.Context(), "user_id")
}

// HandleGetHome displays the landing page of the website.
// Calls GetUserId to determine if it should display the login buttons or the profile button.
func HandleGetHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		user_id := GetUserId(c)
		c.HTML(http.StatusOK, "index.html", gin.H{"user_id": user_id})
	}
}

// HandleGetAbout page just redirects the user to the github readme of the project.
func HandleGetAbout() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "https://github.com/SteveMCWin/golang_tic_tac_toe/blob/master/readme.md")
	}
}

// HandleGetErrorPage is used for unexpected errors.
func HandleGetErrorPage() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Error Encountered :<")
	}
}

// HandleGetLeaderboard displays the top players(by elo).
func HandleGetLeaderboard(lb *models.LeaderBoard) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "leaderboard.html", gin.H{
			"TopPlayers": lb.TopPlayers,
		})
	}
}

// SignInWithProvider begins the gothic authentication process, after which the user will be redirected to the CallbackHandler.
func SignInWithProvider() func(c *gin.Context) {
	return func(c *gin.Context) {
		provider := c.Param("provider") // we need to pass the provider into the gin context url manually for some reason, I forgor why :skull::skull:
		q := c.Request.URL.Query()
		q.Add("provider", provider)
		c.Request.URL.RawQuery = q.Encode()

		user_id := GetUserId(c)
		if user_id != defs.NO_USER_ID {
			c.Redirect(http.StatusPermanentRedirect, "/profile/"+strconv.Itoa(user_id))
			return
		}

		gothic.BeginAuthHandler(c.Writer, c.Request) // this redirects to the providers log in page
	}
}

// CallbackHandler finalizes the OAuth 2.0 process started by the gothic package.
// If the user authentication is completed successfully and the user's email is retrieved, the StoreUser() function is called,
// the session manager stores the user_id in the browser and the user is taken to their profile page.
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

		sessionManager.Put(c.Request.Context(), "user_id", usr.Id)

		c.Redirect(http.StatusTemporaryRedirect, "/profile/"+strconv.Itoa(usr.Id))
	}
}

// HandleGetProfile displays the profile page of the user whose id is in the url.
func HandleGetProfile(db *models.DataBase) func(c *gin.Context) { // displays users profile page
	return func(c *gin.Context) {

		requesting_user_id := GetUserId(c)
		// if requesting_user_id == defs.NO_USER_ID {
		// 	c.Redirect(http.StatusPermanentRedirect, "/") // TODO: add a login page to redirect to
		// 	return
		// }

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

		if requesting_user_id != user_id {
			this_user.Email = ""
		}

		c.HTML(http.StatusOK, "profile.html", this_user)
	}
}

// HandleGetUserGames reads the GameRecords which include the user (whose id was passed in the url)
// as one of the players and displays them as an html page.
func HandleGetUserGames(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		user_id_param := c.Param("user_id")
		user_id, err := strconv.Atoi(user_id_param)
		if err != nil {
			log.Println("Error calling atoi:", err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		records, err := db.ReadGameRecordsForUser(user_id)
		if err != nil {
			log.Println("Error reading game records for user:", err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		winner_ids := make([]int, len(records))

		for i, rec := range records {
			winner_ids[i] = models.GetGameRecordWinner(rec)
		}
		// WARNING: Update the html to handle -1 in winner_ids
		c.HTML(http.StatusOK, "view_games_history.html", gin.H{"records": records, "winner_ids": winner_ids, "user_id": user_id})
	}
}

// HandleGetGameReplay
// TODO
func HandleGetGameReplay(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		requesting_user_id := GetUserId(c)
		if requesting_user_id == defs.NO_USER_ID {
			log.Println("You must be looged in to view replay")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		record_id_param := c.Param("record_id")
		record_id, err := strconv.Atoi(record_id_param)
		if err != nil {
			log.Println("Error getting the game id from url", err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		rec, err := db.ReadGameRecord(record_id)
		if err != nil {
			log.Println("Error reading game record:", err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		game.InitGameReplay(requesting_user_id, rec)

		c.HTML(http.StatusOK, "replay.html", gin.H{csrf.TemplateTag: csrf.TemplateField(c.Request)})
	}
}

func HandlePostGameReplay() func(c *gin.Context) {
	return func(c *gin.Context) {
		// rec_id_param := c.Param("record_id")
		// rec_id, err := strconv.Atoi(rec_id_param)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{})
		// }

		requesting_user_id := GetUserId(c)
		if requesting_user_id == defs.NO_USER_ID {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}

		msg := struct {
			Move int `json:"msg"`
		}{}

		if err := c.BindJSON(&msg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
		}

		var new_board_state []byte
		var err error
		switch msg.Move {
		case defs.NEXT_MOVE:
			new_board_state, err = game.ReplayNextMove(requesting_user_id)
			if err != nil {
				log.Println("Error calling game.ReplayNextMove:", err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		case defs.PREV_MOVE:
			new_board_state, err = game.ReplayPrevMove(requesting_user_id)
			if err != nil {
				log.Println("Error calling game.ReplayPrevMove:", err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		default:
			log.Println("Recieved neither defs.NEXT_MOVE nor defs.PREV_MOVE")
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.String(http.StatusOK, string(new_board_state)) // NOTE: the new_board_state is already json encoded, so just put it in the response body as is
	}
}

// LogoutHandler deletes the session data and calls the gotich logout function.
func LogoutHandler() func(c *gin.Context) { // logging out erases the cookies that are used for remembering the user and telling
	return func(c *gin.Context) {
		// NOTE: not logging out with goth since it makes scs panic
		gothic.Logout(c.Writer, c.Request)
		sessionManager.Destroy(c.Request.Context())
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

// HandleWebsocketConnection connects the player to the socket by calling game.ConnectPlayerToSocketa.
func HandleWebsocketConnection(hub *game.Hub, db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		user_id := GetUserId(c)
		usr, err := db.ReadUser(user_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusPermanentRedirect, "/error-page")
			return
		}

		err = game.ConnectPlayerToSocket(hub, usr, c)

		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusPermanentRedirect, "/error-page")
			return
		}
	}
}

// ServePlay displays the board the players will play on and initializes the websocket connection through the html it serves.
func ServePlay() func(c *gin.Context) {
	return func(c *gin.Context) {
		game_mode := c.Query("game_mode")
		if game_mode == "" { // safeguard
			game_mode = "0"
		}

		c.HTML(http.StatusOK, "board.html", gin.H{ // this sets up the websocket in html and then the html redirects to /ws which calls MakePlayer
			"game_mode": game_mode,
		})
	}
}

func HandleGetSearch(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		query := c.Query("name")
		if query == "" {
			c.HTML(http.StatusOK, "search_users.html", gin.H{})
		} else {
			requesting_user_id := GetUserId(c)
			results, err := db.SearchForUsers(query, requesting_user_id)
			if err != nil {
				log.Println(err)
				c.Redirect(http.StatusPermanentRedirect, "/error-page")
				return
			}
			c.JSON(200, results)
		}
	}
}

// ServeHub displays the pre-game hub where the user selects which game mode they want to play.
func ServeHub() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "hub.html", gin.H{})
	}
}

func HandleGetDeleteProfile() func(c *gin.Context) {
	return func(c *gin.Context) {

		requesting_user_id := GetUserId(c)
		if requesting_user_id == defs.NO_USER_ID {
			log.Println("You cannot delete an account if you are not logged into that account")
			c.Redirect(http.StatusPermanentRedirect, "/")
			return
		}

		user_id_param := c.Param("user_id")
		user_id, err := strconv.Atoi(user_id_param)
		if err != nil {
			log.Println("Error converting user_id_param to int")
			c.Redirect(http.StatusPermanentRedirect, "/error-page")
			return
		}

		c.HTML(http.StatusOK, "delete_account.html", gin.H{"user_id": user_id, csrf.TemplateTag: csrf.TemplateField(c.Request)})
	}
}

func HandleDeleteProfile(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		log.Println("Called DELETE profile")
		requesting_user_id := GetUserId(c)
		if requesting_user_id == defs.NO_USER_ID {
			log.Println("You cannot delete an account if you are not logged into that account")
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		user_id_param := c.Param("user_id")
		user_id, err := strconv.Atoi(user_id_param)
		if err != nil {
			log.Println("Error converting user_id_param to int")
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		if user_id != requesting_user_id {
			log.Println("You cannot delete another users account")
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		err = db.DeleteUser(user_id)
		if err != nil {
			log.Println("Error deleting user profile:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		gothic.Logout(c.Writer, c.Request)
		sessionManager.Destroy(c.Request.Context())

		c.JSON(http.StatusOK, gin.H{})
	}
}

// MIDDLEWARE

// MiddlewareNoCache is used to wipe the cached pages. Useful when you don't want to allow a user to log out and then go back to their profile with the back arrow.
func MiddlewareNoCache() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-store")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Next()
	}
}
