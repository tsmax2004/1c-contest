package utils

import (
	"os"
	"path/filepath"
	"testing"

    "github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

func TestGetRelPath(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tests := []struct {
		name         string
		inputPath    string
		expectedRel  string
		expectingErr bool
	}{
		{
			name:         "Relative Path",
			inputPath:    filepath.Join(pwd, "test", "file.txt"),
			expectedRel:  filepath.Join("test", "file.txt"),
			expectingErr: false,
		},
		{
			name:         "Same Directory",
			inputPath:    pwd,
			expectedRel:  ".",
			expectingErr: false,
		},
	}

	for _, tt := range tests {
		runner.Run(t, tt.name, func(t provider.T) {
			rel, err := GetRelPath(tt.inputPath)
			if (err != nil) != tt.expectingErr {
				t.Errorf("GetRelPath() error = %v, expectingErr %v", err, tt.expectingErr)
				return
			}
			if rel != tt.expectedRel {
				t.Errorf("GetRelPath() got = %v, expected %v", rel, tt.expectedRel)
			}
		})
	}
}
