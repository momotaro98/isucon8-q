package main

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	app      *newrelic.Application
	nrclient = &http.Client{
		Transport: newrelic.NewRoundTripper(nil),
		Timeout:   time.Duration(10) * time.Second,
	}
)

func init() {
	var err error
	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName("ISUCON8Q"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		panic(err)
	}
}

// Middleware to create/end NewRelic transaction
//func nrt(inner http.Handler) http.Handler {
//	mw := func(w http.ResponseWriter, r *http.Request) {
//		txn := app.StartTransaction(r.URL.Path)
//		defer txn.End()
//
//		r = newrelic.RequestWithTransactionContext(r, txn)
//
//		txn.SetWebRequestHTTP(r)
//		w = txn.SetWebResponse(w)
//		inner.ServeHTTP(w, r)
//	}
//	return http.HandlerFunc(mw)
//}

// echo version
func nrt(name string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()

			txn := app.StartTransaction(r.URL.Path)
			defer txn.End()

			r = newrelic.RequestWithTransactionContext(r, txn)

			txn.SetWebRequestHTTP(r)
			w := c.Response().Writer
			w = txn.SetWebResponse(w)

			err := next(c)
			return err
		}
	}
}
