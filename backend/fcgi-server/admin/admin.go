package main

import (
	"time"
	"net/http"
	"crypto/rand"
	"math/big"
	"fmt"
	"encoding/hex"

	"github.com/s2ks/fcgiserver/logger"
	"github.com/s2ks/fcgiserver"
)

type AdminPage struct {
	page *ConfigPage
}

type LoginPage struct {
	ActiveLogins int

	page *ConfigPage
}

type LoginStatus int

const (
	SessionCookieName = "session_token"
)

const (
	LOGIN_INVAL = iota
	LOGIN_NODATA
	LOGIN_TRIES
	LOGIN_OK
)

const (
	SESSION_INVAL = iota
	SESSION_NODATA
	SESSION_EXPIRED
	SESSION_OK
)

func NewSessionCookie() *http.Cookie {
	cookie := new(http.Cookie)

	token := make([]byte, 32)

	_, err := rand.Read(token)

	if err != nil {
		return nil
	}

	cookie.Name = SessionCookieName
	cookie.Value = hex.EncodeToString(token)
	cookie.Expires = time.Now().Add(6 * time.Hour)
	cookie.Path = "/"

	return cookie
}

/* no-op */
func (a *AdminPage) Execute(s *server.FcgiServer) ([]byte, error) {
	return nil, nil
}

func (a *AdminPage) ServeHTTP(w http.ResponseWriter, r *http.Request, h *server.Handle) {
	/* TODO CHECK session cookie */

	server.ServeHTTP(w, r, h)
}

func (l *LoginPage) CheckCredentials(r *http.Request) LoginStatus {
	err := r.ParseForm()
}

func (l *LoginPage) CheckSessionCookie(r *http.Cookie) SessionStatus {

}

func (l *LoginPage) Execute(s *server.FcgiServer) ([]byte, error) {

}

func (l *LoginPage) ServeHTTP(w http.ResponseWriter, r *http.Request, h *server.Handle) {
	/*
		Check the request for a session cookie
	*/

	cookie, err := r.Cookie(SessionCookieName)

	if err != nil {

		switch(l.CheckCredentials(r)) {
			case LOGIN_OK:
				/* LOGIN CREDENTIALS OK */
				/* Set session cookie */
				cookie = NewSessionCookie()

				http.Redirect(w, r, "/login", http.StatusFound)
				break
			case LOGIN_NODATA:
				/* LOGIN FORM NOT FILLED */
				/* SHOW NORMAL LOGIN PAGE */
				break
			case LOGIN_INVAL:
				/* WRONG USER/PASS */
				/* COUNT NUMBER OF TRIES */
				break
			case LOGIN_TRIES:
				/* REACHED MAX TRIES */
				/* TIMEOUT */
				break
			default:
				/* PROGRAMMING ERROR */
				break
		}
	} else {
		/* CHECK IF COOKIE VALUE MATCHES STORED VALUE FOR GIVEN REMOTE ADDRESS (r.RemoteAddr)  */
		switch(l.CheckSessionCookie(cookie)) {
			case SESSION_OK:
				/* SESSION COOKIE OK */
				http.Redirect(w, r, "/", http.StatusFound)
				break
			case SESSION_EXPIRED:
				/* SESSION EXPIRED */
				/* SHOW NORMAL LOGIN PAGE */
				break
			case SESSION_INVAL:
				/* SESSION KEY/REMOTE ADDRESS DO NOT MATCH */
				/* NORMAL LOGIN PAGE*/
				break
			case SESSION_NODATA:
				/* SESSION COOKIE DOES NOT CONTAIN EXPECTED DATA */
				/* NORMAL LOGIN PAGE */
				break
			default:
				/* PROGRAMMING ERROR */
				break
		}
	}
}

func main() {

	var conf *Config

	s, err := server.New()

	if err != nil {
		logger.Fatal(err)
	}

	server.RegisterForExec(s.Register("/login", &LoginPage{}))

	s.Register("/", &AdminPage{})
	s.Register("/add", &AdminPage{})
	s.Register("/edit", &AdminPage{})
	s.Register("/delete", &AdminPage{})

	logger.Fatal(s.Serve())
}
