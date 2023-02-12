package example

import (
	"context"

	versionmanager "github.com/ophum/go-optimistic-locker/version_manager"
)

type PetsPresenter struct {
	versionManager versionmanager.VersionManager
}

func NewPetsPresenter(versionManager versionmanager.VersionManager) *PetsPresenter {
	return &PetsPresenter{versionManager}
}

func (p *PetsPresenter) PetResponse(ctx context.Context, pet *Pet) (*ResponsePet, error) {
	version, err := p.versionManager.Get(ctx, MakePetsKey(pet.ID))
	if err != nil {
		return nil, err
	}
	return &ResponsePet{
		Data:    pet,
		Version: version,
	}, nil
}

func (p *PetsPresenter) PetsResponse(ctx context.Context, pets []*Pet) ([]*ResponsePet, error) {
	res := make([]*ResponsePet, 0, len(pets))
	for _, pet := range pets {
		version, err := p.versionManager.Get(ctx, MakePetsKey(pet.ID))
		if err != nil {
			return nil, err
		}
		res = append(res, &ResponsePet{
			Data:    pet,
			Version: version,
		})
	}
	return res, nil
}
