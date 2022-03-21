package logger

import "go.uber.org/zap"

var SugarLogger *zap.SugaredLogger

func Init() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync() // flushes buffer, if any
	SugarLogger = logger.Sugar()
	return logger, nil
}
