package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginoptimisticlocker "github.com/ophum/go-optimistic-locker/frameworks/gin"
	"github.com/ophum/go-optimistic-locker/internal/example"
)

func h(f func(ctx *gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := f(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
func main() {
	r := gin.Default()

	service := example.NewPetsService()
	r.GET("/pets", func(ctx *gin.Context) {
		pets, err := service.List(ctx.Request.Context())
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		res, err := example.PetsResponse(pets)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, res)

	})

	r.GET("/pets/:id", func(ctx *gin.Context) {
		var reqUri example.Pet
		if err := ctx.ShouldBindUri(&reqUri); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		pet, err := service.Get(ctx.Request.Context(), reqUri.ID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		res, err := example.PetResponse(pet)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, res)
	})

	r.POST("/pets", func(ctx *gin.Context) {
		var req example.Pet
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		pet, err := service.Create(ctx.Request.Context(), &req)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		res, err := example.PetResponse(pet)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusCreated, res)
	})

	petsIDParser := example.IDParser("/pets/:id")
	r.PUT("/pets/:id",
		ginoptimisticlocker.PreconditionCheck(func(r *http.Request) (any, error) {
			id, err := petsIDParser(r)
			if err != nil {
				return nil, err
			}
			return service.Get(r.Context(), id)
		}),
		func(ctx *gin.Context) {
			var reqUri example.Pet
			if err := ctx.ShouldBindUri(&reqUri); err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
			var req example.Pet
			if err := ctx.ShouldBindJSON(&req); err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}

			if reqUri.ID != req.ID {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "id of uri and body is mismatch",
				})
				return
			}

			pet, err := service.Update(ctx.Request.Context(), &req)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			res, err := example.PetResponse(pet)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			ctx.JSON(http.StatusOK, res)
		},
	)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
