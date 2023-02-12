package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	optimisticlocker "github.com/ophum/go-optimistic-locker"
	ginoptimisticlocker "github.com/ophum/go-optimistic-locker/frameworks/gin"
	"github.com/ophum/go-optimistic-locker/internal/example"
	"github.com/ophum/go-optimistic-locker/version_manager/inmemory"
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

	versionManager := inmemory.NewInmemoryStore()
	locker := optimisticlocker.NewLocker(versionManager)

	service := example.NewPetsService()
	presenter := example.NewPetsPresenter(versionManager)

	r.GET("/pets", func(ctx *gin.Context) {
		pets, err := service.List(ctx.Request.Context())
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		res, err := presenter.PetsResponse(ctx.Request.Context(), pets)
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

		res, err := presenter.PetResponse(ctx.Request.Context(), pet)
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

		if _, err := versionManager.Create(ctx.Request.Context(), example.MakePetsKey(pet.ID)); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		res, err := presenter.PetResponse(ctx.Request.Context(), pet)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusCreated, res)
	})

	r.PUT("/pets/:id",
		ginoptimisticlocker.PreconditionCheck(locker, example.KeyParser("/pets/:id")),
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

			if _, err := versionManager.Update(ctx.Request.Context(), example.MakePetsKey(pet.ID)); err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			res, err := presenter.PetResponse(ctx.Request.Context(), pet)
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
