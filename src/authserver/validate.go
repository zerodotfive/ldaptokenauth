package authserver

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/zerodotfive/ldaptokenauth/src/ldapclient"
	"github.com/zerodotfive/ldaptokenauth/src/tokenjwt"
)

func (s *AuthServer) Validate(w http.ResponseWriter, r *http.Request) {
	tokenSigned, err := url.QueryUnescape(r.Header.Get("Authorization"))
	if err != nil {
		log.Printf("Fail: %s, remote %s", err, GetRemote(r))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	cookie, err := r.Cookie("Authorization")
	if err == nil {
		data, err := url.QueryUnescape(cookie.Value)
		if err == nil && len(data) > 0 {
			tokenSigned = data
		}
	}

	if len(tokenSigned) == 0 {
		err := fmt.Errorf(http.StatusText(http.StatusForbidden))
		log.Printf("Fail: %s, remote %s", err, GetRemote(r))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	username, expires, err := tokenjwt.Parse(tokenSigned, s.Secret)
	if err != nil {
		log.Printf("Fail: %s, remote %s", err, GetRemote(r))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if expires.Before(time.Now()) {
		err := fmt.Errorf("Session expired")
		log.Printf("Fail: %s, remote %s", err, GetRemote(r))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	user, err := ldapclient.Search(s.LdapServer, s.BindDn, s.BindPw, s.BaseDN, s.GroupFilter, username)
	if user != nil && err == nil {
		w.Header().Set("username", username)
		log.Printf("Success: valid token for '%s', remote %s", username, GetRemote(r))
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		return
	}

	err = fmt.Errorf(http.StatusText(http.StatusForbidden))
	log.Printf("Fail: %s, remote %s", err, GetRemote(r))
	http.Error(w, err.Error(), http.StatusForbidden)
}
