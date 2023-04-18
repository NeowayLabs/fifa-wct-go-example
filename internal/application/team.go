package application

import (
	"context"
	"fmt"

	"github.com/NeowayLabs/fifa-wct-go-example/internal/domain"
)

type TeamInput struct {
	ID    string
	Name  string
	Group string
}

func NewTeamInput(ID, Name, Group string) *TeamInput {
	return &TeamInput{ID: ID, Name: Name, Group: Group}
}

type TeamOutput struct {
	ID    string
	Name  string
	Group string
}

func NewTeamOutput(team *domain.Team) *TeamOutput {
	return &TeamOutput{ID: team.ID, Name: team.Name, Group: team.Group}
}

type TeamService interface {
	Create(ctx context.Context, input *TeamInput) (*TeamOutput, error)
	Get(ctx context.Context, ID string) (*TeamOutput, error)
	Remove(ctx context.Context, ID string) error
	GetAll(ctx context.Context) ([]*TeamOutput, error)
}

type TeamRepository interface {
	InsertOne(ctx context.Context, team *domain.Team) (*domain.Team, error)
	FindOne(ctx context.Context, ID string) (*domain.Team, error)
	DeleteOne(ctx context.Context, ID string) error
	Find(ctx context.Context) ([]*domain.Team, error)
}

type teamService struct {
	teamRepository TeamRepository
}

func NewTeamService(teamRepository TeamRepository) *teamService {
	return &teamService{
		teamRepository: teamRepository,
	}
}

func (ts *teamService) Create(ctx context.Context, input *TeamInput) (*TeamOutput, error) {
	team, err := domain.NewTeam(input.ID, input.Name, input.Group)
	if err != nil {
		return nil, fmt.Errorf("error on create domain team: %w", err)
	}

	newTem, err := ts.teamRepository.InsertOne(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("error on insert domain team: %w", err)
	}

	return NewTeamOutput(newTem), nil
}

func (ts *teamService) Get(ctx context.Context, ID string) (*TeamOutput, error) {
	if ID == "" {
		return nil, fmt.Errorf("team id is required: %w", domain.ErrInvalidArgument)
	}

	team, err := ts.teamRepository.FindOne(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("error on find team by id: %w", err)
	}

	return NewTeamOutput(team), nil
}

func (ts *teamService) Remove(ctx context.Context, ID string) error {
	if ID == "" {
		return fmt.Errorf("team id is required: %w", domain.ErrInvalidArgument)
	}

	err := ts.teamRepository.DeleteOne(ctx, ID)
	if err != nil {
		return fmt.Errorf("error on delete team by id: %w", err)
	}

	return nil
}

func (ts *teamService) GetAll(ctx context.Context) ([]*TeamOutput, error) {
	teams, err := ts.teamRepository.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("error on find all teams: %w", err)
	}

	outputs := make([]*TeamOutput, 0)
	for _, team := range teams {
		outputs = append(outputs, NewTeamOutput(team))
	}

	return outputs, nil
}
