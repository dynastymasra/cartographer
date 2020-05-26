package country_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dynastymasra/cartographer/config"
	"github.com/dynastymasra/cartographer/country"
	"github.com/dynastymasra/cartographer/domain"
	"github.com/dynastymasra/cartographer/infrastructure/provider"
	"github.com/dynastymasra/cartographer/infrastructure/provider/test"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct {
	suite.Suite
	provider *test.MockNeo4J
	record   *test.MockNeo4JRecord
}

func Test_RepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (r *RepositorySuite) SetupSuite() {
	config.SetupTestLogger()
}

func (r *RepositorySuite) SetupTest() {
	r.provider = &test.MockNeo4J{}
	r.record = &test.MockNeo4JRecord{}
}

func (r *RepositorySuite) Test_Find_ErrorSession() {
	r.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(r.provider, assert.AnError)

	repo := country.NewRepository(r.provider)

	res, err := repo.Find(context.Background(), &provider.Query{})

	assert.Nil(r.T(), res)
	assert.Error(r.T(), err)
}

func (r *RepositorySuite) Test_Find_Error() {
	r.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(r.provider, nil)
	r.provider.On("Close").Return(nil)

	query := provider.NewQuery("Test")
	match, where, _, value := provider.TranslateQuery(query)

	filter := fmt.Sprintf(`MATCH p = %s-[*0..1]->()
			%s
			WITH COLLECT(p) AS v, %s AS node
			CALL apoc.convert.toTree(v) YIELD value
			RETURN value, node`,
		match, where, "test")

	r.provider.On("Run", filter, value).Return(r.provider, assert.AnError)

	repo := country.NewRepository(r.provider)
	res, err := repo.Find(context.Background(), query)

	assert.Nil(r.T(), res)
	assert.Error(r.T(), err)
}

func (r *RepositorySuite) Test_Find_ErrorUnmarshal() {
	r.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(r.provider, nil)
	r.provider.On("Close").Return(nil)

	query := provider.NewQuery("Test")
	match, where, _, value := provider.TranslateQuery(query)

	filter := fmt.Sprintf(`MATCH p = %s-[*0..1]->()
			%s
			WITH COLLECT(p) AS v, %s AS node
			CALL apoc.convert.toTree(v) YIELD value
			RETURN value, node`,
		match, where, "test")

	r.provider.On("Run", filter, value).Return(r.provider, nil)
	r.provider.On("Next").Return()
	r.provider.On("Record").Return(r.record, nil)
	r.provider.On("Err").Return(nil)
	r.record.On("GetByIndex", 0).Return("<-")

	repo := country.NewRepository(r.provider)
	res, err := repo.Find(context.Background(), query)

	assert.Nil(r.T(), res)
	assert.Error(r.T(), err)
}

func (r *RepositorySuite) Test_Find_Success() {
	r.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(r.provider, nil)
	r.provider.On("Close").Return(nil)

	query := provider.NewQuery("Test")
	match, where, _, value := provider.TranslateQuery(query)

	filter := fmt.Sprintf(`MATCH p = %s-[*0..1]->()
			%s
			WITH COLLECT(p) AS v, %s AS node
			CALL apoc.convert.toTree(v) YIELD value
			RETURN value, node`,
		match, where, "test")

	r.provider.On("Run", filter, value).Return(r.provider, nil)
	r.provider.On("Next").Return()
	r.provider.On("Record").Return(r.record, nil)
	r.provider.On("Err").Return(nil)

	timestamp := time.Now().UTC()
	r.record.On("GetByIndex", 0).Return(domain.Country{
		ID:            uuid.NewV4().String(),
		Name:          "Japan",
		ISO3166Alpha2: "JP",
		ISO3166Alpha3: "JPN",
		CreatedAt:     timestamp,
		UpdatedAt:     timestamp,
	})

	repo := country.NewRepository(r.provider)
	res, err := repo.Find(context.Background(), query)

	assert.NotNil(r.T(), res)
	assert.NoError(r.T(), err)
}

func (r *RepositorySuite) Test_FindAll_ErrorSession() {
	r.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(r.provider, assert.AnError)

	repo := country.NewRepository(r.provider)

	res, err := repo.FindAll(context.Background(), &provider.Query{})

	assert.Nil(r.T(), res)
	assert.Error(r.T(), err)
}

func (r *RepositorySuite) Test_FindAll_Error() {
	r.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(r.provider, nil)
	r.provider.On("Close").Return(nil)

	query := provider.NewQuery("Test")
	match, where, order, value := provider.TranslateQuery(query)

	filter := fmt.Sprintf(`MATCH p = %s
			%s
			WITH p, %s
			%s
			WITH COLLECT(p) AS val
			CALL apoc.convert.toTree(val) YIELD value
		RETURN COLLECT(value)`, match, where, "test", order)

	r.provider.On("Run", filter, value).Return(r.provider, assert.AnError)

	repo := country.NewRepository(r.provider)
	res, err := repo.FindAll(context.Background(), query)

	assert.Nil(r.T(), res)
	assert.Error(r.T(), err)
}

func (r *RepositorySuite) Test_FindAll_ErrorUnmarshal() {
	r.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(r.provider, nil)
	r.provider.On("Close").Return(nil)

	query := provider.NewQuery("Test")
	match, where, order, value := provider.TranslateQuery(query)

	filter := fmt.Sprintf(`MATCH p = %s
			%s
			WITH p, %s
			%s
			WITH COLLECT(p) AS val
			CALL apoc.convert.toTree(val) YIELD value
		RETURN COLLECT(value)`, match, where, "test", order)

	r.provider.On("Run", filter, value).Return(r.provider, nil)
	r.provider.On("Next").Return()
	r.provider.On("Record").Return(r.record, nil)
	r.provider.On("Err").Return(nil)
	r.record.On("GetByIndex", 0).Return("<-")

	repo := country.NewRepository(r.provider)
	res, err := repo.FindAll(context.Background(), query)

	assert.Nil(r.T(), res)
	assert.Error(r.T(), err)
}

func (r *RepositorySuite) Test_FindAll_Empty() {
	r.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(r.provider, nil)
	r.provider.On("Close").Return(nil)

	query := provider.NewQuery("Test")
	match, where, order, value := provider.TranslateQuery(query)

	filter := fmt.Sprintf(`MATCH p = %s
			%s
			WITH p, %s
			%s
			WITH COLLECT(p) AS val
			CALL apoc.convert.toTree(val) YIELD value
		RETURN COLLECT(value)`, match, where, "test", order)

	r.provider.On("Run", filter, value).Return(r.provider, nil)
	r.provider.On("Next").Return()
	r.provider.On("Record").Return(r.record, nil)
	r.provider.On("Err").Return(nil)

	response := []domain.Region{{}}
	r.record.On("GetByIndex", 0).Return(response)

	repo := country.NewRepository(r.provider)
	res, err := repo.FindAll(context.Background(), query)

	assert.Nil(r.T(), res)
	assert.Empty(r.T(), res)
	assert.NoError(r.T(), err)
}

func (r *RepositorySuite) Test_FindAll_Success() {
	r.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(r.provider, nil)
	r.provider.On("Close").Return(nil)

	query := provider.NewQuery("Test")
	match, where, order, value := provider.TranslateQuery(query)

	filter := fmt.Sprintf(`MATCH p = %s
			%s
			WITH p, %s
			%s
			WITH COLLECT(p) AS val
			CALL apoc.convert.toTree(val) YIELD value
		RETURN COLLECT(value)`, match, where, "test", order)

	r.provider.On("Run", filter, value).Return(r.provider, nil)
	r.provider.On("Next").Return()
	r.provider.On("Record").Return(r.record, nil)
	r.provider.On("Err").Return(nil)

	timestamp := time.Now().UTC()
	response := []domain.Country{
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
	r.record.On("GetByIndex", 0).Return(response)

	repo := country.NewRepository(r.provider)
	res, err := repo.FindAll(context.Background(), query)

	assert.NotNil(r.T(), res)
	assert.Len(r.T(), res, 5)
	assert.NotEmpty(r.T(), res)
	assert.NoError(r.T(), err)
}
