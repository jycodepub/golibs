package config

import (
	"path/filepath"
	"reflect"
	"testing"
)

type TestConfig struct {
	Key  string   `json:"key"`
	List []string `json:"list"`
	Num  int      `json:"num"`
}

func TestLoad_Success(t *testing.T) {
	// Relative path to test/config.json from config package
	jsonPath := filepath.Join("..", "test", "config.json")

	var cfg TestConfig
	err := Load(jsonPath, &cfg)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if cfg.Key != "value" {
		t.Errorf("expected Key to be 'value', got '%s'", cfg.Key)
	}

	expectedList := []string{"v1", "v2"}
	if !reflect.DeepEqual(cfg.List, expectedList) {
		t.Errorf("expected List to be %v, got %v", expectedList, cfg.List)
	}

	if cfg.Num != 1 {
		t.Errorf("expected Num to be 1, got %d", cfg.Num)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	var cfg TestConfig
	err := Load("non_existent_file.json", &cfg)
	if err == nil {
		t.Fatalf("expected file not found error, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	// Loading non-JSON file (e.g. config.go) to trigger unmarshal error
	var cfg TestConfig
	err := Load("config.go", &cfg)
	if err == nil {
		t.Fatalf("expected unmarshal error, got nil")
	}
}
