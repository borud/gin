# Simple Google Login

Import the following package

    "github.com/borud/gin/pkg/auth"

This is a work in progress.  

I got tired of baroque libraries that pull inn all manner of unwanted
cruft so I made a bare bones library to log in using Google.

This assumes you have your credentials in the environment, so make
sure you have the following environment variables set.

    GOOGLE_CLIENT_ID=<your client id>
	GOOGLE_CLIENT_SECRET=<your client secret>

You can create your client credentials at:
https://console.developers.google.com/apis/credentials



