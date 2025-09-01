package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name string
		URL  string
		want string
	}{
		{
			name: "http protocol",
			URL:  "http://localhost:3000",
			want: "http://localhost:3000",
		},
		{
			name: "https protocol",
			URL:  "https://localhost:3000",
			want: "https://localhost:3000",
		},
		{
			name: "without protocol",
			URL:  "localhost:3000",
			want: "http://localhost:3000",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, GetBaseURL(test.URL))
		})
	}
}
