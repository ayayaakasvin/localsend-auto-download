package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/ayayaakasvin/localsend-auto-download/internal/options"
	"github.com/ayayaakasvin/localsend-auto-download/pkg/slogger"
	"github.com/ayayaakasvin/wpn"
)

// main starts the download process and reports any fatal error.
func main() {
	logger := slogger.NewLogger()
	if err := run(context.Background(), logger); err != nil {
		logger.Error("application failed", "error", err)
		os.Exit(1)
	}
}

// run reads the settings, finds the shared files, and downloads them in parallel.
// It returns an error if setup or any download fails.
func run(ctx context.Context, logger *slog.Logger) error {
	dw, err := options.ParseFlags()
	if err != nil {
		if err.Error() == options.HelpCommand {
			return nil
		}
		return fmt.Errorf("parse flags: %w", err)
	}

	logger.Info("downloader configured", "downloader", dw)

	fmd, err := dw.GetFileList(ctx)
	if err != nil {
		return fmt.Errorf("get file list: %w", err)
	}
	if fmd == nil {
		return fmt.Errorf("get file list: empty response")
	}

	logger.Info("received file list", "file_count", len(fmd.Files))

	wp, err := wpn.NewWorkerPool(3)
	if err != nil {
		return fmt.Errorf("initialize worker pool: %w", err)
	}

	if err := wp.Start(ctx); err != nil {
		return fmt.Errorf("start worker pool: %w", err)
	}

	resultsDone := make(chan []error, 1)
	go func() {
		var downloadErrors []error
		for result := range wp.Results() {
			if result.Error != nil {
				downloadErrors = append(downloadErrors, fmt.Errorf("job %s: %w", result.JobID, result.Error))
				logger.Error(
					"download job failed",
					"job_id", result.JobID,
					"worker_id", result.WorkerID,
					"attempts", result.Attempts,
					"duration", result.TimeConsumed,
					"error", result.Error,
				)
				continue
			}

			logger.Info(
				"download job completed",
				"job_id", result.JobID,
				"worker_id", result.WorkerID,
				"attempts", result.Attempts,
				"duration", result.TimeConsumed,
			)
		}
		resultsDone <- downloadErrors
	}()

	var submitErrors []error
	for _, f := range fmd.Files {
		file := f
		if err := wp.Submit(func(ctx context.Context) error {
			if err := dw.DownloadFile(ctx, &file); err != nil {
				return fmt.Errorf("download %q: %w", file.FileName, err)
			}
			return nil
		}, ctx); err != nil {
			submitErrors = append(submitErrors, fmt.Errorf("submit download for %q: %w", file.FileName, err))
		}
	}

	wp.Shutdown(ctx)
	downloadErrors := <-resultsDone

	logger.Info(
		"download run finished",
		"file_count", len(fmd.Files),
		"failed_downloads", len(downloadErrors),
		"failed_submissions", len(submitErrors),
	)

	allErrors := append(submitErrors, downloadErrors...)
	if len(allErrors) > 0 {
		return errors.Join(allErrors...)
	}
	return nil
}
