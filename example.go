package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/borud/gin/pkg/auth"
)

const (
	listenAddr = ":3000"
)

func main() {
	googleAuth := auth.New(&auth.GoogleAuthConfig{
		ClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		CallbackURL:   "http://localhost:3000/google/callback",
		LoginCallback: loginCallback,
	})

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/error", errorHandler)
	http.HandleFunc("/google/login", googleAuth.GoogleLoginHandler)
	http.HandleFunc("/google/callback", googleAuth.GoogleCallbackHandler)
	http.ListenAndServe(listenAddr, nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
<html>
  <head>
  </head>
  <body>
    <h1>Welcome</h1>
    <a href="/google/login">Login</a>
  </body>
</html>
`)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
<html>
  <head>
  </head>
  <body>
    <h1>Error</h1>
  </body>
</html>
`)
}

func loginCallback(w http.ResponseWriter, r *http.Request, userinfo *auth.Userinfo) {
	html := `<html><head><title>Wohoo></title></head><body>Hello %s<p><img src="%s"</body></html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html, userinfo.Name, userinfo.Picture)
}
