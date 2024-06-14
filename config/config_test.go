package config

import "testing"

func TestGetConfig(t *testing.T) {
	serverVal := "sa"
	postgresVal := "pc"
	t.Setenv("SERVER_ADDRESS", serverVal)
	t.Setenv("POSTGRES_CONNECTION", postgresVal)

	cfg := GetConfig()

	if cfg.PostgresConnection != postgresVal {
		t.Errorf("got %s, want %s", cfg.PostgresConnection, postgresVal)
	}

	if cfg.ServerAddress != serverVal {
		t.Errorf("got %s, want %s", cfg.ServerAddress, serverVal)
	}
}
