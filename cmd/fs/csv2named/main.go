package main

import (
	"context"
	"fmt"
	"iter"
	"log"
	"os"

	cn "github.com/takanoriyanagitani/go-csv-add-namespace"
	fb "github.com/takanoriyanagitani/go-csv-add-namespace/fs/basename"
	cs "github.com/takanoriyanagitani/go-csv-add-namespace/fs/csv/std"
	. "github.com/takanoriyanagitani/go-csv-add-namespace/util"
)

func envValByKey(key string) IO[string] {
	return func(_ context.Context) (string, error) {
		val, found := os.LookupEnv(key)
		switch found {
		case true:
			return val, nil
		default:
			return "", fmt.Errorf("env var %s missing", key)
		}
	}
}

var path2namespace func(string) IO[cn.Namespace] = fb.PathToNamespaceDefault

var filename2rows func(string) IO[iter.Seq2[cn.CsvRow, error]] = cs.
	FilenameToRows

func path2namedRows(filename string) IO[cn.NamedRows] {
	return Bind(
		path2namespace(filename),
		func(ns cn.Namespace) IO[cn.NamedRows] {
			return Bind(
				filename2rows(filename),
				Lift(func(
					rows iter.Seq2[cn.CsvRow, error],
				) (cn.NamedRows, error) {
					return cn.NamedRows{
						Rows:      rows,
						Namespace: ns,
					}, nil
				}),
			)
		},
	)
}

var rows2named func(cn.NamedRows) IO[iter.Seq2[cn.Named, error]] = cs.
	RowsToNamedDefault

func path2namedIter(filename string) IO[iter.Seq2[cn.Named, error]] {
	return Bind(
		path2namedRows(filename),
		rows2named,
	)
}

var csvFilename IO[string] = envValByKey("ENV_INPUT_CSVNAME")

var namedIter IO[iter.Seq2[cn.Named, error]] = Bind(
	csvFilename,
	path2namedIter,
)

var writeNamedRows func(iter.Seq2[cn.Named, error]) IO[Void] = cs.
	WriteNamedRowsToStdout

var filename2iter2stdout IO[Void] = Bind(
	namedIter,
	writeNamedRows,
)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return filename2iter2stdout(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
