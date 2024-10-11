package main

import (
	"beyerleinf/spotify-backup/ent"
	"beyerleinf/spotify-backup/internal/api/handler"
	"beyerleinf/spotify-backup/internal/api/router"
	"beyerleinf/spotify-backup/internal/config"
	"beyerleinf/spotify-backup/internal/global"
	uiHandler "beyerleinf/spotify-backup/internal/ui/handler"
	uiRouter "beyerleinf/spotify-backup/internal/ui/router"
	uiTmpl "beyerleinf/spotify-backup/internal/ui/template"
	logger "beyerleinf/spotify-backup/pkg/log"
	"beyerleinf/spotify-backup/pkg/service/spotify"
	"beyerleinf/spotify-backup/web"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

const StorageDir = ".spotify-backup"

func main() {
	slogger := logger.New("main", logger.LevelInfo)

	err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load config: %w", err)
		panic(err)
	}

	slogger.SetLogLevel(config.AppConfig.Server.LogLevel)

	createStorageDir(slogger)

	dbUrl := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Username,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.Password,
	)

	client, err := ent.Open("postgres", dbUrl)
	if err != nil {
		slogger.Fatal("Failed opening connection to postgres", "err", err)
		panic(err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		slogger.Fatal("Failed creating schema resources", "err", err)
		panic(err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		slogger.Error("Error getting home directory", "err", err)
	}

	storageDir := filepath.Join(homeDir, StorageDir)

	slogger.Info("Connected to database")

	healthHandler := handler.NewHealthHandler(client)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(logger.GetEchoLogger())
	e.Use(middleware.Recover())

	apiBase := e.Group("/api")
	uiBase := e.Group("/ui")

	router.SetupRoutes(apiBase,
		router.HealthRoutes(healthHandler),
	)

	renderer, err := uiTmpl.NewRenderer(web.TemplatesFS)
	if err != nil {
		slogger.Fatal("Failed to initialize renderer", "err", err)
	}

	e.Renderer = renderer
	e.StaticFS("/", web.StaticFS)

	spotifyService := spotify.New(client, storageDir)

	spotifyHandler := uiHandler.NewSpotifyHandler(spotifyService)

	router.SetupRoutes(uiBase,
		uiRouter.SpotifyRoutes(spotifyHandler),
	)

	slogger.Info(fmt.Sprintf("Starting server on [::]:%d", config.AppConfig.Server.Port))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.AppConfig.Server.Port)))
}

func createStorageDir(slogger *logger.Logger) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slogger.Fatal("Error getting home directory", "err", err)
		panic(1)
	}

	storageDir := filepath.Join(homeDir, global.StorageDir)

	err = os.MkdirAll(storageDir, 0755)
	if err != nil {
		slogger.Fatal("Failed to create storage dir", "err", err)
		panic(1)
	}

	slogger.Verbose(fmt.Sprintf("Storage directory at %s created/exists.", storageDir))
}
