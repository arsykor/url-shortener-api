package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"url-shortener-api/internal/composites"
	"url-shortener-api/internal/config"
	"url-shortener-api/pkg/clients/postgresql"
)

func main() {
	storage := flag.String("storage", "im", "Postgres as a data source pass -db; in-memory -im")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)

	var linkComposite *composites.LinkComposite

	switch *storage {
	case "im":
		linkComposite = composites.NewLinkCompositeInMemory(addr)
	case "db":
		postgresqlClient, err := postgresql.NewClient(context.TODO(), 5, conf)
		if err != nil {
			log.Fatal(err)
		}
		linkComposite = composites.NewLinkCompositePostgres(postgresqlClient, addr)
	}

	router := gin.Default()
	linkComposite.Handler.Register(router)
	err = router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
