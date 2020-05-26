package web

import (
	"fmt"
	"net/http"

	countryHandler "github.com/dynastymasra/cartographer/country/handler"

	"github.com/graphql-go/graphql"

	regionHandler "github.com/dynastymasra/cartographer/region/handler"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"gopkg.in/tylerb/graceful.v1"
)

func Run(server *graceful.Server, port string, router *RouterInstance) {
	log := logrus.WithFields(logrus.Fields{
		"port":         port,
		"service_name": router.name,
	})

	log.Infoln("Start run web application")

	regionSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: regionHandler.RegionQuery(router.regionRepo),
	})
	if err != nil {
		log.WithError(err).Fatalln("Cannot create new region graph schema")
	}

	countrySchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: countryHandler.CountryQuery(router.countryRepo),
	})
	if err != nil {
		log.WithError(err).Fatalln("Cannot create new country graph schema")
	}

	router.InsertSchema(NewSchema(regionSchema, countrySchema))

	muxRouter := router.Router()

	server.Server = &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(logrus.StandardLogger()),
		)(muxRouter),
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Fatalln("Failed to start server")
	}
}
