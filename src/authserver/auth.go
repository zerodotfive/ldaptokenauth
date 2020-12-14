package authserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/zerodotfive/ldaptokenauth/src/ldapclient"
	"github.com/zerodotfive/ldaptokenauth/src/tokenjwt"
)

func (s *AuthServer) Auth(w http.ResponseWriter, r *http.Request) {
	callback := r.URL.Query()["callback"]
	if len(callback) == 0 {
		err := fmt.Errorf("Auth request with no callback")
		log.Printf("Fail: %s, remote %s", err, GetRemote(r))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err := strconv.ParseUint(callback[0], 0, 16)
	if err != nil {
		log.Printf("Fail: %s, remote %s", err, GetRemote(r))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		s.authGET(w, r)
	case "POST":
		s.authPOST(w, r)
	default:
		err := fmt.Errorf("Bad request method")
		log.Printf("Fail: %s, remote %s", err, GetRemote(r))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *AuthServer) authGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(
		w,
		"<!doctype html><html><body><form action='/auth?callback=%s' method='post'><label for='username'><b>Username</b></label><input type='text' placeholder='Enter Username' name='username' required><br/><label for='password'><b>Password</b></label><input type='password' placeholder='Enter Password' name='password' required><br/><button type='submit'>Login</button></form></body></html>",
		r.URL.Query()["callback"][0],
	)
}

func (s *AuthServer) authPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, err := ldapclient.Search(s.LdapServer, s.BindDn, s.BindPw, s.BaseDN, s.GroupFilter, r.Form.Get("username"))
	if err != nil {
		log.Printf("Fail: %s, remote %s", err, GetRemote(r))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if !ldapclient.IsValidUser(s.LdapServer, user.DN, r.Form.Get("password")) {
		err := fmt.Errorf("Invalid user '%s'", r.Form.Get("username"))
		log.Printf("Fail: %s, remote %s", err, GetRemote(r))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	token, err := tokenjwt.Generate(user.GetAttributeValue("uid"), s.Secret, s.TTL)
	if err != nil {
		log.Printf("Fail: %s, remote %s", err, GetRemote(r))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "Authorization",
		Domain:  r.Header.Get("Host"),
		Path:    "/",
		Expires: time.Now().Add(s.TTL),
	}
	cookie.Value = token
	http.SetCookie(w, &cookie)
	log.Printf("Success: logged in as '%s', remote %s", r.Form.Get("username"), GetRemote(r))
	http.Redirect(w, r, fmt.Sprintf("http://127.0.0.1:%s/?token=%s", r.URL.Query()["callback"][0], token), http.StatusTemporaryRedirect)
}
