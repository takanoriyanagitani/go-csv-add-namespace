package stdcsv

import (
	"context"
	"encoding/csv"
	"io"
	"iter"
	"os"

	cn "github.com/takanoriyanagitani/go-csv-add-namespace"
	. "github.com/takanoriyanagitani/go-csv-add-namespace/util"
)

func CsvReaderToRows(r *csv.Reader) iter.Seq2[cn.CsvRow, error] {
	return func(yield func(cn.CsvRow, error) bool) {
		for {
			row, e := r.Read()

			if e == io.EOF {
				return
			}

			var record cn.CsvRow = row
			if !yield(record, e) {
				return
			}
		}
	}
}

func ReaderToRows(r io.Reader) iter.Seq2[cn.CsvRow, error] {
	// csv.Reader uses its own bufio(though no spec found)
	var cr *csv.Reader = csv.NewReader(r)
	return CsvReaderToRows(cr)
}

func ReadCloserToRows(r io.ReadCloser) iter.Seq2[cn.CsvRow, error] {
	var original iter.Seq2[cn.CsvRow, error] = ReaderToRows(r)
	return func(yield func(cn.CsvRow, error) bool) {
		defer r.Close()
		for row, e := range original {
			if !yield(row, e) {
				return
			}
		}
	}
}

func FilenameToRows(filename string) IO[iter.Seq2[cn.CsvRow, error]] {
	return func(_ context.Context) (iter.Seq2[cn.CsvRow, error], error) {
		file, e := os.Open(filename)
		if nil != e {
			return nil, e
		}
		return ReadCloserToRows(file), nil
	}
}
