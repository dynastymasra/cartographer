package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/matryer/resync"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	ErrorRecordNotFound    = "result contains no records"
	ErrorRecordMoreThanOne = "result contains more than one record"
)

var (
	driver  neo4j.Driver
	err     error
	runOnce resync.Once
)

type Neo4J struct {
	Address     string
	Username    string
	Password    string
	MaxConnPool int
	Encrypted   bool
	LogEnabled  bool
	LogLevel    int
}

func (n Neo4J) Driver() (neo4j.Driver, error) {
	url := fmt.Sprintf("%s", n.Address)
	auth := neo4j.BasicAuth(n.Username, n.Password, "")

	runOnce.Do(func() {
		driver, err = neo4j.NewDriver(url, auth, func(config *neo4j.Config) {
			config.Encrypted = n.Encrypted
			config.MaxConnectionPoolSize = n.MaxConnPool
			if n.LogEnabled {
				config.Log = neo4j.ConsoleLogger(neo4j.LogLevel(n.LogLevel))
			}
		})
	})

	return driver, err
}

const (
	Equal = "Equal"
	In    = "In"

	Descending = "Descending"
	Ascending  = "Ascending"
)

var (
	validOrdering = map[string]bool{
		Descending: true,
		Ascending:  true,
	}
)

type (
	Query struct {
		Node      string
		Limit     int
		Offset    int
		Incomings []*Query
		Outgoings []*Query
		Filters   []*Filter
		Orderings []*Ordering
	}

	Filter struct {
		Condition string
		Field     string
		Value     interface{}
	}

	Ordering struct {
		Field     string
		Direction string
	}
)

func NewQuery(node string) *Query {
	return &Query{
		Node: node,
	}
}

// Order adds a sort order to the query
func (q *Query) Ordering(property, direction string) *Query {
	order := NewOrdering(property, direction)
	q.Orderings = append(q.Orderings, order)
	return q
}

func (q *Query) Incoming(query *Query) *Query {
	q.Incomings = append(q.Incomings, query)
	return q
}

func (q *Query) Outgoing(query *Query) *Query {
	q.Outgoings = append(q.Outgoings, query)
	return q
}

// Filter adds a filter to the query
func (q *Query) Filter(property, condition string, value interface{}) *Query {
	filter := NewFilter(property, condition, value)
	q.Filters = append(q.Filters, filter)
	return q
}

func (q *Query) Slice(offset, limit int) *Query {
	q.Offset = offset
	q.Limit = limit

	return q
}

// NewFilter creates a new property filter
func NewFilter(field, condition string, value interface{}) *Filter {
	return &Filter{
		Field:     field,
		Condition: condition,
		Value:     value,
	}
}

func NewOrdering(field, direction string) *Ordering {
	d := direction

	if !validOrdering[direction] {
		d = Descending
	}

	return &Ordering{
		Field:     field,
		Direction: d,
	}
}

func TranslateQuery(query *Query) (string, string, string, map[string]interface{}) {
	node := strings.ToLower(query.Node)

	var nodes, q []string
	f := make(map[string]interface{})

	nodes = append(nodes, fmt.Sprintf("(%s:%s)", node, query.Node))
	q, f = TranslateFilter(query, q, f)

	for _, val := range query.Incomings {
		nodes = append(nodes, fmt.Sprintf("(%s)<-[*]-(%s:%s)", node, strings.ToLower(val.Node), val.Node))
		q, f = TranslateFilter(val, q, f)
	}

	for _, val := range query.Outgoings {
		nodes = append(nodes, fmt.Sprintf("(%s)-[*]->(%s:%s)", node, strings.ToLower(val.Node), val.Node))
		q, f = TranslateFilter(val, q, f)
	}

	var o []string
	for _, order := range query.Orderings {
		switch order.Direction {
		case Ascending:
			o = append(o, fmt.Sprintf("ORDER BY %s.%s %s", node, order.Field, "ASC"))
		case Descending:
			o = append(o, fmt.Sprintf("ORDER BY %s.%s %s", node, order.Field, "DESC"))
		default:
			o = append(o, fmt.Sprintf("ORDER BY %s.%s %s", node, order.Field, "ASC"))
		}
	}

	if query.Offset > 0 {
		o = append(o, fmt.Sprintf("SKIP %d", query.Offset))
	}

	if query.Limit > 0 {
		o = append(o, fmt.Sprintf("LIMIT %d", query.Limit))
	}

	var where string
	if len(q) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(q, " AND "))
	}
	match := strings.Join(nodes, ", ")
	order := strings.Join(o, " ")

	return match, where, order, f
}

func TranslateFilter(query *Query, q []string, f map[string]interface{}) ([]string, map[string]interface{}) {
	node := strings.ToLower(query.Node)

	for _, filter := range query.Filters {
		field := fmt.Sprintf("%s.%s", node, filter.Field)
		switch filter.Condition {
		case Equal:
			q = append(q, fmt.Sprintf("%s.%s = $`%s`", node, filter.Field, field))
			f[field] = filter.Value
		case In:
			if val, ok := f[field]; ok {
				if v, ok := val.([]interface{}); ok {
					f[field] = append(v, filter.Value)
				}
			} else {
				q = append(q, fmt.Sprintf("%s.%s IN $`%s`", node, filter.Field, field))
				f[field] = []interface{}{filter.Value}
			}
		default:
			q = append(q, fmt.Sprintf("%s.%s = $`%s`", node, filter.Field, field))
			f[field] = filter.Value
		}
	}

	return q, f
}

func RecordUnmarshal(data interface{}, v interface{}) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(res, v); err != nil {
		return err
	}

	return nil
}
