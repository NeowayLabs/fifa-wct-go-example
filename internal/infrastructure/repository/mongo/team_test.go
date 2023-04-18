//go:build integration

package mongo_test

import (
	"context"
	"testing"

	"github.com/NeowayLabs/fifa-wct-go-example/internal/domain"
	"github.com/NeowayLabs/fifa-wct-go-example/internal/infrastructure/repository/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewTeamRepository(t *testing.T) {
	database := newDatabaseForTest()

	storage, err := mongo.NewTeamRepository(database, testLog)
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}

func TestTeamRepositoryInsertOne(t *testing.T) {
	ctx := context.Background()

	t.Run("should insert a team successfully", func(t *testing.T) {
		database := newDatabaseForTest()
		repository, err := mongo.NewTeamRepository(database, testLog)
		assert.NoError(t, err)

		teamToInsert, err := domain.NewTeam("id", "name", "group")
		assert.NoError(t, err)

		team, err := repository.InsertOne(ctx, teamToInsert)
		assert.NoError(t, err)
		assert.NotNil(t, team)

		result := &mongo.Team{}
		err = database.Collection(mongo.TeamCollectionName).FindOne(ctx, bson.M{"_id": team.ID}).Decode(result)
		assert.NoError(t, err)
		assert.EqualValues(t, teamToInsert, result)
	})

	t.Run("should return duplicate team error", func(t *testing.T) {
		database := newDatabaseForTest()
		repository, err := mongo.NewTeamRepository(database, testLog)
		assert.NoError(t, err)

		teamToInsert, err := domain.NewTeam("id", "name", "group")
		assert.NoError(t, err)

		team, err := repository.InsertOne(ctx, teamToInsert)
		assert.NoError(t, err)
		assert.NotNil(t, team)

		team, err = repository.InsertOne(ctx, teamToInsert)
		assert.Nil(t, team)
		assert.ErrorIs(t, err, domain.ErrDuplicateKey)
	})
}

func TestTeamRepositoryFindOne(t *testing.T) {
	ctx := context.Background()

	database := newDatabaseForTest()
	repository, err := mongo.NewTeamRepository(database, testLog)
	assert.NoError(t, err)

	database.Collection(mongo.TeamCollectionName).InsertMany(ctx, []interface{}{
		bson.M{"_id": "brazil", "name": "Brazil", "group": "A"},
		bson.M{"_id": "portugal", "name": "Portugal", "group": "B"},
	})

	t.Run("should find one team successfully", func(t *testing.T) {
		team, err := repository.FindOne(ctx, "brazil")
		assert.NoError(t, err)
		assert.NotNil(t, team)

		expectedTeam, err := domain.NewTeam("brazil", "Brazil", "A")
		assert.NoError(t, err)
		assert.EqualValues(t, expectedTeam, team)
	})

	t.Run("should return not found error", func(t *testing.T) {
		team, err := repository.FindOne(ctx, "china")
		assert.Nil(t, team)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestTeamRepositoryDeleteOne(t *testing.T) {
	ctx := context.Background()

	database := newDatabaseForTest()
	repository, err := mongo.NewTeamRepository(database, testLog)
	assert.NoError(t, err)

	teamCollection := database.Collection(mongo.TeamCollectionName)
	teamCollection.InsertMany(ctx, []interface{}{
		bson.M{"_id": "brazil", "name": "Brazil", "group": "A"},
		bson.M{"_id": "portugal", "name": "Portugal", "group": "B"},
	})

	t.Run("should delete one team successfully", func(t *testing.T) {
		err := repository.DeleteOne(ctx, "brazil")
		assert.NoError(t, err)

		count, err := teamCollection.CountDocuments(ctx, bson.M{})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		count, err = teamCollection.CountDocuments(ctx, bson.M{"_id": "portugal"})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("should return not found error", func(t *testing.T) {
		err := repository.DeleteOne(ctx, "china")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestTeamRepositoryFind(t *testing.T) {
	ctx := context.Background()

	database := newDatabaseForTest()
	repository, err := mongo.NewTeamRepository(database, testLog)
	assert.NoError(t, err)

	teamCollection := database.Collection(mongo.TeamCollectionName)
	teamCollection.InsertMany(ctx, []interface{}{
		bson.M{"_id": "brazil", "name": "Brazil", "group": "A"},
		bson.M{"_id": "portugal", "name": "Portugal", "group": "B"},
	})

	t.Run("should return teams successfully", func(t *testing.T) {
		teams, err := repository.Find(ctx)
		assert.NoError(t, err)
		assert.Len(t, teams, 2)

		brazilTeam, err := domain.NewTeam("brazil", "Brazil", "A")
		assert.EqualValues(t, brazilTeam, teams[0])

		portugalTeam, err := domain.NewTeam("portugal", "Portugal", "B")
		assert.EqualValues(t, portugalTeam, teams[1])
	})
}
