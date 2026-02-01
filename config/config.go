package config

import (
	"fmt"
	"os"
	"strconv"

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
}

type Config struct {
	Version      string
	ServiceName  string
	HttpPort     int
	JWTSecretkey string
	DB           *DBConfig
}

func loadConfig() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, relying on environment variables", err)
		os.Exit(1)
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
		fmt.Println("HTTP Port must be a number")
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
	enbleSSLMode,err := strconv.ParseBool(enableSSLMode)
	if err != nil {
		fmt.Println("Enable SSL Mode must be a boolean")
		os.Exit(1)
	}
	configurations = &Config{
		Version:      version,
		ServiceName:  serviceName,
		HttpPort:     port,
		JWTSecretkey: jwtSecretKey,
		DB: &DBConfig{
			Host:          host,
			Port:          dbprt,
			Name:          dbName,
			User:          dbUser,
			Password:      dbPassword,
			EnableSSLMode: enbleSSLMode,
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