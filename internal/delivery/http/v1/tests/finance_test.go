package tests

import (
	"bytes"
	"net/http/httptest"
	"testing"

	mock_redis "github.com/cookienyancloud/avito-backend-test/internal/cache/redis/mocks"
	"github.com/cookienyancloud/avito-backend-test/internal/delivery/http/v1"
	"github.com/cookienyancloud/avito-backend-test/internal/domain"
	mock_service "github.com/cookienyancloud/avito-backend-test/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {

	type mockBehavior func(s *mock_service.MockIService, inp *domain.TransactionInput)

	tt := []struct {
		name          string
		inpBody       string
		input         *domain.TransactionInput
		mockB         mockBehavior
		expStatusCode int
		expReqBody    string
	}{
		{
			name:    "ok",
			inpBody: `{"id":"1993f8f2-d580-4fb1-bd8e-5bdfb7ddd7e4","sum":10,"description":"test ok"}`,
			input: &domain.TransactionInput{
				Id:          uuid.MustParse("1993f8f2-d580-4fb1-bd8e-5bdfb7ddd7e4"),
				Sum:         10,
				Description: "test ok",
			},
			mockB: func(s *mock_service.MockIService, inp *domain.TransactionInput) {
				s.
					EXPECT().
					MakeTransaction(gomock.Any(), inp).
					Return(nil)
			},
			expStatusCode: 200,
			expReqBody:    `{"message":"удачная транзакция"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			finService := mock_service.NewMockIService(c)
			subService := mock_service.NewMockICurrency(c)
			cache := mock_redis.NewMockICache(c)
			w := httptest.NewRecorder()
			r := gin.New()
			tc.mockB(finService, tc.input)
			handler := v1.NewHandler(finService, subService, cache)
			r.POST("/transaction", handler.Transaction)
			req := httptest.NewRequest("POST", "/transaction",
				bytes.NewBufferString(tc.inpBody))
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.expStatusCode, w.Code)
			assert.Equal(t, tc.expReqBody, w.Body.String())

		})
	}

}