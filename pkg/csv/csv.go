package csv

import (
	"encoding/csv"
	"os"
	"path/filepath"

	"github.com/hiromaily/aurora-db-data-generator/pkg/logger"
)

type csvGenerator struct {
	logger   logger.Logger
	fileName string
	file     *os.File
	writer   *csv.Writer
}

// csv生成開始前に呼び出し、instanceを生成する
// Note: 並列実行可能なように、1ファイルにつき、1インスタンスを想定
func NewCSVGenerator(logger logger.Logger, fileName string) (*csvGenerator, error) {
	// ディレクトリがない場合は作成する
	if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		return nil, err
	}

	// ファイルオープン 新規作成モード
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// csvGenerator
	return &csvGenerator{
		logger:   logger,
		fileName: fileName,
		file:     file,
		writer:   csv.NewWriter(file),
	}, nil
}

// Note: 必ず呼び出すこと
func (c *csvGenerator) Close() error {
	c.writer.Flush()
	return c.file.Close()
}

// CSVレコード書き込み
func (c *csvGenerator) Generate(data ...string) error {
	return c.writer.Write(data)
}
