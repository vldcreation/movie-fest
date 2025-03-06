package query

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/spf13/cast"
)

func PostgreQueryWhereV2(elems []string, values *[]any, opt *BuilderxOpt, start int) string {
	b := new(strings.Builder)

	// shallow copy values
	_values := *values

	for i := 0; i < len(elems); i++ {
		// typElems := reflect.TypeOf(values[i]).Kind()
		switch _values[i].(type) {
		case string, *string:
			b.WriteString(fmt.Sprintf("LOWER(%s::varchar) = LOWER($%d) AND ", elems[i], i+start))
		default:
			b.WriteString(fmt.Sprintf("%s = $%d AND ", elems[i], i+start))
		}
	}

	// update start
	start = len(elems) + start

	if opt != nil && opt.SearchInColumns != nil {
		// get len values
		end := len(opt.SearchInColumns.Values) - 1 + start

		switch opt.SearchInColumns.Type {
		case SEARCH_TYPE_IN:
			b.WriteString(fmt.Sprintf("%s IN (%s) AND ", opt.SearchInColumns.Col, preparedInStatement(start, end)))

		case SEARCH_TYPE_NOTIN:
			b.WriteString(fmt.Sprintf("%s NOT IN (%s) AND ", opt.SearchInColumns.Col, preparedInStatement(start, end)))

		default:
			b.WriteString(fmt.Sprintf("%s IN (%s) AND ", opt.SearchInColumns.Col, preparedInStatement(start, end)))
		}

		// add values to the query
		for i := 0; i < len(opt.SearchInColumns.Values); i++ {
			if opt.SearchInColumns.Values[i] == nil {
				continue
			}

			*values = append(*values, opt.SearchInColumns.Values[i])
		}

		// update start
		start = end + 1
	}

	if opt != nil && opt.CustomSearch != nil {
		if opt.CustomSearch.AsColumn {
			b.WriteString(fmt.Sprintf("%s %s %s AND ", opt.CustomSearch.Col, opt.CustomSearch.Type, opt.CustomSearch.Value))
		} else {
			b.WriteString(fmt.Sprintf("%s %s $%d AND ", opt.CustomSearch.Col, opt.CustomSearch.Type, start))

			// add values to the query
			*values = append(*values, opt.CustomSearch.Value)

			// update start
			start = start + 1
		}
	}

	if opt != nil && opt.MultipleCustomSearch != nil {
		for _, custom := range opt.MultipleCustomSearch {
			if custom.AsColumn {
				b.WriteString(fmt.Sprintf("%s %s %s AND ", custom.Col, custom.Type, custom.Value))
				continue
			} else {
				b.WriteString(fmt.Sprintf("%s %s $%d AND ", custom.Col, custom.Type, start))

				// add values to the query
				*values = append(*values, custom.Value)

				// update start
				start = start + 1
			}
		}

	}

	if opt != nil && opt.SearchAnyColumn != nil {
		b.WriteString(fmt.Sprintf("$%d %s ANY(\"%s\") AND ", start, opt.SearchAnyColumn.Type, opt.SearchAnyColumn.Col))

		// add values to the query
		*values = append(*values, opt.SearchAnyColumn.Values)

		// update start
		start = start + 1
	}

	if b.Len() == 0 {
		return b.String()
	}

	return b.String()[0 : b.Len()-4]
}

func StructToPostgreQueryWhereWithOpt(iStruct interface{}, tag string, option *BuilderxOpt) (QueryWhere, error) {
	var (
		periodRange string
		startDate,
		endDate string
		qw = QueryWhere{}
	)

	if iStruct == nil {
		return qw, nil
	}

	data, err := StructToKeyValue(iStruct, tag)
	if err != nil {
		return qw, err
	}

	if len(data) == 0 {
		return qw, err
	}

	if option == nil {
		option = &BuilderxOpt{}
	}

	for i := 0; i < len(data); i++ {
		if data[i].Key == "q" {
			qw.Q = cast.ToString(data[i].Value)
			continue
		}

		if data[i].Key == "page" {
			qw.Page = cast.ToInt64(data[i].Value)
			continue
		}

		if data[i].Key == "per_page" {
			qw.Limit = cast.ToInt64(data[i].Value)
			continue
		}

		if data[i].Key == "start_date" {
			startDate = cast.ToString(data[i].Value)
			continue
		}

		if data[i].Key == "end_date" {
			endDate = cast.ToString(data[i].Value)
			continue
		}

		// check if table alias is exists and should be added to each column
		if option.Table != nil && option.Table.Alias != "" && option.Table.MustAliasColumn {
			data[i].Key = fmt.Sprintf("%s.%s", option.Table.Alias, data[i].Key)
		}

		// set values to pq.Array for array data type
		switch data[i].Value.(type) {
		case []string, []*string, []int, []int64, []float64:
			qw.Values = append(qw.Values, pq.Array(data[i].Value))
		default:
			qw.Values = append(qw.Values, data[i].Value)
		}

		qw.Columns = append(qw.Columns, data[i].Key)
	}

	if (len(startDate) > 0 && len(endDate) == 0) || (len(startDate) == 0 && len(endDate) > 0) {
		return qw, fmt.Errorf("invalid date period start %s end %s", startDate, endDate)
	}

	if len(data) == 0 && (len(startDate) == 0 && len(endDate) == 0) {
		return qw, fmt.Errorf("the struct is empty value")
	}

	nw := PostgreQueryWhereV2(qw.Columns, &qw.Values, option, 1)

	nc := len(qw.Columns)
	if len(startDate) > 0 && len(endDate) > 0 {
		qw.Values = append(qw.Values, startDate, endDate)
		periodRange = "(created_at >= $1  AND created_at <= $2 )"
		if len(data) > 0 {
			periodRange = fmt.Sprintf("(created_at >= $%d  AND created_at <= $%d )", nc+1, nc+2)
		}
	}

	if len(nw) > 0 {
		qw.Query = fmt.Sprintf(`WHERE %s`, nw)
	}

	if len(periodRange) > 0 && len(nw) > 0 {
		qw.Query = fmt.Sprintf(`WHERE %s AND %s`, nw, periodRange)
	}

	if qw.Q != "" && option.SearchColumns != nil {
		searchQuery := make([]string, 0)
		for _, v := range option.SearchColumns {
			sanitizedString := strings.ReplaceAll(qw.Q, "'", "")
			searchQuery = append(searchQuery, fmt.Sprintf("LOWER(%s) ILIKE LOWER('%%%s%%')", v, sanitizedString))
		}
		if len(qw.Query) > 0 {
			qw.Query += fmt.Sprintf(` AND (%s)`, strings.Join(searchQuery, " OR "))
		} else {
			qw.Query += fmt.Sprintf(` WHERE (%s)`, strings.Join(searchQuery, " OR "))
		}
	}

	if option.DeletedAtColumn != nil {
		if qw.Query == "" {
			qw.Query += fmt.Sprintf(`WHERE %s IS NULL`, *option.DeletedAtColumn)
		} else {
			qw.Query += fmt.Sprintf(` AND %s IS NULL`, *option.DeletedAtColumn)
		}
	} else {
		if qw.Query == "" {
			qw.Query += `WHERE deleted_at IS NULL`
		} else {
			qw.Query += ` AND deleted_at IS NULL`
		}
	}

	if option.GroupBy != nil && len(option.GroupBy) > 0 {
		qw.Query += fmt.Sprintf(` GROUP BY %s`, strings.Join(option.GroupBy, ", "))
	}

	if option.OrderBy.Column != "" {
		qw.Query += fmt.Sprintf(` ORDER BY %s %s`, option.OrderBy.Column, option.OrderBy.OrderType)
	}

	if option.WithPagination {
		if qw.Limit > 0 {
			qw.Query += fmt.Sprintf(` LIMIT %d`, qw.Limit)
		}

		if qw.Page > 0 {
			qw.Query += fmt.Sprintf(` OFFSET %d`, (qw.Page-1)*qw.Limit)
		}
	}

	return qw, err
}
