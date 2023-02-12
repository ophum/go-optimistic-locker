package echooptimisticlocker

import (
	"github.com/labstack/echo/v4"
	optimisticlocker "github.com/ophum/go-optimistic-locker"
)

func PreconditionCheck(get optimisticlocker.ResourceGetter) echo.MiddlewareFunc {
	return echo.WrapMiddleware(optimisticlocker.PreconditionCheck(get))
}
