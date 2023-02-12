package example

import (
	"errors"
	"sync"
)

type PetStore struct {
	pets          []*Pet
	autoIncrement uint
	mu            *sync.RWMutex
}

func NewPetStore() *PetStore {
	return &PetStore{
		pets:          []*Pet{},
		autoIncrement: 0,
		mu:            &sync.RWMutex{},
	}
}

func (s *PetStore) List() []*Pet {
	ret := make([]*Pet, 0, len(s.pets))
	for _, p := range s.pets {
		ret = append(ret, p.Copy())
	}
	return ret
}

func (s *PetStore) Get(id uint) (*Pet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, pet := range s.pets {
		if pet.ID == id {
			return pet.Copy(), nil
		}
	}
	return nil, errors.New("not found")
}

func (s *PetStore) Add(pet *Pet) uint {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.autoIncrement++
	id := s.autoIncrement
	pet = pet.Copy()
	pet.ID = id
	s.pets = append(s.pets, pet)
	return id
}

func (s *PetStore) Set(id uint, pet *Pet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.pets {
		if p.ID == id {
			s.pets[i] = pet.Copy()
			return nil
		}
	}
	return errors.New("not found")
}

func (s *PetStore) Delete(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.pets {
		if p.ID == id {
			s.pets = append(s.pets[:i], s.pets[i+1:]...)
			return
		}
	}
}
