package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthUser(t *testing.T) {
	nextHandler := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}

	testCases := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{
			name:         "Correct token",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5MjgyMjc2MzksImlzcyI6InRvZG9saXN0YXBwIiwidHlwZSI6ImFjY2VzcyIsInVpZCI6MX0.eWYV4i1unAJmQocgBk6zSN5XeVMcVkQdlF4qJDyWH1Y",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Incorrect token",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInR5cGUiOiJhY2Nlc3MiLCJleHAiOjE3MjgyMjc2Mzl9.urHS-cdM5rzK9qO1T_YGgiD0-CdiG-DooOsjTr8RPAE",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "No token",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Expired token",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInR5cGUiOiJhY2Nlc3MiLCJleHAiOjE0MjgyMjc2Mzl9.rpo908rOvCa6kEN9LNtIT0KiHrG3WOTZ7jRAQiuIaTQ",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Token with invalid claims",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoiYWNjZXNzIiwiZXhwIjoxNDI4MjI3NjM5fQ.1JjIhQdv8B2vxEGeikpuKbXxIcPZnEA31ceetzKSHVo",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(AuthUserMiddleware(nextHandler))
			defer server.Close()

			req, _ := http.NewRequest("GET", server.URL, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.token))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Errorf("expected status code %d, but got: %d", test.expectedCode, resp.StatusCode)
			}
		})
	}

}
