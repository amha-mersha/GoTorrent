package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DB_DRIVER   string
	DB_HOST     string
	DB_PORT     uint16
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	APP_PORT    string
)

func LoadEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		return err
	}
	DB_DRIVER = os.Getenv("DB_DRIVER")
	DB_HOST = os.Getenv("DB_HOST")
	postgresPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return err
	}
	APP_PORT = os.Getenv("APP_PORT")
	DB_PORT = uint16(postgresPort)
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME = os.Getenv("DB_NAME")
	return nil
}
