package main

import (
	"log/slog"
	"os"

	"github.com/alexflint/go-arg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/hiromaily/aurora-db-data-generator/pkg/logger"
)

type args struct {
	Count   int    `arg:"--count,required" help:"Number of data items to be output: Maximum is 999999"`
	AppName string `arg:"--app"            help:"Application Code: If not specified, all test data will be generated e.g. app1"`
}

func main() {
	// logger
	logger := logger.NewConsoleLogger(slog.LevelDebug)

	// parse command line arguments
	var args args
	arg.MustParse(&args)
	if args.Count > 999999 {
		logger.Error("args count must be less 100000", "count", args.Count)
		return
	}
	logger.Info("args parsed", "data-count", args.Count, "app-name", args.AppName)

	// environment variables
	err := godotenv.Load()
	if err != nil {
		logger.Error("failed to parse env", "error", err)
	}
	mysqlDSN := os.Getenv("MYSQL_DSN")

	// Application
	app, err := newApp(logger, mysqlDSN)
	if err != nil {
		logger.Error("failed to call NewApp()", "error", err)
		return
	}
	defer app.close() //nolint: errcheck

	switch args.AppName {
	case "app1":
		if err := app.generateApp01(args.Count); err != nil {
			logger.Error("failed to call generateApp01()", "error", err)
			return
		}
	case "app2":
		// TODO: Implement logic for b002
	}
}
