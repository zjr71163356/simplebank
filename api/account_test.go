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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zjr71163356/simplebank/db/mock"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/utils"
)

func TestAccountAPI(t *testing.T) {
	account := randomAccount()
	testCases := []struct {
		name          string
		accountID     int64
		buildStub     func(*mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, db.Account{})
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name:      "InvalidID",
			accountID: -1,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, db.Account{})
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
				require.NoError(t, nil, err)

				server.router.ServeHTTP(recorder, request)

				tc.checkResponse(recorder)
			},
		)

	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt63(1, 100),
		Owner:    utils.RandomOwnerName(),
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
