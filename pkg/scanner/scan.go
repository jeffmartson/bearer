package scanner

import (
	"context"
	"fmt"
	"os"

	"github.com/bearer/bearer/new/detector/evaluator/stats"
	classification "github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/detectors"
	"github.com/bearer/bearer/pkg/report/writer"
)

func Scan(
	ctx context.Context,
	rootDir string,
	filename string,
	outputPath string,
	classifier *classification.Classifier,
	fileStats *stats.FileStats,
	scanners []string,
) error {
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_TRUNC, 0666)

	if err != nil {
		return fmt.Errorf("fail opening output file %w", err)
	}
	defer file.Close()

	rep := writer.Detectors{
		Classifier: classifier,
		File:       file,
	}

	if err := detectors.Extract(ctx, rootDir, filename, &rep, fileStats, scanners); err != nil {
		return err
	}

	return nil
}
