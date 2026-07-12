package client

import (
	"github.com/pion/logging"
	"github.com/rs/zerolog"

	"codeberg.org/kasefuchs/go-kit/log"
)

type pionLoggerFactory struct{}

func (f *pionLoggerFactory) NewLogger(scope string) logging.LeveledLogger {
	logger := log.Logger().With().Str("pion-scope", scope).CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 1).Logger()

	return &pionLogger{logger}
}

type pionLogger struct {
	logger zerolog.Logger
}

func (l *pionLogger) Trace(msg string)                  { l.logger.Trace().Msg(msg) }
func (l *pionLogger) Tracef(format string, args ...any) { l.logger.Trace().Msgf(format, args...) }

func (l *pionLogger) Debug(msg string)                  { l.logger.Debug().Msg(msg) }
func (l *pionLogger) Debugf(format string, args ...any) { l.logger.Debug().Msgf(format, args...) }

func (l *pionLogger) Info(msg string)                  { l.logger.Info().Msg(msg) }
func (l *pionLogger) Infof(format string, args ...any) { l.logger.Info().Msgf(format, args...) }

func (l *pionLogger) Warn(msg string)                  { l.logger.Warn().Msg(msg) }
func (l *pionLogger) Warnf(format string, args ...any) { l.logger.Warn().Msgf(format, args...) }

func (l *pionLogger) Error(msg string)                  { l.logger.Error().Msg(msg) }
func (l *pionLogger) Errorf(format string, args ...any) { l.logger.Error().Msgf(format, args...) }
