package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gcathelines/tensor-energy-case/config"
	"github.com/gcathelines/tensor-energy-case/graph"
	"github.com/gcathelines/tensor-energy-case/internal/database"
	"github.com/gcathelines/tensor-energy-case/internal/open_meteo"
	"github.com/gcathelines/tensor-energy-case/internal/usecase"
)

func main() {
	cfg := config.Config{}
	var (
		configPath      string
		graphiQLEnabled bool
	)
	flag.StringVar(&configPath, "config", "", "configuration path")
	flag.BoolVar(&graphiQLEnabled, "graphiql", false, "enable graphiql")
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
	usecase := usecase.NewUsecase(weatherAPI, db)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(usecase)}))

	if graphiQLEnabled {
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	http.Handle("/query", srv)

	log.Printf("running server on port: %s", cfg.ServerConfig.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerConfig.Port, nil))
}
