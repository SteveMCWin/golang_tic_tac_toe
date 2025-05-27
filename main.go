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
	
	r.GET("/", home)
	r.GET("/auth/:provider", auth.SignInWithProvider)
	r.GET("/auth/:provider/callback/", auth.CallbackHandler)
    r.GET("/logout/:provider/", auth.LogoutHandler)
	r.GET("/profile", auth.ProfilePageHandler)
    r.GET("/hub", game.EnterHub)
    r.GET("/play", game.Play)
    r.GET("/ws", func (c *gin.Context) { game.MakePlayer(hub, c) })

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



