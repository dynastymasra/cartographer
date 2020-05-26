package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dynastymasra/cartographer/config"
	"github.com/dynastymasra/cartographer/domain"
	"github.com/dynastymasra/cartographer/infrastructure/provider"
	"github.com/dynastymasra/cartographer/region/handler"
	"github.com/dynastymasra/cartographer/region/test"

	"github.com/stretchr/testify/assert"

	"github.com/dynastymasra/cookbook"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"

	graph "github.com/graphql-go/handler"
	"github.com/stretchr/testify/suite"
)

type RegionSuite struct {
	suite.Suite
	repo *test.MockRepository
}

func Test_RegionSuite(t *testing.T) {
	suite.Run(t, new(RegionSuite))
}

func (r *RegionSuite) SetupSuite() {
	config.SetupTestLogger()
}

func (r *RegionSuite) SetupTest() {
	r.repo = &test.MockRepository{}
}

func (r *RegionSuite) Test_FindRegion_BadRequest() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/regions", nil)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.RegionQuery(r.repo),
	})
	if err != nil {
		r.T().Fatal(err)
	}

	handler.FindRegion(schema)(w, req.WithContext(ctx))

	assert.Equal(r.T(), http.StatusBadRequest, w.Code)
}

func (r *RegionSuite) Test_FindRegion_NotFound() {
	body := []byte(`{"query":"{city(id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\") {id name code createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/regions", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.RegionQuery(r.repo),
	})
	if err != nil {
		r.T().Fatal(err)
	}

	query := provider.NewQuery("City")
	query.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")

	r.repo.On("Find", ctx, query).Return((*domain.Region)(nil), errors.New(provider.ErrorRecordNotFound))

	handler.FindRegion(schema)(w, req.WithContext(ctx))

	assert.Equal(r.T(), http.StatusNotFound, w.Code)
}

func (r *RegionSuite) Test_FindRegion_MoreThanOne() {
	body := []byte(`{"query":"{city(id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\") {id name code createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/regions", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.RegionQuery(r.repo),
	})
	if err != nil {
		r.T().Fatal(err)
	}

	query := provider.NewQuery("City")
	query.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")

	r.repo.On("Find", ctx, query).Return((*domain.Region)(nil), errors.New(provider.ErrorRecordMoreThanOne))

	handler.FindRegion(schema)(w, req.WithContext(ctx))

	assert.Equal(r.T(), http.StatusPreconditionFailed, w.Code)
}

func (r *RegionSuite) Test_FindRegion_FilterEmpty() {
	body := []byte(`{"query":"{city {id name code createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/regions", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.RegionQuery(r.repo),
	})
	if err != nil {
		r.T().Fatal(err)
	}

	handler.FindRegion(schema)(w, req.WithContext(ctx))

	assert.Equal(r.T(), http.StatusPreconditionFailed, w.Code)
}

func (r *RegionSuite) Test_FindRegion_Error() {
	body := []byte(`{"query":"{city(id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\") {id name code createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/regions", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.RegionQuery(r.repo),
	})
	if err != nil {
		r.T().Fatal(err)
	}

	query := provider.NewQuery("City")
	query.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")

	r.repo.On("Find", ctx, query).Return((*domain.Region)(nil), assert.AnError)

	handler.FindRegion(schema)(w, req.WithContext(ctx))

	assert.Equal(r.T(), http.StatusInternalServerError, w.Code)
}

func (r *RegionSuite) Test_FindRegion_Success() {
	body := []byte(`{"query":"{city(id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\") {id name code createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/regions", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.RegionQuery(r.repo),
	})
	if err != nil {
		r.T().Fatal(err)
	}

	query := provider.NewQuery("City")
	query.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")

	res := &domain.Region{
		ID:        uuid.NewV4().String(),
		Name:      "Fukuoka",
		Code:      "1",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	r.repo.On("Find", ctx, query).Return(res, nil)

	handler.FindRegion(schema)(w, req.WithContext(ctx))

	assert.Equal(r.T(), http.StatusOK, w.Code)
}

func (r *RegionSuite) Test_FindListRegion_Success() {
	body := []byte(`{"query":"{cities(code: \"1\", country: {id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\"}) {id name code createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/regions", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.RegionQuery(r.repo),
	})
	if err != nil {
		r.T().Fatal(err)
	}

	query2 := provider.NewQuery("Country")
	query2.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")
	query := provider.NewQuery("City")
	query.Filter("code", provider.Equal, "1")
	query.Slice(config.Offset, config.Limit)
	query.Ordering("name", provider.Ascending)
	query.Incoming(query2)

	timestamp := time.Now().UTC()
	res := []*domain.Region{
		{
			ID:        uuid.NewV4().String(),
			Name:      "Fukuoka",
			Code:      "1",
			CreatedAt: timestamp,
			UpdatedAt: timestamp,
		}, {
			ID:        uuid.NewV4().String(),
			Name:      "Kyoto",
			Code:      "2",
			CreatedAt: timestamp,
			UpdatedAt: timestamp,
		}, {
			ID:        uuid.NewV4().String(),
			Name:      "Osaka",
			Code:      "3",
			CreatedAt: timestamp,
			UpdatedAt: timestamp,
		}, {
			ID:        uuid.NewV4().String(),
			Name:      "Tokyo",
			Code:      "4",
			CreatedAt: timestamp,
			UpdatedAt: timestamp,
		},
	}
	r.repo.On("FindAll", ctx, query).Return(res, nil)

	handler.FindRegion(schema)(w, req.WithContext(ctx))

	assert.Equal(r.T(), http.StatusOK, w.Code)
}

func (r *RegionSuite) Test_FindListRegion_Failed() {
	body := []byte(`{"query":"{cities(country: {id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\"}) {id name code createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/regions", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.RegionQuery(r.repo),
	})
	if err != nil {
		r.T().Fatal(err)
	}

	query2 := provider.NewQuery("Country")
	query2.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")
	query := provider.NewQuery("City")
	query.Ordering("name", provider.Ascending)
	query.Incoming(query2)

	r.repo.On("FindAll", ctx, query).Return(([]*domain.Region)(nil), nil)

	handler.FindRegion(schema)(w, req.WithContext(ctx))

	assert.Equal(r.T(), http.StatusBadRequest, w.Code)
}

func (r *RegionSuite) Test_FindListRegion_Error() {
	body := []byte(`{"query":"{cities(country: {id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\"}) {id name code createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/regions", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.RegionQuery(r.repo),
	})
	if err != nil {
		r.T().Fatal(err)
	}

	query2 := provider.NewQuery("Country")
	query2.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")
	query := provider.NewQuery("City")
	query.Slice(config.Offset, config.Limit)
	query.Ordering("name", provider.Ascending)
	query.Incoming(query2)

	r.repo.On("FindAll", ctx, query).Return(([]*domain.Region)(nil), assert.AnError)

	handler.FindRegion(schema)(w, req.WithContext(ctx))

	assert.Equal(r.T(), http.StatusInternalServerError, w.Code)
}
