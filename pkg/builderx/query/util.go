package query

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/vldcreation/movie-fest/pkg/timex/civil"
	"github.com/vldcreation/movie-fest/pkg/util"
)

const (
	placeholder          = "?"
	layoutDateTimeFormat = `2006-01-02 15:04:05`
)

var escapedPlaceholder = strings.Repeat(placeholder, 2)

func isTime(obj reflect.Value) bool {
	_, ok := obj.Interface().(time.Time)
	if ok {
		return ok
	}

	_, ok = obj.Interface().(*time.Time)

	return ok
}

func timeIsZero(obj reflect.Value) bool {
	t, ok := obj.Interface().(time.Time)
	if ok {
		return t.IsZero()
	}

	t2, ok := obj.Interface().(*time.Time)
	if ok {
		return false
	}

	return t2 == nil
}

func ToTime(v reflect.Value) time.Time {
	t, ok := v.Interface().(time.Time)
	if ok {
		return t
	}

	t2, ok := v.Interface().(*time.Time)
	if ok {
		if t2 != nil {
			return *t2
		}

	}

	return time.Time{}
}

func ToDate(v reflect.Value) (civil.Date, error) {

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.String {
		dt := v.String()

		if dt == "0000-00-00" {
			return civil.Date{}, nil
		}

		if dt != "" {
			t, err := util.StringToDateE(dt)
			if err != nil {
				return civil.Date{}, err
			}
			return civil.DateOf(t), nil
		}
	}

	return civil.DateOf(ToTime(v)), nil
}

func isNil(i interface{}) bool {
	if i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil()) {
		return true
	}

	return false
}

func PostgrePlaceholder(n int) string {
	return fmt.Sprintf("$%d", n+1)
}

func MsSqlPlaceholder(n int) string {
	return fmt.Sprintf("@p%d", n+1)
}

func ToPostgrePlaceHolder(query string) string {
	b := strings.Builder{}
	n := 0
	for {
		index := strings.Index(query, placeholder)
		if index == -1 {
			break
		}

		// escape placeholder by repeating it twice
		if strings.HasPrefix(query[index:], escapedPlaceholder) {
			b.WriteString(query[:index]) // Write placeholder once, not twice
			query = strings.TrimSpace(query[index+1:])

			continue
		}

		b.WriteString(query[:index])
		b.WriteString(PostgrePlaceholder(n))
		query = query[index+len(placeholder):]
		n++
	}

	// placeholder not found; write remaining query
	b.WriteString(query)

	return b.String()

}

// StructToKeyValue converts a struct to a key value the struct's tags.
// StructToKeyValue uses tags on struct fields to decide which fields to add to the
// returned slice struct.
func StructToKeyValue(src any, tag string) ([]KeyValue, error) {
	var out []KeyValue
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	vals, ok := src.(map[string]any)
	if ok {
		return MapToKeyValue(vals)
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return out, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		if !fi.IsExported() {
			continue
		}

		isPrimary := false

		tagsv := strings.Split(fi.Tag.Get(tag), ",")

		if len(tagsv) == 0 {
			continue
		}

		if len(tagsv) > 1 {
			isPrimary = util.InArray("primary", tagsv[1:])
		}

		if tagsv[0] != "" && fi.PkgPath == "" {

			if isNil(v.Field(i).Interface()) {
				continue
			}

			// skip if omitempty
			if (len(tagsv) > 1 && tagsv[1] == "omitempty") && util.IsEmptyValue(v.Field(i).Interface()) {
				continue
			}

			if isTime(v.Field(i)) {
				if timeIsZero(v.Field(i)) && (len(tagsv) > 1 && tagsv[1] == "omitempty") {
					continue
				}
			}

			if v.Field(i).Kind() == reflect.Struct {
				continue
			}

			col := tagsv[0]

			// set key value of struct key value interface output
			out = append(out, KeyValue{
				Key:       col,
				Value:     v.Field(i).Interface(),
				IsPrimary: isPrimary,
			})
		}

		if tagsv[0] == "" && v.Field(i).Kind() == reflect.Struct {
			x, err := StructToKeyValue(v.Field(i).Interface(), tag)
			if err != nil {
				return out, err
			}

			out = append(out, x...)
		}
	}

	return out, nil
}

// StructToMap converts a struct to a map using the struct's tags.
// StructToMap uses tags on struct fields to decide which fields to add to the
// returned map.
func StructToMap(src any, tag string) (map[string]any, error) {
	out := map[string]any{}
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return out, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		//field := reflectValue.Field(i).Interface()
		if !fi.IsExported() {
			continue
		}

		tagsv := strings.Split(fi.Tag.Get(tag), ",")

		if len(tagsv) == 0 {
			continue
		}

		if tagsv[0] != "" && fi.PkgPath == "" {

			if isNil(v.Field(i).Interface()) {
				continue
			}

			// skip if omitempty
			if (len(tagsv) > 1 && tagsv[1] == "omitempty") && util.IsEmptyValue(v.Field(i).Interface()) {
				continue
			}

			if isTime(v.Field(i)) {
				if timeIsZero(v.Field(i)) && (len(tagsv) > 1 && tagsv[1] == "omitempty") {
					continue
				}
			}

			if v.Field(i).Kind() == reflect.Struct {
				continue
			}

			col := tagsv[0]

			if util.InArray("date", tagsv) {
				d, err := ToDate(v.Field(i))

				if err != nil {
					return out, fmt.Errorf("column %s value %v, %v", col, v.Field(i).String(), err)
				}

				if !d.IsZero() {
					// set value of string slice to value in struct field
					out[col] = d.String()
				}

				continue
			}

			if util.InArray("ne", tagsv) {
				col = fmt.Sprintf("%s !", col)
			}
			// set key value of map interface output
			out[col] = v.Field(i).Interface()
		}

		if tagsv[0] == "" && v.Field(i).Kind() == reflect.Struct {
			x, err := StructToMap(v.Field(i).Interface(), tag)
			if err != nil {
				return out, err
			}

			for y, z := range x {
				out[y] = z
			}
		}
	}

	return out, nil
}

// ToColumnsValues iterate struct to separate key field and value
func ToColumnsValues(src any, tag string) ([]string, []any, error) {
	var columns []string
	var values []any

	vs, ok := src.(map[string]any)
	if ok {
		columns, values = ToColumnValueFromMap(vs)
		return columns, values, nil
	}

	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		tagsv := strings.Split(fi.Tag.Get(tag), ",")

		if len(tagsv) == 0 {
			continue
		}

		if tagsv[0] != "" && fi.PkgPath == "" {

			if tagsv[0] == "-" {
				continue
			}

			if isNil(v.Field(i).Interface()) {
				continue
			}

			// skip if omitempty
			if (len(tagsv) > 1 && util.InArray("omitempty", tagsv)) && util.IsEmptyValue(v.Field(i).Interface()) {
				continue
			}

			if isTime(v.Field(i)) {
				if timeIsZero(v.Field(i)) && (len(tagsv) > 1 && util.InArray("omitempty", tagsv)) {
					continue
				}
			}

			col := tagsv[0]

			if util.InArray("ne", tagsv) {
				col = fmt.Sprintf("%s !", col)
			}

			if util.InArray("date", tagsv) {
				d, err := ToDate(v.Field(i))

				if err != nil {
					return columns, values, fmt.Errorf("column %s value %v, %v", col, v.Field(i).String(), err)
				}

				if !d.IsZero() {
					// set value of string slice to value in struct field
					columns = append(columns, col)
					values = append(values, d.String())
				}

				continue
			}

			// set value of string slice to value in struct field
			columns = append(columns, col)
			// set value interface of value struct field
			values = append(values, v.Field(i).Interface())

		}
	}

	return columns, values, nil
}

// StructToKeyValueWithSkipOmitEmpty converts a struct to a key value the struct's tags.
// StructToKeyValueWithSkipOmitEmpty uses tags on struct fields to decide which fields to add to the
// returned slice struct.
func StructToKeyValueWithSkipOmitEmpty(src any, tag string, columns []string, skipOmitEmpty bool) ([]KeyValue, error) {
	var out []KeyValue
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return out, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		//field := reflectValue.Field(i).Interface()
		if !fi.IsExported() {
			continue
		}

		tagsv := strings.Split(fi.Tag.Get(tag), ",")

		if len(tagsv) == 0 {
			continue
		}

		col := tagsv[0]
		if col != "" && fi.PkgPath == "" {

			if isNil(v.Field(i).Interface()) {
				continue
			}

			if !util.InArray(col, columns) {
				continue
			}

			// skip if omitempty
			if (len(tagsv) > 1 && tagsv[1] == "omitempty") && util.IsEmptyValue(v.Field(i).Interface()) && skipOmitEmpty {
				continue
			}

			if isTime(v.Field(i)) {
				if timeIsZero(v.Field(i)) && (len(tagsv) > 1 && tagsv[1] == "omitempty") && skipOmitEmpty {
					continue
				}
			}

			if v.Field(i).Kind() == reflect.Struct {
				continue
			}

			// set key value of struct key value interface output
			out = append(out, KeyValue{
				Key:   col,
				Value: v.Field(i).Interface(),
			})
		}

		if col == "" && v.Field(i).Kind() == reflect.Struct {
			x, err := StructToKeyValue(v.Field(i).Interface(), tag)
			if err != nil {
				return out, err
			}

			out = append(out, x...)
		}
	}

	return out, nil
}

// ColumnsFromStruct iterate struct to separate key field and value
func ColumnsFromStruct(src any, tag string, skips ...string) ([]string, error) {
	var columns []string

	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		tagsv := strings.Split(fi.Tag.Get(tag), ",")
		if len(tagsv) == 0 {
			continue
		}

		col := tagsv[0]

		if col != "" && fi.PkgPath == "" {

			if col == "-" {
				continue
			}

			if util.InArray(col, skips) {
				continue
			}

			// set value of string slice to value in struct field
			columns = append(columns, col)

		}
	}

	return columns, nil
}

// ColumnsFromStruct iterate struct to separate key field and value
func PrimaryFieldStruct(src any, tag string) (*string, error) {
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)

		tagsv := strings.Split(fi.Tag.Get(tag), ",")

		if len(tagsv) == 0 {
			continue
		}

		if tagsv[0] != "" && fi.PkgPath == "" {

			if tagsv[0] == "-" {
				continue
			}

			if len(tagsv) > 1 {
				if util.InArray("primary", tagsv[1:]) {
					return &tagsv[0], nil
				}
			}
		}
	}

	return nil, fmt.Errorf(`not found column primary from tag %s`, tag)
}

func MapToKeyValue(src any) ([]KeyValue, error) {
	var out []KeyValue

	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	vals, ok := src.(map[string]any)
	if !ok {
		return nil, fmt.Errorf(`invalid source, the source must be map[string]any or map[string]any`)
	}

	for k, v := range vals {
		out = append(out, KeyValue{Key: k, Value: v})
	}

	return out, nil
}

func ToColumnValueFromMap(src map[string]any) ([]string, []any) {

	cols := []string{}
	vals := []any{}

	for k, v := range src {
		vals = append(vals, v)
		cols = append(cols, k)
	}

	return cols, vals
}

func preparedInStatement(start, end int) string {
	var out []string
	for i := start; i <= end; i++ {
		out = append(out, fmt.Sprintf("$%d", i))
	}

	suffix := []byte(",")
	return strings.Join(out, string(suffix))
}

func getLengthOfSlice(v reflect.Value) int {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Slice {
		return v.Len()
	}

	return 0
}

// fieldPtr
type fieldTypePtr struct {
	value any
}

type FieldTypePtr interface {
	StrPtr() *string
	Int64Ptr() *int64
	BoolPtr() *bool
}

func NewFieldTypePtr(v any) FieldTypePtr {
	return &fieldTypePtr{
		value: v,
	}
}

func (f *fieldTypePtr) StrPtr() *string {
	if f.value == nil {
		return nil
	}

	if v, ok := f.value.(string); ok {
		return &v
	}

	return nil
}

func (f *fieldTypePtr) Int64Ptr() *int64 {
	if f.value == nil {
		return nil
	}

	if v, ok := f.value.(int64); ok {
		return &v
	}

	return nil
}

func (f *fieldTypePtr) BoolPtr() *bool {
	if f.value == nil {
		return nil
	}

	if v, ok := f.value.(bool); ok {
		return &v
	}

	return nil
}
