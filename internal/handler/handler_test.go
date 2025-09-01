package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/Zhukek/metrics/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
	}

	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "counter correct",
			request: "/update/counter/someMetric/527",
			method:  http.MethodPost,
			want: want{
				statusCode:  200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "gauge correct",
			request: "/update/gauge/gaugeMetric/527",
			method:  http.MethodPost,
			want: want{
				statusCode:  200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "method incorrect",
			request: "/update/counter/someMetric/527",
			method:  http.MethodGet,
			want: want{
				statusCode:  405,
				contentType: "",
			},
		},
		{
			name:    "shortlink incorrect",
			request: "/update/counter/someMetric",
			method:  http.MethodPost,
			want: want{
				statusCode:  404,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "metric type incorrect",
			request: "/update/wrongType/someMetric/527",
			method:  http.MethodPost,
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "counter value incorrect",
			request: "/update/counter/someMetric/wrongValue",
			method:  http.MethodPost,
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "gauge value incorrect",
			request: "/update/gauge/someMetric/wrongValue",
			method:  http.MethodPost,
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "get metric",
			request: "/value/counter/someMetric",
			method:  http.MethodGet,
			want: want{
				statusCode:  200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "metric name incorrect",
			request: "/value/counter/wrongName",
			method:  http.MethodGet,
			want: want{
				statusCode:  404,
				contentType: "",
			},
		},
		{
			name:    "get metrics list",
			request: "/",
			method:  http.MethodGet,
			want: want{
				statusCode:  200,
				contentType: "text/html; charset=utf-8",
			},
		},
	}

	storage := models.NewStorage()
	server := httptest.NewServer(NewRouter(&storage))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, server.URL+test.request, nil)
			require.NoError(t, err)

			resp, err := server.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, test.want.statusCode, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}
