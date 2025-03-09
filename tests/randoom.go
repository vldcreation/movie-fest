package tests

import "github.com/vldcreation/movie-fest/pkg/util"

// GenerateRandoomRegisterRequest return []any{"email", "username", "password"}
func GenerateRandoomRegisterRequest() []any {
	return []any{
		util.RandEmail(),
		util.RandString(7),
		util.RandString(10),
	}
}

// GenerateRandoomLoginRequest return []any{"email", "password"}
func GenerateRandoomLoginRequest() []any {
	return []any{
		util.RandEmail(),
		util.RandString(10),
	}
}
