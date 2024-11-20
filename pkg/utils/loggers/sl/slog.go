package sl

import (
	"log/slog"
	"os"

	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	return log
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func RegisterDTO(dto *dtos.RegisterDTO) slog.Attr {
	return slog.Group("register_dto",
		slog.String("username", dto.Username),
		slog.String("email", dto.Email),
	)
}

func LoginDTO(dto *dtos.LoginDTO) slog.Attr {
	return slog.Group("login_dto",
		slog.String("username", dto.Username),
	)
}
