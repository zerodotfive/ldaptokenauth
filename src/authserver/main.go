package authserver

import (
	"net/http"
	"time"
)

type AuthServer struct {
	Listen      string
	Secret      string
	LdapServer  string
	BindDn      string
	BindPw      string
	BaseDN      string
	GroupFilter string
	TTL         time.Duration
}

func (s *AuthServer) Run() {
	http.HandleFunc("/auth", s.Auth)
	http.HandleFunc("/validate", s.Validate)

	http.ListenAndServe(s.Listen, nil)
}
