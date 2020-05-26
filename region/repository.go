package region

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
	Find(context.Context, *provider.Query) (*domain.Region, error)
	FindAll(context.Context, *provider.Query) ([]*domain.Region, error)
}

type RepositoryInstance struct {
	driver neo4j.Driver
}

func NewRepository(driver neo4j.Driver) *RepositoryInstance {
	return &RepositoryInstance{driver: driver}
}

func (r *RepositoryInstance) Find(ctx context.Context, query *provider.Query) (*domain.Region, error) {
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
	p = (city:City)-[*0..1]->()
		WHERE city.code = $`city.code`
		WITH COLLECT(p) AS v, city AS node
		CALL apoc.convert.toTree(v) YIELD value
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

	var region domain.Region
	// Get first index, neo4j.Single handle empty array and array > 1
	if err := provider.RecordUnmarshal(record.GetByIndex(0), &region); err != nil {
		log.WithError(err).Errorln("Failed parse result to struct")
		return nil, err
	}

	return &region, nil
}

func (r *RepositoryInstance) FindAll(ctx context.Context, query *provider.Query) ([]*domain.Region, error) {
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
	MATCH p = (city:City), (city)<-[*]-(province:Province)
		WHERE province.code = $`province.code`
		WITH p, city
		ORDER BY city.name ASC SKIP 0 LIMIT 25
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
			RETURN COLLECT(value) AS value`, match, where, node, order)

	records, err := neo4j.Collect(session.Run(filter, value))
	if err != nil {
		log.WithError(err).Errorln("Failed run action to storage")
		return nil, err
	}

	var results []*domain.Region
	if len(records) > 0 {
		if err := provider.RecordUnmarshal(records[0].GetByIndex(0), &results); err != nil {
			log.WithError(err).Errorln("Failed parse result to struct")
			return nil, err
		}
	}

	if len(results) == 1 {
		if len(results[0].ID) < 1 {
			return nil, nil
		}
	}

	return results, nil
}
