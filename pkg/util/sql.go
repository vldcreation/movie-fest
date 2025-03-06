package util

import "database/sql"

func ToSQLNullableString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func FromSQLNullableString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
