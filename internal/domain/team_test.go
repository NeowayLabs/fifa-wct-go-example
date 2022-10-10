//go:build unit

package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/domain"
)

func TestNewTeam(t *testing.T) {
	t.Run("should create a team successfully", func(t *testing.T) {
		team, err := domain.NewTeam("id", "name", "group")
		assert.NoError(t, err)

		expectedTeam := &domain.Team{ID: "id", Name: "name", Group: "group"}
		assert.EqualValues(t, expectedTeam, team)
	})

	t.Run("should return invalid argument error without id", func(t *testing.T) {
		team, err := domain.NewTeam("", "name", "group")
		assert.Nil(t, team)
		assert.ErrorIs(t, err, domain.ErrInvalidArgument)
		assert.EqualError(t, err, "error creating team without 'id': invalid argument")
	})

	t.Run("should return invalid argument error without name", func(t *testing.T) {
		team, err := domain.NewTeam("id", "", "group")
		assert.Nil(t, team)
		assert.ErrorIs(t, err, domain.ErrInvalidArgument)
		assert.EqualError(t, err, "error creating team without 'name': invalid argument")
	})

	t.Run("should return invalid argument error without group", func(t *testing.T) {
		team, err := domain.NewTeam("id", "name", "")
		assert.Nil(t, team)
		assert.ErrorIs(t, err, domain.ErrInvalidArgument)
		assert.EqualError(t, err, "error creating team without 'group': invalid argument")
	})

}
