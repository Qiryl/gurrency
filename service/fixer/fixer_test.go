package fixer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Qiryl/gurrency"
	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	t.Parallel()
	testTable := []struct {
		name           string
		serverResponse *fixerResponse
		expectedRes    *gurrency.CurrencyRate
		expectedError  error
		mockBehavior   func(t *testing.T, sRes *fixerResponse) http.HandlerFunc
	}{
		{
			name:        "success",
			expectedRes: &gurrency.CurrencyRate{ServiceName: "fixer.io"},
			mockBehavior: func(t *testing.T, res *fixerResponse) http.HandlerFunc {
				byte, err := json.Marshal(res)
				assert.NoError(t, err)
				return func(rw http.ResponseWriter, r *http.Request) {
					_, err := rw.Write(byte)
					assert.NoError(t, err)
				}
			},
		},
		{
			name:          "invalid json",
			expectedError: ErrInvalidResponse,
			mockBehavior: func(t *testing.T, sRes *fixerResponse) http.HandlerFunc {
				byte := []byte("garbage")
				return func(rw http.ResponseWriter, r *http.Request) {
					_, err := rw.Write(byte)
					assert.NoError(t, err)
				}
			},
		},
		{
			name:          "bad response",
			expectedError: ErrBadResponse,
			mockBehavior: func(t *testing.T, sRes *fixerResponse) http.HandlerFunc {
				return func(rw http.ResponseWriter, r *http.Request) {
					rw.WriteHeader(http.StatusInternalServerError)
				}
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			src := httptest.NewServer(tc.mockBehavior(t, tc.serverResponse))
			defer src.Close()

			fs := NewFixerSource("test", src.URL, "test", "test")
			res, err := fs.GetRate()
			assert.Equal(t, tc.expectedRes, res)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
