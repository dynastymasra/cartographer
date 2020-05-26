package handler_test

import (
	"net/http"
	"testing"

	"github.com/dynastymasra/cartographer/config"
	"github.com/dynastymasra/cartographer/infrastructure/provider/test"
	"github.com/dynastymasra/cartographer/infrastructure/web/handler"

	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type PingSuite struct {
	suite.Suite
	provider *test.MockNeo4J
}

func Test_PingSuite(t *testing.T) {
	suite.Run(t, new(PingSuite))
}

func (p *PingSuite) SetupSuite() {
	config.SetupTestLogger()
}

func (p *PingSuite) SetupTest() {
	p.provider = &test.MockNeo4J{}
}

func (p *PingSuite) Test_PingHandler() {
	p.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(p.provider, nil)
	p.provider.On("Close").Return(nil)

	assert.HTTPSuccess(p.T(), handler.Ping(p.provider), http.MethodGet, "/ping", nil)
	assert.HTTPSuccess(p.T(), handler.Ping(p.provider), http.MethodHead, "/ping", nil)
	assert.HTTPBodyContains(p.T(), handler.Ping(p.provider), http.MethodGet, "/ping", nil, "{\"status\":\"success\"}")
}

func (p *PingSuite) Test_PingHandler_Error() {
	p.provider.On("Session", neo4j.AccessModeRead, []string(nil)).Return(p.provider, assert.AnError)

	assert.HTTPError(p.T(), handler.Ping(p.provider), http.MethodGet, "/ping", nil)
	assert.HTTPError(p.T(), handler.Ping(p.provider), http.MethodHead, "/ping", nil)
}
