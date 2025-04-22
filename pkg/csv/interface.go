package csv

type CSVOperator interface {
	Close() error
	Generate(data ...string) error
}
