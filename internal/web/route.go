package web

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(*Context) error
type Middleware func(HandlerFunc) HandlerFunc

func (r *Router) toHttpHandler(h HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		for i := range middlewares {
			h = middlewares[i](h)
		}

		ctx := Context{
			Response: rw,
			Request:  req,
			Data:     make(map[string]any),
			renderer: r.renderer,
		}

		ctx.Data["Path"] = req.URL.Path

		if err := h(&ctx); err != nil {
			// Assumption is made that error should be consumed by some middleware
			panic(fmt.Sprintf("http handler: %v", err))
		}
	}
}

func WrapHandlerFunc(h http.HandlerFunc) HandlerFunc {
	return func(ctx *Context) error {
		h.ServeHTTP(ctx.Response, ctx.Request)
		return nil
	}
}
