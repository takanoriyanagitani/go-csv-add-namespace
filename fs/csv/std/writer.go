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

type HeaderConfig struct {
	HasHeader       bool
	NamespaceHeader string
}

var HeaderConfigEmpty HeaderConfig

var HeaderConfigDefault HeaderConfig = HeaderConfig{
	HasHeader:       true,
	NamespaceHeader: "namespace",
}

type IsHeader bool

type RowToNamed func(cn.CsvRow, HeaderConfig, IsHeader) IO[cn.Named]

type RowsToNamed func(cn.NamedRows) IO[iter.Seq2[cn.Named, error]]

func (c RowToNamed) ToRowsToNamed() RowsToNamed {
	var cfg HeaderConfig = HeaderConfigDefault

	return func(named cn.NamedRows) IO[iter.Seq2[cn.Named, error]] {
		return func(ctx context.Context) (iter.Seq2[cn.Named, error], error) {
			return func(yield func(cn.Named, error) bool) {
				var empty cn.Named
				var isHeader IsHeader = true

				for original, e := range named.Rows {
					if nil != e {
						yield(empty, e)
						return
					}

					namedRow, e := c(original, cfg, isHeader)(ctx)
					isHeader = false

					if !yield(namedRow, e) {
						return
					}
				}
			}, nil
		}
	}
}

func RowToNamedFromNamespace(ns cn.Namespace) RowToNamed {
	var buf []string
	return func(original cn.CsvRow, cfg HeaderConfig, h IsHeader) IO[cn.Named] {
		return func(_ context.Context) (cn.Named, error) {
			buf = buf[:0]

			if bool(h) && cfg.HasHeader {
				buf = append(buf, cfg.NamespaceHeader)
			} else {
				buf = append(buf, string(ns))
			}

			for _, col := range original {
				buf = append(buf, col)
			}

			return buf, nil
		}
	}
}

func Rows2named(rows cn.NamedRows) IO[iter.Seq2[cn.Named, error]] {
	var row2named RowToNamed = RowToNamedFromNamespace(rows.Namespace)
	var rows2named RowsToNamed = row2named.ToRowsToNamed()
	return rows2named(rows)
}

var RowsToNamedDefault RowsToNamed = Rows2named

type WriteNamedRows func(iter.Seq2[cn.Named, error]) IO[Void]

func NamedPairsToCsvWriter(
	ctx context.Context,
	pairs iter.Seq2[cn.Named, error],
	wtr *csv.Writer,
) error {
	for cols, e := range pairs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if nil != e {
			return e
		}

		e := wtr.Write(cols)
		if nil != e {
			return e
		}
	}

	wtr.Flush()

	return wtr.Error()
}

func NamedPairsToWriter(
	ctx context.Context,
	pairs iter.Seq2[cn.Named, error],
	wtr io.Writer,
) error {
	return NamedPairsToCsvWriter(
		ctx,
		pairs,
		csv.NewWriter(wtr), // csv.Writer uses its own bufio(no spec found)
	)
}

func NamedPairs2stdout(pairs iter.Seq2[cn.Named, error]) IO[Void] {
	return func(ctx context.Context) (Void, error) {
		return Empty, NamedPairsToWriter(
			ctx,
			pairs,
			os.Stdout,
		)
	}
}

var WriteNamedRowsToStdout WriteNamedRows = NamedPairs2stdout
