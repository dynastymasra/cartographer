package country

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/dynastymasra/cartographer/domain"
	"github.com/dynastymasra/cartographer/infrastructure/provider"

	"github.com/dynastymasra/cookbook"
	"github.com/sirupsen/logrus"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Repository interface {
	Find(context.Context, *provider.Query) (*domain.Country, error)
	FindAll(context.Context, *provider.Query) ([]*domain.Country, error)
}

type RepositoryInstance struct {
	driver neo4j.Driver
}

func NewRepository(driver neo4j.Driver) *RepositoryInstance {
	return &RepositoryInstance{driver: driver}
}

func (r *RepositoryInstance) Find(ctx context.Context, query *provider.Query) (*domain.Country, error) {
	log := logrus.WithFields(logrus.Fields{
		cookbook.RequestID: ctx.Value(cookbook.RequestID),
		"package":          runtime.FuncForPC(reflect.ValueOf(r.Find).Pointer()).Name(),
	})

	session, err := r.driver.Session(neo4j.AccessModeRead)
	if err != nil {
		log.WithError(err).Errorln("Failed create new session")
		return nil, err
	}
	defer session.Close()

	node := strings.ToLower(query.Node)
	match, where, _, value := provider.TranslateQuery(query)

	/**
	MATCH p = (country:Country)-[*0..1]->()
		WHERE country.ISO3166Alpha2 = "ID"
		WITH COLLECT(p) AS val, country AS node
		CALL apoc.convert.toTree(val) YIELD value
	RETURN value, node
	*/
	filter := fmt.Sprintf(`MATCH p = %s-[*0..1]->()
			%s
			WITH COLLECT(p) AS v, %s AS node
			CALL apoc.convert.toTree(v) YIELD value
			RETURN value, node`,
		match, where, node)

	record, err := neo4j.Single(session.Run(filter, value))
	if err != nil {
		log.WithError(err).Warnln("Failed run action to storage")
		return nil, err
	}

	var country domain.Country
	// Get first index, neo4j.Single handle empty array and array > 1
	if err := provider.RecordUnmarshal(record.GetByIndex(0), &country); err != nil {
		log.WithError(err).Errorln("Failed parse result to struct")
		return nil, err
	}

	return &country, nil
}

func (r *RepositoryInstance) FindAll(ctx context.Context, query *provider.Query) ([]*domain.Country, error) {
	log := logrus.WithFields(logrus.Fields{
		cookbook.RequestID: ctx.Value(cookbook.RequestID),
		"package":          runtime.FuncForPC(reflect.ValueOf(r.FindAll).Pointer()).Name(),
	})

	session, err := r.driver.Session(neo4j.AccessModeRead)
	if err != nil {
		log.WithError(err).Errorln("Failed create new session")
		return nil, err
	}
	defer session.Close()

	node := strings.ToLower(query.Node)
	match, where, order, value := provider.TranslateQuery(query)

	/**
	MATCH p = (country:Country), (country)-[*]->(province:Province)
		WHERE province.code="34"
		WITH p, country
		ORDER BY country.name ASC SKIP 0 LIMIT 5
		WITH COLLECT(p) AS val
		CALL apoc.convert.toTree(val) YIELD value
	RETURN COLLECT(value) AS value
	*/
	filter := fmt.Sprintf(`MATCH p = %s
			%s
			WITH p, %s
			%s
			WITH COLLECT(p) AS val
			CALL apoc.convert.toTree(val) YIELD value
		RETURN COLLECT(value)`, match, where, node, order)

	records, err := neo4j.Collect(session.Run(filter, value))
	if err != nil {
		log.WithError(err).Errorln("Failed run action to storage")
		return nil, err
	}

	var countries []*domain.Country
	if len(records) > 0 {
		if err := provider.RecordUnmarshal(records[0].GetByIndex(0), &countries); err != nil {
			log.WithError(err).Errorln("Failed parse result to struct")
			return nil, err
		}
	}

	if len(countries) == 1 {
		if len(countries[0].ID) < 1 {
			return nil, nil
		}
	}

	return countries, nil
}
