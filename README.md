# Simple Google Login

**This is a work in progress.**

I got tired of baroque libraries that pull inn all manner of unwanted
cruft so I made a bare bones library to log in using Google. This code
takes care of only the Google login.  

Import the following package

    "github.com/borud/gin/pkg/auth"

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

## Examples

Please see the [example.go](example.go) file for a minimal example.
