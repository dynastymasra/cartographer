package handler

import (
	"fmt"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/dynastymasra/cookbook"
	"github.com/sirupsen/logrus"
)

func Ping(driver neo4j.Driver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log := logrus.WithField(cookbook.RequestID, r.Context().Value(cookbook.RequestID))

		session, err := driver.Session(neo4j.AccessModeRead)
		if err != nil {
			log.WithError(err).Errorln("Failed connect to Neo4J")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, cookbook.ErrorResponse(err.Error(), r.Context().Value(cookbook.RequestID)).Stringify())
			return
		}
		defer session.Close()

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, cookbook.SuccessResponse().Stringify())
	}
}
