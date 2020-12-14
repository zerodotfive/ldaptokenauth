package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/zerodotfive/ldaptokenauth/src/browser"
	"github.com/zerodotfive/ldaptokenauth/src/callbackserver"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage:\n# %s <auth_url>\n", os.Args[0])
		os.Exit(1)
	}

	httpSync := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv, port := callbackserver.Run(httpSync, ctx)
	if err := browser.Run(fmt.Sprintf("%s?callback=%s", os.Args[1], port)); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		srv.Shutdown(ctx)
		os.Exit(4)
	}

	httpSync.Wait()
	srv.Shutdown(ctx)
	cancel()
}
