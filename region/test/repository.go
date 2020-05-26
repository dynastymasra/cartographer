package test

import (
	"context"

	"github.com/dynastymasra/cartographer/domain"
	"github.com/dynastymasra/cartographer/infrastructure/provider"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Find(ctx context.Context, query *provider.Query) (*domain.Region, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*domain.Region), args.Error(1)
}

func (m *MockRepository) FindAll(ctx context.Context, query *provider.Query) ([]*domain.Region, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*domain.Region), args.Error(1)
}
