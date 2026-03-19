package url

import "testing"

func TestNewURL(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:      "valid https",
			input:     "https://google.com",
			expected:  "https://google.com",
			expectErr: false,
		},
		{
			name:      "valid http",
			input:     "http://example.com",
			expected:  "http://example.com",
			expectErr: false,
		},
		{
			name:      "missing scheme",
			input:     "google.com",
			expected:  "https://google.com",
			expectErr: false,
		},
		{
			name:      "trim whitespace",
			input:     "   google.com   ",
			expected:  "https://google.com",
			expectErr: false,
		},
		{
			name:      "invalid random string",
			input:     "dsfsdf",
			expectErr: true,
		},
		{
			name:      "missing host",
			input:     "https://",
			expectErr: true,
		},
		{
			name:      "unsupported scheme",
			input:     "ftp://example.com",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			u, err := NewUrl(tt.input)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if u.Value != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, u.Value)
			}
		})
	}
}
