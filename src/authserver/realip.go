package authserver

import "net/http"

func GetRemote(r *http.Request) string {
	remote := r.Header.Get("X-Forwarded-For")
	if len(remote) == 0 {
		remote = r.RemoteAddr
	}

	return remote
}
