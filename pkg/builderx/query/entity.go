package query

import "reflect"

type (
	KeyValue struct {
		Key       string
		Value     any
		IsPrimary bool
	}

	QueryWhere struct {
		Columns []string
		Values  []any
		Limit   int64
		Page    int64
		Query   string
		Q       string
	}

	QueryUpdate struct {
		Columns []string
		Values  []any
		Query   string
	}

	QueryInsert struct {
		Columns []string
		Values  []any
		Query   string
	}

	OrderBy struct {
		Column    string
		OrderType string
	}

	// Custom search for custom query
	// Example: searchColumns = {Type: ">", Col: "created_at", Values: "time1"}
	// valid type: "=", "!=", ">", "<" ,">=", "<="
	// valid column: "created_at"
	// valid value: "time1"
	// @AsColumn is optional, if you want to compare with another column
	// Example: searchColumns = {Type: ">", Col: "qty", Values: "usage", AsColumn: true}
	// That query will be like this: SELECT * FROM table WHERE qty > usage
	CustomSearch struct {
		Type     string
		Col      string
		Value    interface{}
		AsColumn bool
	}

	// SearchIn is a struct to search in or not in
	// Example: searchInColumns = {Type: "in", Col: "column", Values: ["value1", "value2"]}
	// valid type: "in", "not_in"
	// valid value type is slice of string, number or boolean
	SearchIn struct {
		Type   string
		Col    string
		Values []interface{}
	}

	// SearchAny is a struct to search any column
	// Example: 'value' = any(column)
	SearchAny struct {
		Type   string
		Col    string
		Values interface{}
	}

	// Table is a struct to define table name and alias
	// Example: table = {Name: "table_name", Alias: "t", MustAliasColumn: true}
	// valid name: "table_name"
	// valid alias: "t"
	// valid mustAliasColumn: true, false
	// if mustAliasColumn is true, then all column must use alias
	// example: SELECT t.column FROM table_name t
	Table struct {
		Name            string
		Alias           string
		MustAliasColumn bool
	}
)

const (
	ORDER_TYPE_ASC  = "asc"
	ORDER_TYPE_DESC = "desc"

	SEARCH_TYPE_IN    = "in"
	SEARCH_TYPE_NOTIN = "not_in"

	SEARCH_TYPE_EQUAL                 = "="
	SEARCH_TYPE_NOT_EQUAL             = "!="
	SEARCH_TYPE_GREATHER_THAN         = ">"
	SEARCH_TYPE_LOWER_THAN            = "<"
	SEARCH_TYPE_GREATER_THAN_OR_EQUAL = ">="
	SEARCH_TYPE_LOWER_THAN_OR_EQUAL   = "<="
)

func ValidCustomSearchType(searchType string) bool {
	return searchType == SEARCH_TYPE_EQUAL ||
		searchType == SEARCH_TYPE_NOT_EQUAL ||
		searchType == SEARCH_TYPE_GREATHER_THAN ||
		searchType == SEARCH_TYPE_LOWER_THAN ||
		searchType == SEARCH_TYPE_GREATER_THAN_OR_EQUAL ||
		searchType == SEARCH_TYPE_LOWER_THAN_OR_EQUAL
}

func ValidSearchInType(searchType string) bool {
	return searchType == SEARCH_TYPE_IN || searchType == SEARCH_TYPE_NOTIN
}

func ValidSearchAnyType(searchType string) bool {
	return searchType == SEARCH_TYPE_EQUAL ||
		searchType == SEARCH_TYPE_NOT_EQUAL
}

func ValidSearchAnyValues(values interface{}) bool {
	return values != nil ||
		reflect.TypeOf(values).Kind() != reflect.Slice ||
		reflect.TypeOf(values).Kind() != reflect.Array ||
		reflect.TypeOf(values).Kind() != reflect.Struct
}
