package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"go.autokitteh.dev/demodriver/internal/ddsvc"
)

const appName = "dd"

var (
	configDir = filepath.Join(xdg.ConfigHome, appName)

	opts struct {
		verbose  bool
		logLevel string // source of truth for logging.
	}
)

var cmd = cobra.Command{
	Use:   appName,
	Short: "demo driver for temporal workflows",
	Args:  cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if opts.verbose {
			opts.logLevel = "debug"
		}

		if err := initLogger(opts.logLevel); err != nil {
			return fmt.Errorf("logger: %w", err)
		}

		if err := loadDotEnv(); err != nil {
			return fmt.Errorf("dotenv: %w", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ddsvc.New(slog.Default(), appName).Run()
		return nil
	},
}

func init() {
	cmd.PersistentFlags().StringVarP(&opts.logLevel, "log-level", "L", "info", `explicit log level`)
	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, `verbose output`)
	cmd.MarkFlagsMutuallyExclusive("log-level", "verbose")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func loadDotEnv() error {
	err := godotenv.Load(".env", filepath.Join(configDir, ".env"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func initLogger(lvl string) error {
	var lv slog.Level
	if err := lv.UnmarshalText([]byte(lvl)); err != nil {
		return fmt.Errorf("level: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lv}))
	slog.SetDefault(logger)

	return nil
}
