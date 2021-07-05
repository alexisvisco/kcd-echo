package kcdecho

import (
	"context"
	"github.com/alexisvisco/kcd"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Setup() {
	kcd.Config.StringsExtractors = append(kcd.Config.StringsExtractors, EchoPathExtractor{})
	kcd.Config.ValueExtractors = append(kcd.Config.ValueExtractors, EchoContextExtractor{})
}

type EchoPathExtractor struct{}

func (g EchoPathExtractor) Extract(req *http.Request, res http.ResponseWriter, valueOfTag string) ([]string, error) {
	params := req.Context().Value("echo-ctx")

	if params != nil {
		p, ok := params.(echo.Context)
		if ok {
			name := p.Param(valueOfTag)

			if name != "" {
				return []string{name}, nil
			}
		}
	}

	return nil, nil
}

func (g EchoPathExtractor) Tag() string {
	return "path"
}

type EchoContextExtractor struct {}

func (e EchoContextExtractor) Extract(req *http.Request, res http.ResponseWriter, valueOfTag string) (interface{}, error) {
	iectx := req.Context().Value("echo-ctx")
	if iectx == nil {
		return nil, nil
	}

	ectx, ok := iectx.(echo.Context)
	if !ok {
		return nil, nil
	}

	return ectx.Get(valueOfTag), nil
}

func (e EchoContextExtractor) Tag() string {
	return "echoctx"
}

func Handler(h interface{}, defaultStatusCode int) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "echo-ctx", ctx))
		kcd.Handler(h, defaultStatusCode)(ctx.Response(), req)
		return nil
	}
}
