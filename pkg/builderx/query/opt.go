package query

type (
	// BuilderxOpt is a struct to store options while building query
	// Note: This approach only works for postgresql implementation
	// TODO: Add support for other database
	BuilderxOpt struct {
		// Query search use Q as search query param
		// searchColumns is a list of columns that will be used for searching
		// Example: searchColumns = ["name", "description"]
		SearchColumns []string
		// OrderBy is a struct to store order by column and order type
		// Example: orderBy = {Column: "name", OrderType: "asc"}
		OrderBy OrderBy
		// GroupBy is a string to store group by column
		// Example: groupBy = "name"
		GroupBy []string
		// WithPagination is a boolean to determine whether to use pagination or not
		// Example: withPagination = true
		WithPagination bool
		// SearchInColumns is a list of columns-i idx that will be used for searching with IN operator
		// Example: searchInColumns = {Type: "in", Col: "column", Values: ["value1", "value2"]}
		// query will be like this: SELECT * FROM table WHERE column IN (value1, value2)
		SearchInColumns *SearchIn
		// CustomSearch is a list of columns that will be used for searching with BETWEEN operator
		// Example: searchBetweenColumns = {Col: "created_at", Values: "time1"}
		// valid type: "=", "!=", ">", "<" ,">=", "<="
		CustomSearch *CustomSearch
		// MultipleCustomSearch is a list of CustomSearch
		MultipleCustomSearch []CustomSearch
		// SearchAnyColumn is a string to store search any column
		// Example: 'value' = any(column)
		SearchAnyColumn *SearchAny
		// deletedAtColumn is a function to determine whether to include deleted data or not
		// Example: deletedAtColumn = collected_datas.deleted_at
		DeletedAtColumn *string
		// Table is a struct to define table name and alias
		// Example: table = {Name: "table_name", Alias: "t", MustAliasColumn: true}
		Table *Table
	}
)

// NewBuilderxOpt is a function to create new BuilderxOpt that will be used in builderx
func NewBuilderxOpt(opts ...funcBuilderxOpt) *BuilderxOpt {
	b := &BuilderxOpt{}

	// Set default value
	// We assume that we always use pagination
	b.WithPagination = true

	b.AddOptions(opts...)

	return b
}

// add options
func (b *BuilderxOpt) AddOptions(opts ...funcBuilderxOpt) {
	for _, opt := range opts {
		opt(b)
	}
}

type funcBuilderxOpt func(*BuilderxOpt)

func WithSearchColumns(columns ...string) funcBuilderxOpt {
	return func(b *BuilderxOpt) {
		b.SearchColumns = columns
	}
}

func WithOrderBy(column string, orderType string) funcBuilderxOpt {
	return func(b *BuilderxOpt) {
		b.OrderBy = OrderBy{
			Column:    column,
			OrderType: orderType,
		}
	}
}

func WithGroupBy(columns ...string) funcBuilderxOpt {
	return func(b *BuilderxOpt) {
		b.GroupBy = columns
	}
}

func WithPagination(v bool) funcBuilderxOpt {
	return func(b *BuilderxOpt) {
		b.WithPagination = v
	}
}

func WithSearchInColumns(search SearchIn) funcBuilderxOpt {
	return func(b *BuilderxOpt) {
		if !ValidSearchInType(search.Type) {
			search.Type = SEARCH_TYPE_IN
		}

		if search.Col == "" {
			return
		}

		if search.Values == nil {
			return
		}

		b.SearchInColumns = &SearchIn{
			Type:   search.Type,
			Col:    search.Col,
			Values: search.Values,
		}
	}
}

func WithCustomSearch(custom CustomSearch) funcBuilderxOpt {
	return func(b *BuilderxOpt) {
		if !ValidCustomSearchType(custom.Type) {
			custom.Type = SEARCH_TYPE_EQUAL
		}

		b.CustomSearch = &CustomSearch{
			Type:     custom.Type,
			Col:      custom.Col,
			Value:    custom.Value,
			AsColumn: custom.AsColumn,
		}
	}
}

func WithMultipleCustomSearch(customs []CustomSearch) funcBuilderxOpt {
	return func(b *BuilderxOpt) {
		for _, custom := range customs {
			if !ValidCustomSearchType(custom.Type) {
				custom.Type = SEARCH_TYPE_EQUAL
			}

			b.MultipleCustomSearch = append(b.MultipleCustomSearch, CustomSearch{
				Type:     custom.Type,
				Col:      custom.Col,
				Value:    custom.Value,
				AsColumn: custom.AsColumn,
			})
		}
	}
}

func WithSearchAnyColumn(search SearchAny) funcBuilderxOpt {
	return func(b *BuilderxOpt) {
		if !ValidSearchAnyType(search.Type) {
			search.Type = SEARCH_TYPE_EQUAL
		}

		if !ValidSearchAnyValues(search.Values) {
			return
		}

		b.SearchAnyColumn = &search
	}
}

func SetDeletedAtColumn(column string) funcBuilderxOpt {
	return func(b *BuilderxOpt) {
		b.DeletedAtColumn = &column
	}
}

func SetTable(table Table) funcBuilderxOpt {
	return func(b *BuilderxOpt) {
		b.Table = &table
	}
}
