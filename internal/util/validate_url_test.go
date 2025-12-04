package util

import "testing"

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"http localhost with port", "http://localhost:8083", false},
		{"https localhost", "https://localhost", false},
		{"plain localhost with port", "localhost:8083", false},
		{"domain without scheme", "example.com", false},
		{"domain with scheme", "https://example.com", false},
		{"ip address with port", "127.0.0.1:8083", false},
		{"ip with scheme", "http://127.0.0.1:8083", false},

		{"empty string", "", true},
		{"single letter", "d", true},
		{"no host", "http://", true},
		{"invalid scheme", "ftp://example.com", true},
		{"missing host", "http://:8083", true},
		{"invalid host chars", "http://@@@", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
