package mongo

import (
	"context"
	"fmt"
	"log"

	"gitlab.neoway.com.br/diogo.giassi/fifa-wct-go-example/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TeamCollectionName = "teams"
)

type Team struct {
	ID    string `bson:"_id"`
	Name  string `bson:"name"`
	Group string `bson:"group"`
}

func NewTeam(t *domain.Team) *Team {
	return &Team{ID: t.ID, Name: t.Name, Group: t.Group}
}

func (t *Team) ToDomainTeam() (*domain.Team, error) {
	team, err := domain.NewTeam(t.ID, t.Name, t.Group)
	if err != nil {
		return nil, fmt.Errorf("error on create domain team: %w ", err)
	}

	return team, nil
}

type teamRepository struct {
	collection *mongo.Collection
	log        *log.Logger
}

func NewTeamRepository(db *mongo.Database, log *log.Logger) (*teamRepository, error) {
	collection := db.Collection(TeamCollectionName)

	return &teamRepository{
		collection: collection,
		log:        log,
	}, nil
}

func (tr *teamRepository) InsertOne(ctx context.Context, t *domain.Team) (*domain.Team, error) {
	team := NewTeam(t)
	if _, err := tr.collection.InsertOne(ctx, team); err != nil {
		if IsDuplicateKeyError(err) {
			return nil, domain.ErrDuplicateKey
		}

		return nil, fmt.Errorf("error to insert team: %w", err)
	}

	return team.ToDomainTeam()
}

func (tr *teamRepository) FindOne(ctx context.Context, ID string) (*domain.Team, error) {
	var result *Team

	key := bson.M{"_id": ID}

	err := tr.collection.FindOne(ctx, key).Decode(&result)
	if err != nil {
		if IsNotFoundError(err) {
			return nil, fmt.Errorf("error to get team with id '%s': %w", ID, domain.ErrNotFound)
		}

		return nil, fmt.Errorf("error to get team with id '%s': %w", ID, err)
	}

	return result.ToDomainTeam()
}

func (tr *teamRepository) DeleteOne(ctx context.Context, ID string) error {
	key := bson.M{"_id": ID}

	result, err := tr.collection.DeleteOne(ctx, key)
	if err != nil {
		return fmt.Errorf("error to delete team with id '%s': %w", ID, err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (tr *teamRepository) Find(ctx context.Context) ([]*domain.Team, error) {
	cur, err := tr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error to find all teams: %w", err)
	}

	defer func() {
		if err := cur.Close(ctx); err != nil {
			tr.log.Printf("error closing mongo cursor: %v", err)
		}
	}()

	teams := make([]*domain.Team, 0)

	for cur.Next(ctx) {
		var t Team
		if err := cur.Decode(&t); err != nil {
			return nil, fmt.Errorf("error on decode team: %w", err)
		}

		team, err := t.ToDomainTeam()
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}
