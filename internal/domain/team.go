package domain

import "fmt"

type Team struct {
	ID    string
	Name  string
	Group string
}

func NewTeam(ID, Name, Group string) (*Team, error) {
	if ID == "" {
		return nil, fmt.Errorf("error creating team without 'id': %w", ErrInvalidArgument)
	}

	if Name == "" {
		return nil, fmt.Errorf("error creating team without 'name': %w", ErrInvalidArgument)
	}

	if Group == "" {
		return nil, fmt.Errorf("error creating team without 'group': %w", ErrInvalidArgument)
	}

	return &Team{ID: ID, Name: Name, Group: Group}, nil
}
