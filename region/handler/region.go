package handler

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/dynastymasra/cartographer/config"
	"github.com/dynastymasra/cartographer/domain"
	"github.com/dynastymasra/cartographer/infrastructure/provider"
	"github.com/dynastymasra/cartographer/region"

	"github.com/graphql-go/graphql/gqlerrors"

	"github.com/sirupsen/logrus"

	"github.com/dynastymasra/cookbook"
	"github.com/graphql-go/graphql"
	graph "github.com/graphql-go/handler"
)

func FindRegion(schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		requestID := r.Context().Value(cookbook.RequestID).(string)
		log := logrus.WithFields(logrus.Fields{
			cookbook.RequestID: requestID,
			"package":          runtime.FuncForPC(reflect.ValueOf(FindRegion).Pointer()).Name(),
		})

		req := graph.NewRequestOptions(r)

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  req.Query,
			VariableValues: req.Variables,
			OperationName:  req.OperationName,
			Context:        r.Context(),
		})

		if len(result.Errors) > 0 {
			var errs []cookbook.JSON
			for _, err := range result.Errors {
				log.WithError(err).Warnln("Failed process request")

				switch err := err.OriginalError().(type) {
				case *gqlerrors.Error:
					if err, ok := err.OriginalError.(*config.ServiceError); ok {
						if err.Code() >= http.StatusInternalServerError {
							log.WithError(err).Errorln("Failed process data from storage")
						}

						config.ParseToJSON(err, w, requestID)
						return
					}
					errs = append(errs, cookbook.JSON{
						"message": err.Error(),
					})
				default:
					errs = append(errs, cookbook.JSON{
						"message": err.Error(),
					})
				}
			}

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, cookbook.FailResponse(&cookbook.JSON{"errors": errs}, requestID).Stringify())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, cookbook.SuccessDataResponse(result.Data, nil).Stringify())
	}
}

func RegionQuery(repo region.Repository) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Query",
			Description: "Query region data from storage",
			Fields: graphql.Fields{
				"province": &graphql.Field{
					Type:    domain.ProvinceType,
					Args:    domain.RegionArgs,
					Resolve: RegionResolver(domain.ProvinceNode, repo),
				},
				"provinces": &graphql.Field{
					Type:    graphql.NewList(domain.ProvinceType),
					Args:    domain.ListRegionArgs,
					Resolve: ListRegionResolver(domain.ProvinceNode, repo),
				},
				"city": &graphql.Field{
					Type:    domain.CityType,
					Args:    domain.RegionArgs,
					Resolve: RegionResolver(domain.CityNode, repo),
				},
				"cities": &graphql.Field{
					Type:    graphql.NewList(domain.CityType),
					Args:    domain.ListRegionArgs,
					Resolve: ListRegionResolver(domain.CityNode, repo),
				},
				"regency": &graphql.Field{
					Type:    domain.RegencyType,
					Args:    domain.RegionArgs,
					Resolve: RegionResolver(domain.RegencyNode, repo),
				},
				"regencies": &graphql.Field{
					Type:    graphql.NewList(domain.RegencyType),
					Args:    domain.ListRegionArgs,
					Resolve: ListRegionResolver(domain.RegencyNode, repo),
				},
				"district": &graphql.Field{
					Type:    domain.DistrictType,
					Args:    domain.RegionArgs,
					Resolve: RegionResolver(domain.DistrictNode, repo),
				},
				"districts": &graphql.Field{
					Type:    graphql.NewList(domain.DistrictType),
					Args:    domain.ListRegionArgs,
					Resolve: ListRegionResolver(domain.DistrictNode, repo),
				},
				"village": &graphql.Field{
					Type:    domain.VillageType,
					Args:    domain.RegionArgs,
					Resolve: RegionResolver(domain.VillageNode, repo),
				},
				"villages": &graphql.Field{
					Type:    graphql.NewList(domain.VillageType),
					Args:    domain.ListRegionArgs,
					Resolve: ListRegionResolver(domain.VillageNode, repo),
				},
			},
		})
}

func RegionResolver(node string, repo region.Repository) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		log := logrus.WithFields(logrus.Fields{
			cookbook.RequestID: p.Context.Value(cookbook.RequestID),
			"package":          runtime.FuncForPC(reflect.ValueOf(RegionResolver).Pointer()).Name(),
			"arguments":        cookbook.Stringify(p.Args),
		})

		query := provider.NewQuery(node)

		for key, field := range p.Args {
			query.Filter(key, provider.Equal, field)
		}

		if len(query.Filters) < 1 {
			log.WithField("query", cookbook.Stringify(query)).Warnln("Query is empty")
			return nil, config.NewError(http.StatusPreconditionFailed, strings.ToLower(node), "need min one argument")
		}

		res, err := repo.Find(p.Context, query)
		if err != nil {
			if err.Error() == provider.ErrorRecordNotFound {
				return nil, config.NewError(http.StatusNotFound, strings.ToLower(node), err.Error())
			}

			if err.Error() == provider.ErrorRecordMoreThanOne {
				return nil, config.NewError(http.StatusPreconditionFailed, strings.ToLower(node), err.Error())
			}

			log.WithField("query", cookbook.Stringify(query)).WithError(err).Errorln("Failed find region from storage")
			return nil, config.NewError(http.StatusInternalServerError, "", err.Error())
		}

		return res, nil
	}
}

func ListRegionResolver(node string, repo region.Repository) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		log := logrus.WithFields(logrus.Fields{
			cookbook.RequestID: p.Context.Value(cookbook.RequestID),
			"package":          runtime.FuncForPC(reflect.ValueOf(ListRegionResolver).Pointer()).Name(),
			"arguments":        cookbook.Stringify(p.Args),
		})

		limit := p.Args["limit"].(int)
		offset := p.Args["offset"].(int)

		// Delete value from map avoid duplicate value
		delete(p.Args, "limit")
		delete(p.Args, "offset")
		delete(p.Args, strings.ToLower(node))

		query := provider.NewQuery(node)
		query.Slice(offset, limit)
		query.Ordering("name", provider.Ascending)

		for key, field := range p.Args {
			switch val := field.(type) {
			case map[string]interface{}:
				if node, ok := domain.Incoming[key]; ok {
					q := provider.NewQuery(node)
					for k, v := range val {
						q.Filter(k, provider.Equal, v)
					}
					query.Incoming(q)
				}
			default:
				query.Filter(key, provider.Equal, field)
			}
		}

		results, err := repo.FindAll(p.Context, query)
		if err != nil {
			log.WithField("query", cookbook.Stringify(query)).WithError(err).Errorln("Failed find region from storage")
			return nil, config.NewError(http.StatusInternalServerError, "", err.Error())
		}

		return results, nil
	}
}
