# Simple Google Login

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/borud/gin/pkg/auth)

**This is a work in progress.**

I got tired of baroque libraries that pull inn all manner of unwanted
cruft so I made a bare bones library to log in using Google. This code
takes care of only the Google login.  

Import the following package

    "github.com/borud/gin/pkg/auth"

## Abbreviated example

Here is an abbreviated example of how to use this module.  The
`loginCallback` is where you would create the session etc, but since
people use different libraries for handling sessions we leave this
part up to you.

    package main
    
    import (
    	"fmt"
    	"net/http"
    	"os"
    
    	"github.com/borud/gin/pkg/auth"
    )
    
    func main() {
    	googleAuth := auth.New(&auth.GoogleAuthConfig{
    		ClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
    		ClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
    		CallbackURL:   "http://localhost:3000/google/callback",
    		LoginCallback: loginCallback,
    	})
    
    	http.HandleFunc("/google/login", googleAuth.GoogleLoginHandler)
    	http.HandleFunc("/google/callback", googleAuth.GoogleCallbackHandler)
    	http.ListenAndServe(":3000", nil)
    }
    
    func loginCallback(w http.ResponseWriter, r *http.Request, userinfo *auth.Userinfo) {
    	html := `Hello %s<p><img src="%s"`
    	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    	fmt.Fprintf(w, html, userinfo.Name, userinfo.Picture)
    }
        

Please see the [example.go](example.go) file for a more complete example.


## Session management

You have to do your own session creation and management.  You can do
this by implementing your own callback of the following type:

    type LoginFunc func(w http.ResponseWriter, r *http.Request, u *Userinfo)
	
The `Userinfo` struct that is provided contains the following fields:

    type Userinfo struct {
    	ID            string `json:"id"`
    	Email         string `json:"email"`
    	VerifiedEmail bool   `json:"verified_email"`
    	Name          string `json:"name"`
    	GivenName     string `json:"given_name"`
    	FamilyName    string `json:"family_name"`
    	Picture       string `json:"picture"`
    	Locale        string `json:"locale"`
    	HostDomain    string `json:"hd"`
    }


## Providing credentials
This assumes you have your credentials in the environment, so make
sure you have the following environment variables set.

    GOOGLE_CLIENT_ID=<your client id>
	GOOGLE_CLIENT_SECRET=<your client secret>

You can create your client credentials at:
https://console.developers.google.com/apis/credentials

