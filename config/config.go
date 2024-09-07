package config

import (
	"errors"
	"fmt"
	"time"
)

// Config represents the application configuration.
type Config struct {
	ServerConfig    ServerConfig    `yaml:"server"`
	DBConfig        DBConfig        `yaml:"db"`
	OpenMeteoConfig OpenMeteoConfig `yaml:"openmeteo"`
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if err := c.DBConfig.Validate(); err != nil {
		return err
	}
	if err := c.OpenMeteoConfig.Validate(); err != nil {
		return err
	}
	if err := c.ServerConfig.Validate(); err != nil {
		return err
	}
	return nil
}

// ServerConfig represents the server configuration.
type ServerConfig struct {
	Port string `yaml:"port"`
}

// Validate validates the server configuration.
func (c ServerConfig) Validate() error {
	if c.Port == "" {
		return errors.New("serverconfig port is required")
	}
	return nil
}

// DBConfig represents the database configuration.
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// Validate validates the database configuration.
func (c DBConfig) Validate() error {
	if c.Host == "" {
		return errors.New("dbconfig host is required")
	}
	if c.Port == "" {
		return errors.New("dbconfig port is required")
	}
	if c.Username == "" {
		return errors.New("dbconfig username is required")
	}
	if c.Password == "" {
		return errors.New("dbconfig password is required")
	}
	if c.DBName == "" {
		return errors.New("dbconfig dbname is required")
	}
	return nil
}

// DSN returns the data source name for the database connection.
func (c DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Username, c.Password, c.Host, c.Port, c.DBName)
}

// OpenMeteoConfig represents the OpenMeteo configuration.
type OpenMeteoConfig struct {
	APIURL  string        `yaml:"api_url"`
	Timeout time.Duration `yaml:"timeout"`
}

// Validate validates the OpenMeteo configuration.
func (c *OpenMeteoConfig) Validate() error {
	if c.APIURL == "" {
		return errors.New("openmeteoconfig api_url is required")
	}
	if c.Timeout == 0 {
		c.Timeout = 15 * time.Second
	}
	return nil
}
