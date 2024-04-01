package db

import (
	"testing"

	configs "github.com/MatthewAraujo/vacation-backend/config"
	"github.com/go-sql-driver/mysql"
)

func TestNewMySQLStorage(t *testing.T) {
	// Test case 1
	cfg := mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	_, err := NewMySQLStorage(cfg)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
