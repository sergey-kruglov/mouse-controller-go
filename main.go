package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vgo/robotgo"
	socketio "github.com/googollee/go-socket.io"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	// r.Static("/", "public")

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "app:click", func(s socketio.Conn, msg string) {
		robotgo.Click()
	})

	server.OnEvent("/", "app:touchmove", func(s socketio.Conn, msg []int) {
		robotgo.MoveRelative(msg[0], msg[1])
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("error:", e)
	})

	go server.Serve()
	defer server.Close()

	fs := http.FileServer(http.Dir("./public"))

	r.GET("/", gin.WrapH(fs))
	r.Any("/socket.io/", gin.WrapH(server))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
