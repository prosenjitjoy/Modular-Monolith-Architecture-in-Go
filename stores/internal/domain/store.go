package domain

import (
	"errors"
)

var (
	ErrStoreNameIsBlank               = errors.New("store name cannot be blank")
	ErrStoreLocationIsBlank           = errors.New("store location cannot be blank")
	ErrStoreIsAlreadyParticipating    = errors.New("store is already participating")
	ErrStoreIsAlreadyNotParticipating = errors.New("store is already not participating")
)

type Store struct {
	ID            string
	Name          string
	Location      string
	Participating bool
}

func CreateStore(id, name, location string) (*Store, error) {
	if name == "" {
		return nil, ErrStoreNameIsBlank
	}

	if location == "" {
		return nil, ErrStoreLocationIsBlank
	}

	store := &Store{
		ID:       id,
		Name:     name,
		Location: location,
	}

	return store, nil
}

func (s *Store) EnableParticipation() error {
	if s.Participating {
		return ErrStoreIsAlreadyParticipating
	}

	s.Participating = true

	return nil
}

func (s *Store) DisableParticipation() error {
	if !s.Participating {
		return ErrStoreIsAlreadyNotParticipating
	}

	s.Participating = false

	return nil
}
