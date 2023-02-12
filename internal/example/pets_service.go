package example

import (
	"context"
)

type PetsService struct {
	petStore *PetStore
}

func NewPetsService() *PetsService {
	return &PetsService{
		petStore: NewPetStore(),
	}
}

func (c *PetsService) List(ctx context.Context) ([]*Pet, error) {
	return c.petStore.List(), nil
}

func (c *PetsService) Get(ctx context.Context, id uint) (*Pet, error) {
	return c.petStore.Get(id)
}

func (c *PetsService) Create(ctx context.Context, req *Pet) (*Pet, error) {
	id := c.petStore.Add(&Pet{
		Name: req.Name,
		Age:  req.Age,
	})
	return c.petStore.Get(id)
}

func (c *PetsService) Update(ctx context.Context, req *Pet) (*Pet, error) {
	if err := c.petStore.Set(req.ID, req); err != nil {
		return nil, err
	}
	return c.petStore.Get(req.ID)
}
