package main

import (
	"html/template"
	"net/http"

	"tic_tac_toe.fun/auth"
	"tic_tac_toe.fun/game"
	"tic_tac_toe.fun/users"

	"github.com/gin-gonic/gin"
)

// used by gin to load template funcs
func templateFuncs() template.FuncMap {
    return template.FuncMap{
        "add1": func(i int) int {
            return i + 1
        },
    }
}

func main() {

    users.InitDb()

    hub := game.MakeHub()
    go hub.HandleGames()

	r := gin.Default()

    r.SetFuncMap(templateFuncs())
	r.LoadHTMLGlob("templates/*")   // loads all templates from the templates directory
	
	r.GET("/", func (c *gin.Context) { c.HTML(http.StatusOK, "index.html", gin.H{}) })
    r.GET("/about", func (c *gin.Context) { c.Redirect(http.StatusPermanentRedirect, "https://github.com/SteveMCWin/golang_tic_tac_toe/blob/master/readme.md")})
    r.GET("/leaderboard", users.ServeLeaderboard)
    r.GET("/hub", game.ServeHub)
    r.GET("/play", game.ServePlay)
	r.GET("/profile", auth.ServeProfile)
    r.GET("/logout/:provider/", auth.LogoutHandler)
	r.GET("/auth/:provider", auth.SignInWithProvider)
	r.GET("/auth/:provider/callback/", auth.CallbackHandler)
    r.GET("/ws", func (c *gin.Context) { game.MakePlayer(hub, c) })

	r.Run(":5000")
}

