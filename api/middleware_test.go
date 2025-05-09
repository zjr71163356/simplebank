package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/zjr71163356/simplebank/token"
)

func addAuthorization(t *testing.T, request *http.Request, authorizationType string, tokenMaker token.Maker, username string, duration time.Duration) {
	token, _, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	authorizationValue := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationValue)
}

func TestMiddleware(t *testing.T) {

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, "bearer", tokenMaker, "username", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code, "unexpected status code, response body: %s", recorder.Body.String())
			},
		},
		{
			name: "NoAuthorizationHeader",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// Intentionally do nothing to simulate no header
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// Optionally check response body for specific error message
				var errorBody struct {
					Error string `json:"error"`
				}
				err := json.Unmarshal(recorder.Body.Bytes(), &errorBody)
				require.NoError(t, err)
				require.Contains(t, errorBody.Error, "authorization header is not provided")
			},
		},
		{
			name: "UnsupportedAuthorizationType",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// Use a different authorization type like "Basic" or anything other than "Bearer"
				addAuthorization(t, request, "Basic", tokenMaker, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// Optionally check response body
				var errorBody struct {
					Error string `json:"error"`
				}
				err := json.Unmarshal(recorder.Body.Bytes(), &errorBody)
				require.NoError(t, err)
				require.Contains(t, errorBody.Error, "unsupported authorization type") // Adjust expected message based on middleware.go
			},
		},
		{
			name: "InvalidAuthorizationFormatMissingToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// Provide only the type, missing the token part
				request.Header.Set(authorizationHeaderKey, authorizationTypeBearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// Optionally check response body
				var errorBody struct {
					Error string `json:"error"`
				}
				err := json.Unmarshal(recorder.Body.Bytes(), &errorBody)
				require.NoError(t, err)
				require.Contains(t, errorBody.Error, "invalid authorization header format")
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// Create a token with a negative duration, making it instantly expired
				addAuthorization(t, request, authorizationTypeBearer, tokenMaker, "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				// Optionally check response body
				var errorBody struct {
					Error string `json:"error"`
				}
				err := json.Unmarshal(recorder.Body.Bytes(), &errorBody)
				require.NoError(t, err)
				require.Contains(t, errorBody.Error, "token is expired") // Adjust expected message based on token verification logic
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			server, err := newTestServer(t, nil)
			require.NoError(t, err)

			authUrl := "/auth"
			server.router.GET(authUrl, authMiddleWare(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, nil)
			})

			request, err := http.NewRequest(http.MethodGet, authUrl, nil)
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)

			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}

}
