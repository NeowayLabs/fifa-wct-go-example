//go:build unit

package application

import (
	"context"

	"github.com/NeowayLabs/fifa-wct-go-example/internal/domain"
)

type TeamServiceMock struct {
	CreateFn           func(ctx context.Context, input *TeamInput) (*TeamOutput, error)
	CreateInvokedCount int
	GetFn              func(ctx context.Context, ID string) (*TeamOutput, error)
	GetInvokedCount    int
	RemoveFn           func(ctx context.Context, ID string) error
	RemoveInvokedCount int
	GetAllFn           func(ctx context.Context) ([]*TeamOutput, error)
	GetAllInvokedCount int
}

func (tsm *TeamServiceMock) Create(ctx context.Context, input *TeamInput) (*TeamOutput, error) {
	tsm.CreateInvokedCount++
	return tsm.CreateFn(ctx, input)
}

func (tsm *TeamServiceMock) Get(ctx context.Context, ID string) (*TeamOutput, error) {
	tsm.GetInvokedCount++
	return tsm.GetFn(ctx, ID)
}

func (tsm *TeamServiceMock) Remove(ctx context.Context, ID string) error {
	tsm.RemoveInvokedCount++
	return tsm.RemoveFn(ctx, ID)
}

func (tsm *TeamServiceMock) GetAll(ctx context.Context) ([]*TeamOutput, error) {
	tsm.GetAllInvokedCount++
	return tsm.GetAllFn(ctx)
}

type TeamRepositoryMock struct {
	InsertOneFn           func(ctx context.Context, team *domain.Team) (*domain.Team, error)
	InsertOneInvokedCount int
	FindOneFn             func(ctx context.Context, ID string) (*domain.Team, error)
	FindOneInvokedCount   int
	DeleteOneFn           func(ctx context.Context, ID string) error
	DeleteOneInvokedCount int
	FindFn                func(ctx context.Context) ([]*domain.Team, error)
	FindInvokedCount      int
}

func (trm *TeamRepositoryMock) InsertOne(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	trm.InsertOneInvokedCount++
	return trm.InsertOneFn(ctx, team)
}

func (trm *TeamRepositoryMock) FindOne(ctx context.Context, ID string) (*domain.Team, error) {
	trm.FindOneInvokedCount++
	return trm.FindOneFn(ctx, ID)
}

func (trm *TeamRepositoryMock) DeleteOne(ctx context.Context, ID string) error {
	trm.DeleteOneInvokedCount++
	return trm.DeleteOneFn(ctx, ID)
}

func (trm *TeamRepositoryMock) Find(ctx context.Context) ([]*domain.Team, error) {
	trm.FindInvokedCount++
	return trm.FindFn(ctx)
}
