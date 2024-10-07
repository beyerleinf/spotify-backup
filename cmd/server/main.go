package main

import (
	"beyerleinf/spotify-backup/ent"
	"beyerleinf/spotify-backup/internal/api/handler"
	"beyerleinf/spotify-backup/internal/api/router"
	"beyerleinf/spotify-backup/internal/config"
	"beyerleinf/spotify-backup/internal/ui"
	uiHandler "beyerleinf/spotify-backup/internal/ui/handler"
	uiRouter "beyerleinf/spotify-backup/internal/ui/router"
	uiTmpl "beyerleinf/spotify-backup/internal/ui/template"
	logger "beyerleinf/spotify-backup/pkg/log"
	"beyerleinf/spotify-backup/pkg/service/spotify"
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	slog := logger.New("main", logger.LevelInfo)

	err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load config: %w", err)
		panic(err)
	}

	slog.SetLogLevel(config.AppConfig.Server.LogLevel)

	dbUrl := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Username,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.Password,
	)

	client, err := ent.Open("postgres", dbUrl)
	if err != nil {
		slog.Fatal("Failed opening connection to postgres", "err", err)
		panic(err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		slog.Fatal("Failed creating schema resources", "err", err)
		panic(err)
	}

	slog.Info("Connected to database")

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

	renderer, err := uiTmpl.NewRenderer(ui.PublicFS)
	if err != nil {
		slog.Fatal("Failed to initialize renderer", "err", err)
	}

	e.Renderer = renderer
	e.StaticFS("/", ui.StaticFS)

	spotifyService := spotify.New(client)

	spotifyHandler := uiHandler.NewSpotifyHandler(spotifyService)

	router.SetupRoutes(uiBase,
		uiRouter.SpotifyRoutes(spotifyHandler),
	)

	slog.Info(fmt.Sprintf("Starting server on [::]:%d", config.AppConfig.Server.Port))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.AppConfig.Server.Port)))
}
