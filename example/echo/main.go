package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	optimisticlocker "github.com/ophum/go-optimistic-locker"
	echooptimisticlocker "github.com/ophum/go-optimistic-locker/frameworks/echo"
	"github.com/ophum/go-optimistic-locker/internal/example"
	"github.com/ophum/go-optimistic-locker/version_manager/inmemory"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	versionManager := inmemory.NewInmemoryStore()
	locker := optimisticlocker.NewLocker(versionManager)

	service := example.NewPetsService()
	presenter := example.NewPetsPresenter(versionManager)

	e.GET("/pets", func(c echo.Context) error {
		pets, err := service.List(c.Request().Context())
		if err != nil {
			return err
		}
		res, err := presenter.PetsResponse(c.Request().Context(), pets)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, res)
	})

	e.GET("/pets/:id", func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return err
		}
		pet, err := service.Get(c.Request().Context(), uint(id))
		if err != nil {
			return err
		}
		res, err := presenter.PetResponse(c.Request().Context(), pet)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, res)
	})

	e.POST("/pets", func(c echo.Context) error {
		var req example.Pet
		if err := c.Bind(&req); err != nil {
			return err
		}

		pet, err := service.Create(c.Request().Context(), &req)
		if err != nil {
			return err
		}
		if _, err := versionManager.Create(c.Request().Context(), example.MakePetsKey(pet.ID)); err != nil {
			return err
		}
		res, err := presenter.PetResponse(c.Request().Context(), pet)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, res)
	})

	e.PUT("/pets/:id", func(c echo.Context) error {
		var req example.Pet
		if err := c.Bind(&req); err != nil {
			return err
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return err
		}
		if req.ID != uint(id) {
			return errors.New("id of uri and body is mismatch")
		}

		pet, err := service.Update(c.Request().Context(), &req)
		if err != nil {
			return err
		}
		if _, err := versionManager.Update(c.Request().Context(), example.MakePetsKey(pet.ID)); err != nil {
			return err
		}
		res, err := presenter.PetResponse(c.Request().Context(), pet)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, res)
	}, echooptimisticlocker.PreconditionCheck(locker, example.KeyParser("/pets/:id")))

	e.Logger.Fatal(e.Start(":1234"))
}
