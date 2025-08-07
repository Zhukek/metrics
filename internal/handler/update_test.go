package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
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
				contentType: "text/plain; charset=utf-8",
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			recorder := httptest.NewRecorder()
			Update(recorder, request)
			result := recorder.Result()
			result.Body.Close()
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}
