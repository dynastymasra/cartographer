package web

import (
	"fmt"
	"net/http"

	"github.com/dynastymasra/cartographer/country"
	"github.com/dynastymasra/cartographer/infrastructure/web/handler"
	"github.com/dynastymasra/cartographer/region"

	"github.com/graphql-go/graphql"

	countryHandler "github.com/dynastymasra/cartographer/country/handler"
	regionHandler "github.com/dynastymasra/cartographer/region/handler"

	"github.com/dynastymasra/cookbook"
	"github.com/dynastymasra/cookbook/negroni/middleware"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/urfave/negroni"
)

const DefaultResponseNotFound = "the requested resource doesn't exists"

type RouterInstance struct {
	name        string
	driver      neo4j.Driver
	regionRepo  region.Repository
	countryRepo country.Repository
	schema      *GraphSchema
}

type GraphSchema struct {
	region  graphql.Schema
	country graphql.Schema
}

func NewSchema(region graphql.Schema, country graphql.Schema) *GraphSchema {
	return &GraphSchema{
		region:  region,
		country: country,
	}
}

func NewRouter(name string, driver neo4j.Driver, regionRepo region.Repository, countryRepo country.Repository) *RouterInstance {
	return &RouterInstance{
		name:        name,
		driver:      driver,
		regionRepo:  regionRepo,
		countryRepo: countryRepo,
	}
}

func (r *RouterInstance) InsertSchema(schema *GraphSchema) {
	r.schema = schema
}

func (r *RouterInstance) Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true).UseEncodedPath()

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, cookbook.FailResponse(&cookbook.JSON{
			"endpoint": DefaultResponseNotFound,
		}, nil).Stringify())
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, cookbook.FailResponse(&cookbook.JSON{
			"method": DefaultResponseNotFound,
		}, nil).Stringify())
	})

	commonHandlers := negroni.New(
		middleware.RequestID(),
	)

	// Probes
	router.Handle("/ping", commonHandlers.With(
		negroni.WrapFunc(handler.Ping(r.driver)),
	)).Methods(http.MethodGet, http.MethodHead)

	router.Handle("/ping", commonHandlers.With(
		negroni.WrapFunc(handler.Ping(r.driver)),
	)).Methods(http.MethodGet, http.MethodHead)

	subRouter := router.PathPrefix("/v1/").Subrouter().UseEncodedPath()
	commonHandlers.Use(middleware.LogrusLog(r.name))

	subRouter.Handle("/regions", commonHandlers.With(
		negroni.WrapFunc(regionHandler.FindRegion(r.schema.region)),
	)).Methods(http.MethodPost)

	subRouter.Handle("/countries", commonHandlers.With(
		negroni.WrapFunc(countryHandler.FindCountry(r.schema.country)),
	)).Methods(http.MethodPost)

	return router
}
