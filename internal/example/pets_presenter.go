package example

import (
	optimisticlocker "github.com/ophum/go-optimistic-locker"
)

func PetResponse(pet *Pet) (*ResponsePet, error) {
	etag, err := optimisticlocker.GenerateEtag(pet)
	if err != nil {
		return nil, err
	}
	return &ResponsePet{
		Data: pet,
		Etag: etag,
	}, nil
}

func PetsResponse(pets []*Pet) ([]*ResponsePet, error) {
	res := make([]*ResponsePet, 0, len(pets))
	for _, pet := range pets {
		etag, err := optimisticlocker.GenerateEtag(pet)
		if err != nil {
			return nil, err
		}
		res = append(res, &ResponsePet{
			Data: pet,
			Etag: etag,
		})
	}
	return res, nil
}
