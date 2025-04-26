package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zjr71163356/simplebank/db/mock"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/utils"
)

type eqCreateUserMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (eqM eqCreateUserMatcher) Matches(x interface{}) bool {
	v, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := utils.MatchPassWord(v.HashedPassword, eqM.password)
	if err != nil {
		return false
	}

	eqM.arg.HashedPassword = v.HashedPassword

	return reflect.DeepEqual(eqM.arg, v)
}

// String describes what the matcher matches.
func (eqM eqCreateUserMatcher) String() string {
	return fmt.Sprintf("is equal to arg: %v (%T) and  password: %v (%T)", eqM.arg, eqM.arg, eqM.password, eqM.password)
}

func NewEqCreateUserMatcher(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserMatcher{
		arg:      arg,
		password: password,
	}
}
func TestGetUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStub     func(*mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStub: func(store *mockdb.MockStore) {
				userInput := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), NewEqCreateUserMatcher(userInput, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStub: func(store *mockdb.MockStore) {

				userInput := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), NewEqCreateUserMatcher(userInput, password)).Times(1).Return(db.User{}, sql.ErrConnDone)
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

				url := "/User/Create"
				body, err := json.Marshal(tc.body)
				require.NoError(t, err)

				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)

				tc.checkResponse(recorder)
			},
		)

	}

}

func randomUser(t *testing.T) (db.User, string) {
	password := utils.RandomString(6)
	HashedPassword, err := utils.HashPassWord(password)
	require.NoError(t, err)
	return db.User{
		Username:       utils.RandomOwnerName(),
		FullName:       utils.RandomOwnerName(),
		Email:          utils.RandomEmail(),
		HashedPassword: HashedPassword,
	}, password
}

func requireBodyMatchUser(t *testing.T, Body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(Body)
	require.NoError(t, err)
	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
