package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echooptimisticlocker "github.com/ophum/go-optimistic-locker/frameworks/echo"
	"github.com/ophum/go-optimistic-locker/internal/example"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	service := example.NewPetsService()

	e.GET("/pets", func(c echo.Context) error {
		pets, err := service.List(c.Request().Context())
		if err != nil {
			return err
		}
		res, err := example.PetsResponse(pets)
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
		res, err := example.PetResponse(pet)
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
		res, err := example.PetResponse(pet)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, res)
	})

	petsIDParser := example.IDParser("/pets/:id")
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
		res, err := example.PetResponse(pet)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, res)
	}, echooptimisticlocker.PreconditionCheck(func(r *http.Request) (any, error) {
		id, err := petsIDParser(r)
		if err != nil {
			return nil, err
		}
		return service.Get(r.Context(), id)
	}))

	e.Logger.Fatal(e.Start(":1234"))
}
