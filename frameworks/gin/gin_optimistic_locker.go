package ginoptimisticlocker

import (
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	optimisticlocker "github.com/ophum/go-optimistic-locker"
)

func PreconditionCheck(get optimisticlocker.ResourceGetter) gin.HandlerFunc {
	nextHandler, wrapper := adapter.New()
	return wrapper(optimisticlocker.PreconditionCheck(get)(nextHandler))
}
