package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/yiqing95/freechat/server/socketio"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {

	/**
		*  -------------------------------------------------------------------   +|
							静态文件服务
	     @see https://github.com/YuriyNasretdinov/social-net/blob/part1/main.go # serveStatic

		 http.HandleFunc("/static2/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		})
	*/
	logger := log.New(os.Stdout, "server: ", log.Lshortfile)
	// http.Handle("/", Notify(logger)(indexHandler))
	fs := http.FileServer(http.Dir("assets/"))
	// 对于形如 /static/css/some-style.css   的文件 对应的寻找位置是 assets/css/some-style.css
	http.Handle("/static/", Notify(logger)(http.StripPrefix("/static/", fs)))
	/**
	*  -------------------------------------------------------------------   +|
	 */

	flag.Parse()

	ioserver := socketio.InitHandler()
	http.Handle("/socket.io/", ioserver)

	// http.HandleFunc("/", serveHome)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		NewView(w).Rend("home/index", nil)
	})

	/*
		http.HandleFunc("/socketio", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "home-socketio.html")
		})
	*/

	// err := http.ListenAndServe(*addr, nil)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
