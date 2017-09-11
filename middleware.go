package main

import (
	"log"
	"net/http"
)

// @see https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81

type Adapter func(http.Handler) http.Handler

func Notify1() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("before")
			defer log.Println("after")
			h.ServeHTTP(w, r)
		})
	}
}

/**
* 中间件函数示例
 */
func Notify(logger *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("request URL " + r.URL.Path)
			logger.Println("before")
			defer logger.Println("after")
			h.ServeHTTP(w, r)
		})
	}
}

/*
func Notify(logger *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("before")
			defer logger.Println("after")
			h.ServeHTTP(w, r)
		})
	}
}
*/

/*
logger := log.New(os.Stdout, "server: ", log.Lshortfile)
http.Handle("/", Notify(logger)(indexHandler))
*/

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

/*
http.Handle("/", Adapt(indexHandler, AddHeader("Server", "Mine"),
CheckAuth(providers),
CopyMgoSession(db),
Notify(logger),
)
*/
