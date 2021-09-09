package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewFolder(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		wantF   Folder
		wantErr bool
	}{
		{"test_0", "alpha:ALPHA:SPARE", Folder{}, true},
		{"test_1", "alpha:ALPHA", Folder{Source: "alpha", Target: "ALPHA"}, false},
		{"test_2", "../alpha/beta", Folder{Source: "../alpha/beta", Target: "beta"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotF, err := NewFolder(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.wantF, gotF)
		})
	}
}
