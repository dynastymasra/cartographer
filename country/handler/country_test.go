package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	graph "github.com/graphql-go/handler"

	"github.com/dynastymasra/cartographer/config"
	"github.com/dynastymasra/cartographer/country/handler"
	"github.com/dynastymasra/cartographer/country/test"
	"github.com/dynastymasra/cartographer/domain"
	"github.com/dynastymasra/cartographer/infrastructure/provider"
	"github.com/dynastymasra/cookbook"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type CountrySuite struct {
	suite.Suite
	repo *test.MockRepository
}

func Test_RegionSuite(t *testing.T) {
	suite.Run(t, new(CountrySuite))
}

func (c *CountrySuite) SetupSuite() {
	config.SetupTestLogger()
}

func (c *CountrySuite) SetupTest() {
	c.repo = &test.MockRepository{}
}

func (c *CountrySuite) Test_FindCountry_BadRequest() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/countries", nil)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.CountryQuery(c.repo),
	})
	if err != nil {
		c.T().Fatal(err)
	}

	handler.FindCountry(schema)(w, req.WithContext(ctx))

	assert.Equal(c.T(), http.StatusBadRequest, w.Code)
}

func (c *CountrySuite) Test_CountryRegion_NotFound() {
	body := []byte(`{"query":"{country(id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\") {id name createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.CountryQuery(c.repo),
	})
	if err != nil {
		c.T().Fatal(err)
	}

	query := provider.NewQuery("Country")
	query.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")

	c.repo.On("Find", ctx, query).Return((*domain.Country)(nil), errors.New(provider.ErrorRecordNotFound))

	handler.FindCountry(schema)(w, req.WithContext(ctx))

	assert.Equal(c.T(), http.StatusNotFound, w.Code)
}

func (c *CountrySuite) Test_FindCountry_MoreThanOne() {
	body := []byte(`{"query":"{country(id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\") {id name createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.CountryQuery(c.repo),
	})
	if err != nil {
		c.T().Fatal(err)
	}

	query := provider.NewQuery("Country")
	query.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")

	c.repo.On("Find", ctx, query).Return((*domain.Country)(nil), errors.New(provider.ErrorRecordMoreThanOne))

	handler.FindCountry(schema)(w, req.WithContext(ctx))

	assert.Equal(c.T(), http.StatusPreconditionFailed, w.Code)
}

func (c *CountrySuite) Test_FindCountry_FilterEmpty() {
	body := []byte(`{"query":"{country {id name createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.CountryQuery(c.repo),
	})
	if err != nil {
		c.T().Fatal(err)
	}

	handler.FindCountry(schema)(w, req.WithContext(ctx))

	assert.Equal(c.T(), http.StatusPreconditionFailed, w.Code)
}

func (c *CountrySuite) Test_FindCountry_Error() {
	body := []byte(`{"query":"{country(id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\") {id name createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.CountryQuery(c.repo),
	})
	if err != nil {
		c.T().Fatal(err)
	}

	query := provider.NewQuery("Country")
	query.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")

	c.repo.On("Find", ctx, query).Return((*domain.Country)(nil), assert.AnError)

	handler.FindCountry(schema)(w, req.WithContext(ctx))

	assert.Equal(c.T(), http.StatusInternalServerError, w.Code)
}

func (c *CountrySuite) Test_FindCountry_Success() {
	body := []byte(`{"query":"{country(id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\") {id name createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.CountryQuery(c.repo),
	})
	if err != nil {
		c.T().Fatal(err)
	}

	query := provider.NewQuery("Country")
	query.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")

	res := &domain.Country{
		ID:            uuid.NewV4().String(),
		Name:          "Japan",
		ISO3166Alpha2: "JP",
		ISO3166Alpha3: "JPN",
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	c.repo.On("Find", ctx, query).Return(res, nil)

	handler.FindCountry(schema)(w, req.WithContext(ctx))

	assert.Equal(c.T(), http.StatusOK, w.Code)
}

func (c *CountrySuite) Test_FindListCountry_Success() {
	body := []byte(`{"query":"{countries(dialCode: \"1\", currencies: {id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\"}) {id name createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.CountryQuery(c.repo),
	})
	if err != nil {
		c.T().Fatal(err)
	}

	query2 := provider.NewQuery("Currency")
	query2.Filter("id", provider.In, "e81f509f-38ec-42e8-9a1c-8e527977e526")
	query := provider.NewQuery("Country")
	query.Filter("dialCode", provider.Equal, "1")
	query.Slice(config.Offset, config.Limit)
	query.Ordering("name", provider.Ascending)
	query.Outgoing(query2)

	timestamp := time.Now().UTC()
	res := []*domain.Country{
		{
			ID:            uuid.NewV4().String(),
			Name:          "Japan",
			ISO3166Alpha2: "JP",
			ISO3166Alpha3: "JPN",
			CreatedAt:     timestamp,
			UpdatedAt:     timestamp,
		}, {
			ID:            uuid.NewV4().String(),
			Name:          "Germany",
			ISO3166Alpha2: "DE",
			ISO3166Alpha3: "DEU",
			CreatedAt:     timestamp,
			UpdatedAt:     timestamp,
		}, {
			ID:            uuid.NewV4().String(),
			Name:          "New Zealand",
			ISO3166Alpha2: "NZ",
			ISO3166Alpha3: "NZL",
			CreatedAt:     timestamp,
			UpdatedAt:     timestamp,
		}, {
			ID:            uuid.NewV4().String(),
			Name:          "Australia",
			ISO3166Alpha2: "AU",
			ISO3166Alpha3: "AUS",
			CreatedAt:     timestamp,
			UpdatedAt:     timestamp,
		}, {
			ID:            uuid.NewV4().String(),
			Name:          "Turkey",
			ISO3166Alpha2: "TR",
			ISO3166Alpha3: "TUR",
			CreatedAt:     timestamp,
			UpdatedAt:     timestamp,
		},
	}
	c.repo.On("FindAll", ctx, query).Return(res, nil)

	handler.FindCountry(schema)(w, req.WithContext(ctx))

	assert.Equal(c.T(), http.StatusOK, w.Code)
}

func (c *CountrySuite) Test_FindListCountry_Error() {
	body := []byte(`{"query":"{countries(dialCode: \"1\", currencies: {id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\"}) {id name createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.CountryQuery(c.repo),
	})
	if err != nil {
		c.T().Fatal(err)
	}

	query2 := provider.NewQuery("Currency")
	query2.Filter("id", provider.In, "e81f509f-38ec-42e8-9a1c-8e527977e526")
	query := provider.NewQuery("Country")
	query.Filter("dialCode", provider.Equal, "1")
	query.Slice(config.Offset, config.Limit)
	query.Ordering("name", provider.Ascending)
	query.Outgoing(query2)

	c.repo.On("FindAll", ctx, query).Return(([]*domain.Country)(nil), assert.AnError)

	handler.FindCountry(schema)(w, req.WithContext(ctx))

	assert.Equal(c.T(), http.StatusInternalServerError, w.Code)
}

func (c *CountrySuite) Test_FindListCountry_Failed() {
	body := []byte(`{"query":"{countries(dialCode: \"1\", currencies: {id: \"e81f509f-38ec-42e8-9a1c-8e527977e526\"}) {id name createdAt updatedAt}}"}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", graph.ContentTypeJSON)

	ctx := context.WithValue(req.Context(), cookbook.RequestID, uuid.NewV4().String())

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: handler.CountryQuery(c.repo),
	})
	if err != nil {
		c.T().Fatal(err)
	}

	query2 := provider.NewQuery("Currency")
	query2.Filter("id", provider.Equal, "e81f509f-38ec-42e8-9a1c-8e527977e526")
	query := provider.NewQuery("Country")
	query.Filter("dialCode", provider.Equal, "1")
	query.Slice(config.Offset, config.Limit)
	query.Ordering("name", provider.Ascending)
	query.Outgoing(query2)

	timestamp := time.Now().UTC()
	res := []*domain.Country{
		{
			ID:            uuid.NewV4().String(),
			Name:          "Japan",
			ISO3166Alpha2: "JP",
			ISO3166Alpha3: "JPN",
			CreatedAt:     timestamp,
			UpdatedAt:     timestamp,
		}, {
			ID:            uuid.NewV4().String(),
			Name:          "Germany",
			ISO3166Alpha2: "DE",
			ISO3166Alpha3: "DEU",
			CreatedAt:     timestamp,
			UpdatedAt:     timestamp,
		}, {
			ID:            uuid.NewV4().String(),
			Name:          "New Zealand",
			ISO3166Alpha2: "NZ",
			ISO3166Alpha3: "NZL",
			CreatedAt:     timestamp,
			UpdatedAt:     timestamp,
		}, {
			ID:            uuid.NewV4().String(),
			Name:          "Australia",
			ISO3166Alpha2: "AU",
			ISO3166Alpha3: "AUS",
			CreatedAt:     timestamp,
			UpdatedAt:     timestamp,
		}, {
			ID:            uuid.NewV4().String(),
			Name:          "Turkey",
			ISO3166Alpha2: "TR",
			ISO3166Alpha3: "TUR",
			CreatedAt:     timestamp,
			UpdatedAt:     timestamp,
		},
	}
	c.repo.On("FindAll", ctx, query).Return(res, nil)

	handler.FindCountry(schema)(w, req.WithContext(ctx))

	assert.Equal(c.T(), http.StatusBadRequest, w.Code)
}
