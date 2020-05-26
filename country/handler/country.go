package handler

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/dynastymasra/cartographer/config"
	"github.com/dynastymasra/cartographer/country"
	"github.com/dynastymasra/cartographer/domain"
	"github.com/dynastymasra/cartographer/infrastructure/provider"
	"github.com/dynastymasra/cookbook"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	graph "github.com/graphql-go/handler"
	"github.com/sirupsen/logrus"
)

func FindCountry(schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		requestID := r.Context().Value(cookbook.RequestID).(string)
		log := logrus.WithFields(logrus.Fields{
			cookbook.RequestID: requestID,
			"package":          runtime.FuncForPC(reflect.ValueOf(FindCountry).Pointer()).Name(),
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

func CountryQuery(repo country.Repository) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Query",
			Description: "Query country from storage",
			Fields: graphql.Fields{
				"country": &graphql.Field{
					Type: domain.CountryType,
					Args: domain.CountryArgs,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						log := logrus.WithFields(logrus.Fields{
							cookbook.RequestID: p.Context.Value(cookbook.RequestID),
							"package":          runtime.FuncForPC(reflect.ValueOf(CountryQuery).Pointer()).Name(),
							"arguments":        cookbook.Stringify(p.Args),
						})

						query := provider.NewQuery(domain.CountryNode)

						for key, field := range p.Args {
							query.Filter(key, provider.Equal, field)
						}

						if len(query.Filters) < 1 {
							log.WithField("query", cookbook.Stringify(query)).Warnln("Query is empty")
							return nil, config.NewError(http.StatusPreconditionFailed, strings.ToLower(query.Node), "need min one argument")
						}

						res, err := repo.Find(p.Context, query)
						if err != nil {
							if err.Error() == provider.ErrorRecordNotFound {
								return nil, config.NewError(http.StatusNotFound, strings.ToLower(query.Node), err.Error())
							}

							if err.Error() == provider.ErrorRecordMoreThanOne {
								return nil, config.NewError(http.StatusPreconditionFailed, strings.ToLower(query.Node), err.Error())
							}

							log.WithField("query", cookbook.Stringify(query)).WithError(err).Errorln("Failed find country from storage")
							return nil, config.NewError(http.StatusInternalServerError, "", err.Error())
						}

						return res, nil
					},
				},
				"countries": &graphql.Field{
					Type: graphql.NewList(domain.CountryType),
					Args: domain.ListCountryArgs,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						log := logrus.WithFields(logrus.Fields{
							cookbook.RequestID: p.Context.Value(cookbook.RequestID),
							"package":          runtime.FuncForPC(reflect.ValueOf(CountryQuery).Pointer()).Name(),
							"arguments":        cookbook.Stringify(p.Args),
						})

						limit := p.Args["limit"].(int)
						offset := p.Args["offset"].(int)

						// Delete value from map avoid duplicate value
						delete(p.Args, "limit")
						delete(p.Args, "offset")

						query := provider.NewQuery(domain.CountryNode)
						query.Slice(offset, limit)
						query.Ordering("name", provider.Ascending)

						for key, field := range p.Args {
							switch val := field.(type) {
							case []interface{}:
								if node, ok := domain.Outgoing[key]; ok {
									q := provider.NewQuery(node)
									for _, v := range val {
										if va, ok := v.(map[string]interface{}); ok {
											for k, v := range va {
												q.Filter(k, provider.In, v)
											}
										}
									}
									query.Outgoing(q)
								}
							default:
								query.Filter(key, provider.Equal, field)
							}
						}

						results, err := repo.FindAll(p.Context, query)
						if err != nil {
							log.WithField("query", cookbook.Stringify(query)).WithError(err).Errorln("Failed find country from storage")
							return nil, config.NewError(http.StatusInternalServerError, "", err.Error())
						}

						return results, nil
					},
				},
			},
		})
}
