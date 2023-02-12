package ginoptimisticlocker

import (
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	optimisticlocker "github.com/ophum/go-optimistic-locker"
)

func PreconditionCheck(l optimisticlocker.Locker, k optimisticlocker.VersionKeyParser) gin.HandlerFunc {
	nextHandler, wrapper := adapter.New()
	return wrapper(l.PreconditionCheck(k)(nextHandler))
}
