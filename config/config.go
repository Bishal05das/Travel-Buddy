package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lpernett/godotenv"
)

var configurations *Config

type DBConfig struct {
	Host          string
	Port          int
	Name          string
	User          string
	Password      string
	EnableSSLMode bool
	// MigrationsURL is the golang-migrate source, e.g. file://migrations.
	MigrationsURL string
}

type Config struct {
	Version      string
	ServiceName  string
	HttpPort     int
	JWTSecretkey string
	JWTTTL       time.Duration
	// TrustProxyHeaders makes the rate limiter use X-Real-IP /
	// X-Forwarded-For. Enable only behind a reverse proxy that sets them.
	TrustProxyHeaders bool
	// RateLimitPerMinute is the request budget per client IP (default 30).
	RateLimitPerMinute int
	DB                 *DBConfig
}

func loadConfig() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, relying on environment variables")
	}
	version := os.Getenv("VERSION")
	if version == "" {
		fmt.Println("Version is required")
		os.Exit(1)
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		fmt.Println("Service Name is required")
		os.Exit(1)
	}
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		fmt.Println("HTTP Port is required")
		os.Exit(1)
	}
	port, err := strconv.Atoi(httpPort)
	if err != nil {
		fmt.Println("HTTP Port must be a number")
		os.Exit(1)
	}
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		fmt.Println("JWT Secret Key is required")
		os.Exit(1)
	}
	jwtTTL := 24 * time.Hour
	if v := os.Getenv("JWT_TTL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 {
			fmt.Println("JWT_TTL must be a positive duration, e.g. 24h")
			os.Exit(1)
		}
		jwtTTL = d
	}
	trustProxyHeaders := false
	if v := os.Getenv("TRUST_PROXY_HEADERS"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			fmt.Println("TRUST_PROXY_HEADERS must be a boolean")
			os.Exit(1)
		}
		trustProxyHeaders = b
	}
	rateLimit := 0 // middleware default
	if v := os.Getenv("RATE_LIMIT_PER_MINUTE"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			fmt.Println("RATE_LIMIT_PER_MINUTE must be a positive number")
			os.Exit(1)
		}
		rateLimit = n
	}
	host := os.Getenv("DBHOST")
	if host == "" {
		fmt.Println("HOST is required")
		os.Exit(1)
	}
	dbport := os.Getenv("DBPORT")
	if dbport == "" {
		fmt.Println("DBPort is required")
		os.Exit(1)
	}
	dbprt, err := strconv.Atoi(dbport)
	if err != nil {
		fmt.Println("DB Port must be a number")
		os.Exit(1)
	}
	dbName := os.Getenv("DBNAME")
	if dbName == "" {
		fmt.Println("NAME is required")
		os.Exit(1)
	}
	dbUser := os.Getenv("DBUSER")
	if dbUser == "" {
		fmt.Println("User is required")
		os.Exit(1)
	}
	dbPassword := os.Getenv("DBPASSWORD")
	if dbPassword == "" {
		fmt.Println("Password is required")
		os.Exit(1)
	}
	enableSSLMode := os.Getenv("ENABLE_SSL_MODE")
	if enableSSLMode == "" {
		fmt.Println("Enable SSL Mode is required")
		os.Exit(1)
	}
	enbleSSLMode, err := strconv.ParseBool(enableSSLMode)
	if err != nil {
		fmt.Println("Enable SSL Mode must be a boolean")
		os.Exit(1)
	}
	// Relative to the working directory by default so local runs work; the
	// Docker image sets MIGRATIONS_PATH=file:///migrations.
	migrationsURL := os.Getenv("MIGRATIONS_PATH")
	if migrationsURL == "" {
		migrationsURL = "file://migrations"
	}
	configurations = &Config{
		Version:            version,
		ServiceName:        serviceName,
		HttpPort:           port,
		JWTSecretkey:       jwtSecretKey,
		JWTTTL:             jwtTTL,
		TrustProxyHeaders:  trustProxyHeaders,
		RateLimitPerMinute: rateLimit,
		DB: &DBConfig{
			Host:          host,
			Port:          dbprt,
			Name:          dbName,
			User:          dbUser,
			Password:      dbPassword,
			EnableSSLMode: enbleSSLMode,
			MigrationsURL: migrationsURL,
		},
	}
}

func GetConfig() *Config {
	//singleton design pattern
	if configurations == nil {
		loadConfig()
	}

	return configurations
}
