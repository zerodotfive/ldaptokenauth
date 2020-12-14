package callbackserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

func Run(httpSync *sync.WaitGroup, ctx context.Context) (*http.Server, string) {
	httpSync.Add(1)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(2)
	}

	srv := &http.Server{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer httpSync.Done()
		token := r.URL.Query()["token"]
		fmt.Fprintf(w, "<html><body><h2>Successfully logged in!</h2>You may close the window.</body></html>")
		fmt.Fprintf(os.Stdout, "%s", token[0])
	})

	go func() {
		if err := srv.Serve(listener); err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(3)
		}
	}()

	return srv, strings.Replace(listener.Addr().String(), "127.0.0.1:", "", -1)
}
