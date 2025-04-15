package api

import (
	"fmt"
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

	ctrl := gomock.NewController(t)

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	server := NewServer(store)

	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/GetAccount/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, nil, err)

	server.router.ServeHTTP(recorder, request)

}

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt63(1, 100),
		Owner:    utils.RandomOwnerName(1, 5),
		Balance:  utils.RandomInt63(1, 10000),
		Currency: utils.RandomCurrency(),
	}
}
