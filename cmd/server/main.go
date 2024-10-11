package main

import (
	"beyerleinf/spotify-backup/ent"
	"beyerleinf/spotify-backup/internal/server/api/handler"
	apiRouter "beyerleinf/spotify-backup/internal/server/api/router"
	"beyerleinf/spotify-backup/internal/server/config"
	uiHandler "beyerleinf/spotify-backup/internal/server/ui/handler"
	uiRouter "beyerleinf/spotify-backup/internal/server/ui/router"
	uiTmpl "beyerleinf/spotify-backup/internal/server/ui/template"
	"beyerleinf/spotify-backup/pkg/logger"
	"beyerleinf/spotify-backup/pkg/router"
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

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load config: %w", err)
		panic(err)
	}

	slogger.SetLogLevel(cfg.Server.LogLevel)

	storageDir := createStorageDir(slogger)

	dbURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.DBName,
		cfg.Database.Password,
	)

	client, err := ent.Open("postgres", dbURL)
	if err != nil {
		slogger.Fatal("Failed opening connection to postgres", "err", err)
		panic(err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		slogger.Fatal("Failed creating schema resources", "err", err)
		panic(err)
	}

	slogger.Info("Connected to database")

	healthHandler := handler.NewHealthHandler(client, cfg)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(logger.GetEchoLogger())
	e.Use(middleware.Recover())

	apiBase := e.Group("/api")
	uiBase := e.Group("/ui")

	router.SetupRoutes(apiBase,
		apiRouter.HealthRoutes(healthHandler),
	)

	renderer, err := uiTmpl.NewRenderer(web.TemplatesFS)
	if err != nil {
		slogger.Fatal("Failed to initialize renderer", "err", err)
	}

	e.Renderer = renderer
	e.StaticFS("/", web.StaticFS)

	spotifyService := spotify.New(client, cfg, storageDir)

	spotifyHandler := uiHandler.NewSpotifyHandler(spotifyService, cfg)

	router.SetupRoutes(uiBase,
		uiRouter.SpotifyRoutes(spotifyHandler),
	)

	slogger.Info(fmt.Sprintf("Starting server on [::]:%d", cfg.Server.Port))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Server.Port)))
}

func createStorageDir(slogger *logger.Logger) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slogger.Error("Error getting home directory", "err", err)
	}

	storageDir := filepath.Join(homeDir, StorageDir)

	err = os.MkdirAll(storageDir, 0755)
	if err != nil {
		slogger.Fatal("Failed to create storage dir", "err", err)
		panic(1)
	}

	slogger.Verbose(fmt.Sprintf("Using storage directory at %s.", storageDir))

	return storageDir
}
