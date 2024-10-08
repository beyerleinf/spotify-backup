package handler

import (
	"beyerleinf/spotify-backup/ent"
	"beyerleinf/spotify-backup/internal/config"
	logger "beyerleinf/spotify-backup/pkg/log"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	slogger *logger.Logger
	db      *ent.Client
}

func NewHealthHandler(db *ent.Client) *HealthHandler {
	return &HealthHandler{
		slogger: logger.New("health-check", config.AppConfig.Server.LogLevel),
		db:      db,
	}
}

func (h *HealthHandler) GetHealthStatus(c echo.Context) error {
	db_err := h.testDBConnection()

	res := map[string]string{
		"status": "ok",
	}

	if db_err != nil {
		res["database"] = "err"
	} else {
		res["database"] = "ok"
	}

	return c.JSON(http.StatusOK, res)
}

func (h *HealthHandler) testDBConnection() error {
	ctx := context.Background()

	err := h.db.Ping(ctx)
	if err != nil {
		h.slogger.Error("failed to query database", "err", err)
		return fmt.Errorf("database connection failed %v", err)
	}

	return nil
}
