package main

import (
	"html/template"
	"net/http"
    "tic_tac_toe.fun/auth"
    "tic_tac_toe.fun/game"

	"github.com/gin-gonic/gin"
)

func main() {

    hub := game.MakeHub()
    go hub.HandleGames()

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	
	r.GET("/", ServeHome)
    // r.GET("/about", ServeAbout)
    // r.GET("/leaderboard", ServeLeaderboard)
    r.GET("/hub", game.ServeHub)
    r.GET("/play", game.ServePlay)
	r.GET("/profile", auth.ServeProfile)
    r.GET("/logout/:provider/", auth.LogoutHandler)
	r.GET("/auth/:provider", auth.SignInWithProvider)
	r.GET("/auth/:provider/callback/", auth.CallbackHandler)
    r.GET("/ws", func (c *gin.Context) { game.MakePlayer(hub, c) })

	r.RunTLS(":5000", "./testdata/server.pem", "./testdata/server.key")
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



