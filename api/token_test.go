package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zjr71163356/simplebank/db/mock"
	db "github.com/zjr71163356/simplebank/db/sqlc"
	"github.com/zjr71163356/simplebank/token"
	"github.com/zjr71163356/simplebank/utils"
)

// 创建随机用户用于令牌测试

func TestRenewTokenAPI(t *testing.T) {
	user, _ := randomUser(t)
	accessTokenDuration := time.Minute

	// 创建随机会话的辅助函数
	randomSession := func(refreshToken string, username string, isBlocked bool, expiresAt time.Time) db.Session {
		return db.Session{
			ID:           uuid.New(),
			Username:     username,
			RefreshToken: refreshToken,
			UserAgent:    utils.RandomString(10),
			ClientIp:     utils.RandomString(10), // 使用 RandomString 替代 RandomIPV4
			IsBlocked:    isBlocked,
			ExpiresAt:    expiresAt,
			CreatedAt:    time.Now(),
		}
	}

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, tokenMaker token.Maker) string // 返回用于请求的刷新令牌
		buildStubs    func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload, session db.Session)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{}, // 刷新令牌将由 setupAuth 设置
			setupAuth: func(t *testing.T, tokenMaker token.Maker) string {
				refreshToken, _, err := tokenMaker.CreateToken(user.Username, time.Hour) // 刷新令牌持续时间较长
				require.NoError(t, err)
				require.NotEmpty(t, refreshToken)
				return refreshToken
			},
			buildStubs: func(store *mockdb.MockStore, refreshTokenString string, refreshPayload *token.Payload, session db.Session) {
				// 会话与刷新令牌和载荷匹配
				validSession := randomSession(refreshTokenString, refreshPayload.Username, false, time.Now().Add(time.Hour))

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.Id)).
					Times(1).
					Return(validSession, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusAccepted, recorder.Code)

				var rsp NewAccessTokenResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &rsp)
				require.NoError(t, err)
				require.NotEmpty(t, rsp.AccessToken)
				require.WithinDuration(t, time.Now().Add(accessTokenDuration), rsp.AccessTokenExpiresAt, time.Second)
			},
		},
		{
			name: "BadRequestMissingRefreshToken",
			body: gin.H{}, // 故意为空
			setupAuth: func(t *testing.T, tokenMaker token.Maker) string {
				return "" // 无刷新令牌
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload, session db.Session) {
				// 不期望调用
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UnauthorizedInvalidRefreshToken",
			body: gin.H{"refresh_token": "invalid-token"},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) string {
				return "invalid-token"
			},
			buildStubs: func(store *mockdb.MockStore, refreshToken string, refreshPayload *token.Payload, session db.Session) {
				// 不期望 GetSession 被调用，因为令牌验证会失败
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "SessionNotFound", // 服务器当前对任何 GetSession 错误返回 500
			body: gin.H{},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) string {
				refreshToken, _, err := tokenMaker.CreateToken(user.Username, time.Hour)
				require.NoError(t, err)
				return refreshToken
			},
			buildStubs: func(store *mockdb.MockStore, refreshTokenString string, refreshPayload *token.Payload, session db.Session) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.Id)).
					Times(1).
					Return(db.Session{}, sql.ErrNoRows) // 或 db.ErrRecordNotFound（如果已定义）
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code) // 根据当前 token.go，任何 GetSession 错误都返回 500
			},
		},
		{
			name: "SessionRefreshTokenMismatch",
			body: gin.H{},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) string {
				refreshToken, _, err := tokenMaker.CreateToken(user.Username, time.Hour)
				require.NoError(t, err)
				return refreshToken
			},
			buildStubs: func(store *mockdb.MockStore, refreshTokenString string, refreshPayload *token.Payload, session db.Session) {
				mismatchedSession := randomSession("different-refresh-token", refreshPayload.Username, false, time.Now().Add(time.Hour))
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.Id)).
					Times(1).
					Return(mismatchedSession, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				require.Contains(t, recorder.Body.String(), "session refresh token does not match request")
			},
		},
		{
			name: "SessionUserMismatch",
			body: gin.H{},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) string {
				refreshToken, _, err := tokenMaker.CreateToken(user.Username, time.Hour) // 用户的令牌
				require.NoError(t, err)
				return refreshToken
			},
			buildStubs: func(store *mockdb.MockStore, refreshTokenString string, refreshPayload *token.Payload, session db.Session) {
				// 会话属于不同的用户
				mismatchedUserSession := randomSession(refreshTokenString, "anotheruser", false, time.Now().Add(time.Hour))
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.Id)).
					Times(1).
					Return(mismatchedUserSession, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				require.Contains(t, recorder.Body.String(), "session user does not match token user")
			},
		},
		{
			name: "SessionIsBlocked",
			body: gin.H{},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) string {
				refreshToken, _, err := tokenMaker.CreateToken(user.Username, time.Hour)
				require.NoError(t, err)
				return refreshToken
			},
			buildStubs: func(store *mockdb.MockStore, refreshTokenString string, refreshPayload *token.Payload, session db.Session) {
				blockedSession := randomSession(refreshTokenString, refreshPayload.Username, true, time.Now().Add(time.Hour)) // IsBlocked = true
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.Id)).
					Times(1).
					Return(blockedSession, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				require.Contains(t, recorder.Body.String(), "session is blocked")
			},
		},
		{
			name: "SessionIsExpired",
			body: gin.H{},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) string {
				refreshToken, _, err := tokenMaker.CreateToken(user.Username, time.Hour)
				require.NoError(t, err)
				return refreshToken
			},
			buildStubs: func(store *mockdb.MockStore, refreshTokenString string, refreshPayload *token.Payload, session db.Session) {
				expiredSession := randomSession(refreshTokenString, refreshPayload.Username, false, time.Now().Add(-time.Hour)) // ExpiresAt 在过去
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(refreshPayload.Id)).
					Times(1).
					Return(expiredSession, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				require.Contains(t, recorder.Body.String(), "session is expired")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// 创建一个真实的 PasetoMaker 用于生成初始刷新令牌
			// 使用固定密钥以便可重现，或者使用随机密钥以保证唯一性
			symmetricKey := utils.RandomString(32)
			tokenMaker, err := token.NewPasetoMaker(symmetricKey)
			require.NoError(t, err)

			// 使用真实的 token maker 获取刷新令牌字符串（因为 setupAuth 需要 token.Maker）
			var refreshTokenString string
			if tc.setupAuth != nil {
				refreshTokenString = tc.setupAuth(t, tokenMaker)
				if refreshTokenString != "" {
					if _, ok := tc.body["refresh_token"]; !ok && len(tc.body) > 0 {
						tc.body["refresh_token"] = refreshTokenString
					} else if refreshTokenString != "" {
						tc.body = gin.H{"refresh_token": refreshTokenString}
					}
				}
			}

			// 如果生成了刷新令牌，解析它以获取其载荷用于模拟设置
			var refreshPayload *token.Payload
			if refreshTokenString != "" && tc.name != "UnauthorizedInvalidRefreshToken" {
				refreshPayload, err = tokenMaker.VerifyToken(refreshTokenString)
				require.NoError(t, err)
			}

			tc.buildStubs(store, refreshTokenString, refreshPayload, db.Session{})

			// 创建测试服务器
			// 假设 newTestServer 定义在 main_test.go 或类似文件中并返回 (*Server, error)
			server, err := newTestServer(t, store)
			require.NoError(t, err)

			// 替换 tokenMaker 以便我们可以控制它的行为
			server.tokenMaker = tokenMaker

			// 确保 AccessTokenDuration 与测试设置一致
			server.config.AccessTokenDuration = accessTokenDuration

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/RenewToken" // 确保与 server.go 中的路由匹配
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
