package csv2named

import (
	"iter"
)

type CsvRow []string

type Named []string

type Namespace string

type NamedRows struct {
	Rows iter.Seq2[CsvRow, error]
	Namespace
}
