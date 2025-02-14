package config

import (
    "log"
    "os"
)

type Config struct {
    MongoURI  string
    DBName    string
    JWTSecret string
    SMTPHost  string
    SMTPPort  string
    SMTPEmail string
    SMTPPass  string
}

func LoadConfig() *Config {
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Fatal("MONGO_URI is not set in environment variables")
    }

    dbName := os.Getenv("MONGO_DB")
    if dbName == "" {
        log.Fatal("MONGO_DB is not set in environment variables")
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET is not set in environment variables")
    }

    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
    smtpEmail := os.Getenv("SMTP_EMAIL")
    smtpPass := os.Getenv("SMTP_PASSWORD")

    return &Config{
        MongoURI:  mongoURI,
        DBName:    dbName,
        JWTSecret: jwtSecret,
        SMTPHost:  smtpHost,
        SMTPPort:  smtpPort,
        SMTPEmail: smtpEmail,
        SMTPPass:  smtpPass,
    }
}