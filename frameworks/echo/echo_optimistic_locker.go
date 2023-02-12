package echooptimisticlocker

import (
	"github.com/labstack/echo/v4"
	optimisticlocker "github.com/ophum/go-optimistic-locker"
)

func PreconditionCheck(locker optimisticlocker.Locker, keyParser optimisticlocker.VersionKeyParser) echo.MiddlewareFunc {
	return echo.WrapMiddleware(locker.PreconditionCheck(keyParser))
}
