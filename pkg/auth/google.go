package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleAuth is the Google authentication setup
type GoogleAuth struct {
	config        *oauth2.Config
	errorURL      string
	mu            sync.Mutex
	states        map[string]*state
	loginCallback LoginFunc
	pruneInterval time.Duration
}

// GoogleAuthConfig represents the configuration for GoogleAuth
type GoogleAuthConfig struct {
	ClientID      string
	ClientSecret  string
	CallbackURL   string
	ErrorURL      string
	Scopes        []string
	LoginCallback LoginFunc
}

// Userinfo contains the information we get from Google
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

// LoginFunc is called when a user has successfully logged in
type LoginFunc func(w http.ResponseWriter, r *http.Request, u *Userinfo)

// state tracks the state of outstanding requests
type state struct {
	stateString string
	created     time.Time
	remoteAddr  string
	referer     string
}

const (
	// If the user does not complete the login within expireState we
	// will forget their state.
	expireStateAfter     = 5 * time.Minute
	defaultPruneInterval = expireStateAfter / 2
	defaultErrorURL      = "/error"
)

var (
	defaultScopes = []string{"profile", "email", "openid"}
)

// New creates a new GoogleAuth instance
func New(authConfig *GoogleAuthConfig) *GoogleAuth {
	a := &GoogleAuth{
		config: &oauth2.Config{
			ClientID:     authConfig.ClientID,
			ClientSecret: authConfig.ClientSecret,
			RedirectURL:  authConfig.CallbackURL,
			Endpoint:     google.Endpoint,
			Scopes:       authConfig.Scopes,
		},
		errorURL:      authConfig.ErrorURL,
		states:        make(map[string]*state),
		pruneInterval: defaultPruneInterval,
	}

	if a.errorURL == "" {
		a.errorURL = defaultErrorURL
	}

	if a.config.Scopes == nil {
		a.config.Scopes = defaultScopes
	}

	if a.loginCallback == nil {
		a.loginCallback = defaultLoginCallback
	}

	// Start the expiry goroutine
	go a.expireOutstandingStates()

	return a
}

// expireOutstandingStates periodically goes through the state table
// and removes any old states that may be lingering.
func (a *GoogleAuth) expireOutstandingStates() {
	for {
		time.Sleep(a.pruneInterval)
		now := time.Now()
		count := 0
		a.mu.Lock()
		for k := range a.states {
			if now.Sub(a.states[k].created) > expireStateAfter {
				log.Printf("Pruned state for remoteAddr='%s' referer='%s'", a.states[k].remoteAddr, a.states[k].referer)
				delete(a.states, k)
				count++
			}
		}
		a.mu.Unlock()

		if count > 0 {
			log.Printf("Pruned %d old Google login states", count)
		}
	}
}

// NewAuthState generates a new state string and stores it for future comparison
func (a *GoogleAuth) NewAuthState(r *http.Request) string {
	authState := &state{
		stateString: createRandomStateString(),
		created:     time.Now(),
		referer:     r.Header.Get("Referer"),
		remoteAddr:  r.RemoteAddr,
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.states[authState.stateString] = authState
	return authState.stateString
}

// GetAndDelete looks up an auth state and deletes it if found.  If
// the auth state is found it returns true, otherwise it returns
// false.
func (a *GoogleAuth) GetAndDelete(state string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, ok := a.states[state]
	if ok {
		delete(a.states, state)
		return true
	}
	return false
}

// GoogleLoginHandler handles the redirect to Google
func (a *GoogleAuth) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := a.NewAuthState(r)
	url := a.config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallbackHandler handles the callback from Google
func (a *GoogleAuth) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")

	if !a.GetAndDelete(state) {
		log.Printf("Invalid oauth state, possibly timeout. remoteAddr='%s'", r.RemoteAddr)
		http.Redirect(w, r, a.errorURL, http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := a.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed: %v", err)
		http.Redirect(w, r, a.errorURL, http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("Getting userinfo failed: %v", err)
		http.Redirect(w, r, a.errorURL, http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	u := &Userinfo{}
	err = json.Unmarshal(contents, u)
	if err != nil {
		log.Printf("Unable to unmarshal userinfo: %v", err)
		http.Redirect(w, r, a.errorURL, http.StatusTemporaryRedirect)
		return
	}

	a.loginCallback(w, r, u)
}

func defaultLoginCallback(w http.ResponseWriter, r *http.Request, u *Userinfo) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
<html>
  <head>
  </head>
  <body>
    <h1>Default login callback not set in application</h1>
    But here is the data we got back from Google:

   <p>
   %+v
  </body>
</html>
`, u)
}
