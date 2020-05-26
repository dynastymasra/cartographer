package test

import (
	"net/url"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/mock"
)

type MockNeo4J struct {
	mock.Mock
	Count int
}

type MockNeo4JRecord struct {
	mock.Mock
}

func (n *MockNeo4J) Target() url.URL {
	args := n.Called()
	return args.Get(0).(url.URL)
}

func (n *MockNeo4J) Session(accessMode neo4j.AccessMode, bookmarks ...string) (neo4j.Session, error) {
	args := n.Called(accessMode, bookmarks)
	return args.Get(0).(neo4j.Session), args.Error(1)
}

func (n *MockNeo4J) Close() error {
	args := n.Called()
	return args.Error(0)
}

func (n *MockNeo4J) LastBookmark() string {
	args := n.Called()
	return args.String(0)
}

func (n *MockNeo4J) BeginTransaction(configurers ...func(*neo4j.TransactionConfig)) (neo4j.Transaction, error) {
	args := n.Called(configurers)
	return args.Get(0).(neo4j.Transaction), args.Error(1)
}

func (n *MockNeo4J) ReadTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	args := n.Called(work, configurers)
	return args.Get(0).(interface{}), args.Error(1)
}

func (n *MockNeo4J) WriteTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	args := n.Called(work, configurers)
	return args.Get(0).(interface{}), args.Error(1)
}

func (n *MockNeo4J) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	args := n.Called(cypher, params)
	return args.Get(0).(neo4j.Result), args.Error(1)
}

func (n *MockNeo4J) Keys() ([]string, error) {
	args := n.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (n *MockNeo4J) Next() bool {
	_ = n.Called()
	if n.Count > 0 {
		return false
	}
	n.Count++
	return true
}

func (n *MockNeo4J) Err() error {
	args := n.Called()
	return args.Error(0)
}

func (n *MockNeo4J) Record() neo4j.Record {
	args := n.Called()
	return args.Get(0).(neo4j.Record)
}

func (n *MockNeo4J) Summary() (neo4j.ResultSummary, error) {
	args := n.Called()
	return args.Get(0).(neo4j.ResultSummary), args.Error(1)
}

func (n *MockNeo4J) Consume() (neo4j.ResultSummary, error) {
	args := n.Called()
	return args.Get(0).(neo4j.ResultSummary), args.Error(1)
}

func (n *MockNeo4JRecord) GetByIndex(index int) interface{} {
	args := n.Called(index)
	return args.Get(0).(interface{})
}

func (n *MockNeo4JRecord) Keys() []string {
	args := n.Called()
	return args.Get(0).([]string)
}

func (n *MockNeo4JRecord) Values() []interface{} {
	args := n.Called()
	return args.Get(0).([]interface{})
}

func (n *MockNeo4JRecord) Get(key string) (interface{}, bool) {
	args := n.Called(key)
	return args.Get(0).(interface{}), args.Bool(1)
}
