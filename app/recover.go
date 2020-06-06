package app

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/ca17/go-common/log"
)

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func ServerRecover(debug bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					if debug {
						log.Errorf("%+v", r)
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

