// This file will run automated tests for API.
package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const ApiUrl = "http://localhost:8000"

func TestApi(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip API tests")
	}

	testcases := getTestCases()
	ctx := context.Background()
	client := &http.Client{}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			for idx := range tc.Steps {
				step := &tc.Steps[idx]
				request, err := step.Request(t, ctx, &tc)
				request.Header.Set("Content-Type", "application/json")
				request.Header.Set("Accept", "application/json")
				require.NoError(t, err)

				// Send request
				response, err := client.Do(request)

				require.NoError(t, err)
				defer response.Body.Close()

				// Check response
				ReadJsonResult(t, response, step)
				step.Expect(t, ctx, &tc, response, step.Result)
			}
		})
	}
}

func getTestCases() []TestCase {
	return []TestCase{
		{
			Name: "Test Health",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest("GET", ApiUrl+"/health", nil)
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
					},
				},
			},
		},
		CreateNormalTestCase("Test Auth Register", []any{AuthRegister, "user1@example.com", "user1", "user1@123"}),
		CreateNormalTestCase("Test Auth Login", []any{AuthLogin, "user1@example.com", "user1@123"}),
	}
}

type TestCase struct {
	Name  string
	Steps []TestCaseStep
}

type RequestFunc func(*testing.T, context.Context, *TestCase) (*http.Request, error)
type ExpectFunc func(*testing.T, context.Context, *TestCase, *http.Response, map[string]any)

type TestCaseStep struct {
	Request RequestFunc
	Expect  ExpectFunc
	Result  map[string]any
}

func ResponseContains(t *testing.T, resp *http.Response, text string) {
	body, err := io.ReadAll(resp.Body)
	bodyStr := string(body)
	require.NoError(t, err)
	require.Contains(t, bodyStr, text)
}

func ReadJsonResult(t *testing.T, resp *http.Response, step *TestCaseStep) {
	var result map[string]any
	err := json.NewDecoder(resp.Body).Decode(&result)
	step.Result = result
	require.NoError(t, err)
}

func RequireIsUUID(t *testing.T, value string) {
	_, err := uuid.Parse(value)
	require.NoError(t, err)
}

func RequireIsValidToken(t *testing.T, testsuite *TestSuite, value string) {
	payload, err := testsuite.TokenMaker.VerifyToken(value)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NoError(t, payload.Valid())
}

const (
	AuthRegister = iota
	AuthLogin
)

func CreateNormalTestCase(name string, a []any) TestCase {
	tc := TestCase{}
	tc.Name = name

	switch a[0].(int) {
	case AuthRegister:
		tc.Steps = append(tc.Steps, TestCaseStep{
			Request: SendRequestRegister(a[1].(string), a[2].(string), a[3].(string)),
			Expect:  ExpectRegisterOk(),
		})
	case AuthLogin:
		tc.Steps = append(tc.Steps, TestCaseStep{
			Request: SendRequestLogin(a[1].(string), a[2].(string)),
			Expect:  ExpectLoginOk(),
		})

	}
	return tc
}

func SendRequestRegister(email string, username string, password string) RequestFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
		body := map[string]any{
			"email":    email,
			"username": username,
			"password": password,
		}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		return http.NewRequest("POST", ApiUrl+"/auth/register", bytes.NewBuffer(jsonBody))
	}
}

func ExpectRegisterOk() ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		RequireIsUUID(t, data["id"].(string))
	}
}

func SendRequestLogin(email string, password string) RequestFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
		body := map[string]any{
			"email":    email,
			"password": password,
		}
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		return http.NewRequest("POST", ApiUrl+"/auth/login", bytes.NewBuffer(jsonBody))
	}
}

func ExpectLoginOk() ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		require.Equal(t, http.StatusOK, resp.StatusCode)
		RequireIsValidToken(t, NewTestSuite(), data["token"].(string))
	}
}

func RequireReturnIsUUID(t *testing.T, resp *http.Response, data map[string]any) {
	require.Equal(t, http.StatusOK, resp.StatusCode)
	RequireIsUUID(t, data["id"].(string))
}

func RequireReturnIsUUIDCreated(t *testing.T, resp *http.Response, data map[string]any) {
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	fmt.Printf("data: %v\n", data)
	RequireIsUUID(t, data["id"].(string))
}

func ExpectBadRequest() ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
}

func ExpectTokenValid(suite *TestSuite) ExpectFunc {
	return func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
		require.Equal(t, http.StatusOK, resp.StatusCode)
		token := data["token"].(string)
		RequireIsValidToken(t, suite, token)
	}
}
