//go:build unit

package application_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/application"
	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/domain"
)

func Test_teamService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("should create a team successfully", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{
			InsertOneFn: func(ctx context.Context, team *domain.Team) (*domain.Team, error) {
				assert.NotNil(t, ctx)
				assert.Equal(t, "brazil", team.ID)
				assert.Equal(t, "Brazil", team.Name)
				assert.Equal(t, "A", team.Group)

				return domain.NewTeam(team.ID, team.Name, team.Group)
			},
		}

		service := application.NewTeamService(repository)

		input := application.NewTeamInput("brazil", "Brazil", "A")
		output, err := service.Create(ctx, input)
		assert.NoError(t, err)

		expectedTeam, err := domain.NewTeam("brazil", "Brazil", "A")
		assert.NoError(t, err)

		expected := application.NewTeamOutput(expectedTeam)
		assert.EqualValues(t, expected, output)

		assert.Equal(t, 1, repository.InsertOneInvokedCount)
	})

	t.Run("should return new domain team error", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{}
		service := application.NewTeamService(repository)

		input := application.NewTeamInput("", "Brazil", "A")
		output, err := service.Create(ctx, input)
		assert.Nil(t, output)
		assert.ErrorIs(t, err, domain.ErrInvalidArgument)
		assert.EqualError(t, err, "error on create domain team: error creating team without 'id': invalid argument")

		assert.Equal(t, 0, repository.InsertOneInvokedCount)
	})

	t.Run("should return repository error", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{
			InsertOneFn: func(ctx context.Context, team *domain.Team) (*domain.Team, error) {
				return nil, fmt.Errorf("repository error")
			},
		}

		service := application.NewTeamService(repository)

		input := application.NewTeamInput("brazil", "Brazil", "A")
		output, err := service.Create(ctx, input)
		assert.Nil(t, output)
		assert.EqualError(t, err, "error on insert domain team: repository error")

		assert.Equal(t, 1, repository.InsertOneInvokedCount)
	})
}

func Test_teamService_Get(t *testing.T) {
	ctx := context.Background()

	t.Run("should get a team successfully", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{
			FindOneFn: func(ctx context.Context, ID string) (*domain.Team, error) {
				assert.NotNil(t, ctx)
				assert.Equal(t, "brazil", ID)

				return domain.NewTeam("brazil", "Brazil", "A")
			},
		}

		service := application.NewTeamService(repository)

		output, err := service.Get(ctx, "brazil")
		assert.NoError(t, err)

		expectedTeam, err := domain.NewTeam("brazil", "Brazil", "A")
		assert.NoError(t, err)

		expected := application.NewTeamOutput(expectedTeam)
		assert.EqualValues(t, expected, output)

		assert.Equal(t, 1, repository.FindOneInvokedCount)
	})

	t.Run("should return required id error", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{}
		service := application.NewTeamService(repository)

		output, err := service.Get(ctx, "")
		assert.Nil(t, output)
		assert.ErrorIs(t, err, domain.ErrInvalidArgument)
		assert.EqualError(t, err, "team id is required: invalid argument")

		assert.Equal(t, 0, repository.FindOneInvokedCount)
	})

	t.Run("should return repository error", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{
			FindOneFn: func(ctx context.Context, ID string) (*domain.Team, error) {
				return nil, fmt.Errorf("repository error")
			},
		}

		service := application.NewTeamService(repository)

		output, err := service.Get(ctx, "brazil")
		assert.Nil(t, output)
		assert.EqualError(t, err, "error on find team by id: repository error")

		assert.Equal(t, 1, repository.FindOneInvokedCount)
	})
}

func Test_teamService_Remove(t *testing.T) {
	ctx := context.Background()

	t.Run("should remove a team successfully", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{
			DeleteOneFn: func(ctx context.Context, ID string) error {
				assert.NotNil(t, ctx)
				assert.Equal(t, "brazil", ID)

				return nil
			},
		}

		service := application.NewTeamService(repository)

		err := service.Remove(ctx, "brazil")
		assert.NoError(t, err)

		assert.Equal(t, 1, repository.DeleteOneInvokedCount)
	})

	t.Run("should return required id error", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{}
		service := application.NewTeamService(repository)

		err := service.Remove(ctx, "")
		assert.ErrorIs(t, err, domain.ErrInvalidArgument)
		assert.EqualError(t, err, "team id is required: invalid argument")

		assert.Equal(t, 0, repository.DeleteOneInvokedCount)
	})

	t.Run("should return repository error", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{
			DeleteOneFn: func(ctx context.Context, ID string) error {
				return fmt.Errorf("repository error")
			},
		}

		service := application.NewTeamService(repository)

		err := service.Remove(ctx, "brazil")
		assert.EqualError(t, err, "error on delete team by id: repository error")

		assert.Equal(t, 1, repository.DeleteOneInvokedCount)
	})
}

func Test_teamService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("should get all teams successfully", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{
			FindFn: func(ctx context.Context) ([]*domain.Team, error) {
				assert.NotNil(t, ctx)

				team, err := domain.NewTeam("brazil", "Brazil", "A")
				assert.NoError(t, err)

				return []*domain.Team{team}, nil
			},
		}

		service := application.NewTeamService(repository)

		teams, err := service.GetAll(ctx)
		assert.NoError(t, err)

		expectedTeam, err := domain.NewTeam("brazil", "Brazil", "A")
		assert.NoError(t, err)

		expectedTeams := []*application.TeamOutput{
			application.NewTeamOutput(expectedTeam),
		}
		assert.EqualValues(t, expectedTeams, teams)

		assert.Equal(t, 1, repository.FindInvokedCount)
	})

	t.Run("should return repository error", func(t *testing.T) {
		repository := &application.TeamRepositoryMock{
			FindFn: func(ctx context.Context) ([]*domain.Team, error) {
				return nil, fmt.Errorf("repository error")
			},
		}

		service := application.NewTeamService(repository)

		teams, err := service.GetAll(ctx)
		assert.Nil(t, teams)
		assert.EqualError(t, err, "error on find all teams: repository error")

		assert.Equal(t, 1, repository.FindInvokedCount)
	})
}
