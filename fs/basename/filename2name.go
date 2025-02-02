package basename

import (
	"context"
	"path/filepath"
	"strings"

	cn "github.com/takanoriyanagitani/go-csv-add-namespace"
	. "github.com/takanoriyanagitani/go-csv-add-namespace/util"
)

type Basename string

type PathToNamespace func(string) IO[cn.Namespace]

type PathToBasename func(string) IO[Basename]

func Path2base(p string) IO[Basename] {
	return func(_ context.Context) (Basename, error) {
		return Basename(filepath.Base(p)), nil
	}
}

var Path2baseDefault PathToBasename = Path2base

type BasenameToName func(Basename) IO[cn.Namespace]

func (b BasenameToName) ToPathToNamespace(p2b PathToBasename) PathToNamespace {
	return func(p string) IO[cn.Namespace] {
		return Bind(
			p2b(p),
			b,
		)
	}
}

var PathToNamespaceDefault PathToNamespace = RemoveExtDefault.
	ToBasenameToName().
	ToPathToNamespace(Path2baseDefault)

type RemoveExt func(Basename) IO[string]

func (r RemoveExt) ToBasenameToName() BasenameToName {
	return func(b Basename) IO[cn.Namespace] {
		return Bind(
			r(b),
			Lift(func(noext string) (cn.Namespace, error) {
				return cn.Namespace(noext), nil
			}),
		)
	}
}

func Noext(b Basename) IO[string] {
	return func(_ context.Context) (string, error) {
		noext, _, _ := strings.Cut(string(b), ".")
		return noext, nil
	}
}

var RemoveExtDefault RemoveExt = Noext
