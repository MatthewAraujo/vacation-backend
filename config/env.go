package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type CloudflareConfig struct {
	BucketName      string
	AccountID       string
	AccessKeyID     string
	AccessKeySecret string
}
type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	Cloudflare             CloudflareConfig
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "ecom"),
		JWTSecret:              getEnv("JWT_SECRET", "not-that-secret"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
		Cloudflare: CloudflareConfig{
			BucketName:      getEnv("CLOUDFLARE_BUCKETNAME", "my_bucket"),
			AccountID:       getEnv("CLOUDFLARE_ACCOUNT_ID", "my_account_id"),
			AccessKeyID:     getEnv("CLOUDFLARE_ACESS_KEY_ID", "my_access_key_id"),
			AccessKeySecret: getEnv("CLOUDFLARE_ACESS_SECRET_ID", "my_access_key_secret"),
		},
	}

}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
