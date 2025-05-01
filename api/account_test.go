package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zjr71163356/simplebank/db/mock"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/token"
	"github.com/zjr71163356/simplebank/utils"
)

func TestGetAccount(t *testing.T) {
	username := utils.RandomOwnerName()
	account := randomAccount(username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(*mockdb.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, authorizationTypeBearer, tokenMaker, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, authorizationTypeBearer, tokenMaker, "unauthorized_user", time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code, "unexpected status code, response body: %s", recorder.Body.String())
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, authorizationTypeBearer, tokenMaker, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, authorizationTypeBearer, tokenMaker, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name:      "InvalidID",
			accountID: -1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, authorizationTypeBearer, tokenMaker, account.Owner, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(
			tc.name,
			func(t *testing.T) {
				ctrl := gomock.NewController(t)

				store := mockdb.NewMockStore(ctrl)
				tc.buildStub(store)

				server, err := newTestServer(t, store)

				require.NoError(t, err)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/GetAccount/%d", tc.accountID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				tc.setupAuth(t, request, server.tokenMaker)
				server.router.ServeHTTP(recorder, request)

				tc.checkResponse(t, recorder)
			},
		)

	}

}

func TestCreateAccount(t *testing.T) {
	username := utils.RandomOwnerName()
	account := randomAccount(username)
	CreateAccount := db.CreateAccountParams{
		Owner:    account.Owner,
		Currency: account.Currency,
		Balance:  0,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStub     func(*mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
				"balance":  0,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(CreateAccount)).Times(1).Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
				"balance":  0,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(CreateAccount)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(
			tc.name,
			func(t *testing.T) {
				ctrl := gomock.NewController(t)

				store := mockdb.NewMockStore(ctrl)
				tc.buildStub(store)

				server, err := newTestServer(t, store)
				require.NoError(t, err)
				recorder := httptest.NewRecorder()

				url := "/CreateAccount"
				// Marshal body data to JSON
				data, err := json.Marshal(tc.body)
				require.NoError(t, err)

				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)

				tc.checkResponse(recorder)
			},
		)

	}
}

func randomAccount(username string) db.Account {
	return db.Account{
		ID:       utils.RandomInt63(1, 100),
		Owner:    username,
		Balance:  utils.RandomInt63(1, 10000),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, Body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(Body)
	require.NoError(t, err)
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
