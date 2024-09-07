package main

import (
	"database/sql"
	"flag"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"

	"github.com/gcathelines/tensor-energy-case/config"
	"github.com/gcathelines/tensor-energy-case/internal/database"
	"github.com/gcathelines/tensor-energy-case/internal/open_meteo"
	"github.com/gcathelines/tensor-energy-case/internal/usecase"
)

func main() {
	cfg := config.Config{}
	configPath := ""
	flag.StringVar(&configPath, "config", "", "configuration path")
	flag.Parse()

	out, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(out, &cfg); err != nil {
		panic(err)
	}

	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	// initialize database
	sqlDB, err := sql.Open("postgres", cfg.DBConfig.DSN())
	if err != nil {
		panic(err)
	}
	db := database.NewDatabase(sqlDB)

	// initialize open_meteo client
	weatherAPI := open_meteo.NewOpenMeteoClient(cfg.OpenMeteoConfig)
	_ = usecase.NewUsecase(weatherAPI, db)
}
