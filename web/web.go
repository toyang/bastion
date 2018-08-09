package web

import (
	"fmt"
	"github.com/novakit/nova"
	"github.com/novakit/static"
	"github.com/novakit/view"
	"github.com/yankeguo/bastion/types"
	"net/http"
	"google.golang.org/grpc/status"
)

func NewServer(opts types.WebOptions) *http.Server {
	n := nova.New()
	if opts.Dev {
		n.Env = nova.Development
	} else {
		n.Env = nova.Production
	}
	n.Error(func(c *nova.Context, err error) {
		if s, ok := status.FromError(err); ok {
			// if it's a grpc status error, extract description
			http.Error(c.Res, s.Message(), http.StatusInternalServerError)
		} else {
			// just render any error as 500 and expose the message
			http.Error(c.Res, err.Error(), http.StatusInternalServerError)
		}
	})
	// mount static module
	n.Use(static.Handler(static.Options{
		Directory: "public",
		BinFS:     !opts.Dev,
	}))
	// mount view module for json rendering only
	n.Use(view.Handler(view.Options{
		Directory: "views",
		BinFS:     !opts.Dev,
	}))
	// mount rpc module
	n.Use(rpcModule(opts))
	// mount auth module
	n.Use(authModule())
	// mount all routes
	mountRoutes(n)
	// build the http.Server
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Handler: n,
	}
}