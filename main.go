package main

import (
	"html/template"
	"net/http"
    "tic_tac_toe.fun/auth"
    "tic_tac_toe.fun/game"
    "tic_tac_toe.fun/users"

	"github.com/gin-gonic/gin"
)

func main() {

    users.InitDb()

    hub := game.MakeHub()
    go hub.HandleGames()

	r := gin.Default()

    r.SetFuncMap(templateFuncs())
	r.LoadHTMLGlob("templates/*")
	
	r.GET("/", ServeHome)
    // r.GET("/about", ServeAbout)
    r.GET("/leaderboard", users.ServeLeaderboard)
    r.GET("/hub", game.ServeHub)
    r.GET("/play", game.ServePlay)
	r.GET("/profile", auth.ServeProfile)
    r.GET("/logout/:provider/", auth.LogoutHandler)
	r.GET("/auth/:provider", auth.SignInWithProvider)
	r.GET("/auth/:provider/callback/", auth.CallbackHandler)
    r.GET("/ws", func (c *gin.Context) { game.MakePlayer(hub, c) })

	r.RunTLS(":5000", "./testdata/server.pem", "./testdata/server.key")
}

func templateFuncs() template.FuncMap {
    return template.FuncMap{
        "add1": func(i int) int {
            return i + 1
        },
    }
}

func ServeHome(c *gin.Context) {

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



