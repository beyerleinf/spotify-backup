package handler

import (
	"beyerleinf/spotify-backup/ent"
	"beyerleinf/spotify-backup/internal/server/config"
	"beyerleinf/spotify-backup/pkg/logger"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// A HealthHandler instance.
type HealthHandler struct {
	slogger *logger.Logger
	db      *ent.Client
	config  *config.Config
}

// NewHealthHandler creates a new instance of the [HealthHandler].
func NewHealthHandler(db *ent.Client, config *config.Config) *HealthHandler {
	return &HealthHandler{
		slogger: logger.New("health-check", config.Server.LogLevel),
		db:      db,
		config:  config,
	}
}

// GetHealthStatus checks the health status of various components and
// returns an API response.
func (h *HealthHandler) GetHealthStatus(c echo.Context) error {
	dbErr := h.testDBConnection()

	res := map[string]string{
		"status": "ok",
	}

	if dbErr != nil {
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
