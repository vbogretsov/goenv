package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const (
	usage = `usage: goenv ENVPATH IMPORT

Create a new isolated development environment for Go.

Version: 1.0.0

Arguments:

    ENVPATH     Path where a new Go environment should be initialized. This
                value will be set to the GOPATH variable.

    IMPORT      Project import path. Example: github.com/author/projectname

Options:

    --help      Print this message and exit.
`
	script = `
# This file must be used with "source activate" or ". activate"

if [[ -n "${GOENV+1}" ]]; then
	deactivate
fi

export GOENV={{.ProjectName}}
export GOENV_OLDPS1=$PS1
export GOENV_OLDGOPATH=$GOPATH
export GOENV_OLDPATH=$PATH

export GOPATH={{.GoPath}}
export PATH="$GOPATH/bin:$PATH"
export PS1="($(basename $GOPATH))$PS1"

mkdir -p $(dirname $GOPATH/src/{{.ImportPath}})
rm -f $GOPATH/src/{{.ImportPath}}
ln -s {{.ProjectPath}} $GOPATH/src/{{.ImportPath}}

deactivate() {
	export PS1=$GOENV_OLDPS1
	export GOPATH=$GOENV_OLDGOPATH
	export PATH=$GOENV_OLDPATH

	unset GOENV GOENV_OLDPS1 GOENV_OLDPATH GOENV_OLDGOPATH
	unset -f deactivate
}
`

	pathExists = `unable to initialize new environemnt, path exists`
)

type params struct {
	ProjectName string
	ProjectPath string
	GoPath      string
	ImportPath  string
}

func isHelpRequired() bool {
	return len(os.Args) == 2 && os.Args[1] == "--help" || len(os.Args) != 3
}

func initEnv(envpath string, importpath string) error {
	if _, err := os.Stat(envpath); err == nil {
		return errors.New(pathExists)
	}

	dirMode := os.ModeDir | 0755

	err := os.MkdirAll(envpath, dirMode)
	if err != nil {
		return err
	}

	binDir := filepath.Join(envpath, "bin")

	err = os.MkdirAll(binDir, dirMode)
	if err != nil {
		return err
	}

	scriptTemplate, err := template.New("activate").Parse(script)
	if err != nil {
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	goPath, err := filepath.Abs(envpath)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	err = scriptTemplate.Execute(&buf, params{
		ProjectName: filepath.Base(pwd),
		ProjectPath: pwd,
		GoPath:      goPath,
		ImportPath:  importpath,
	})

	if err != nil {
		return err
	}

	scriptContent := []byte(buf.String())
	scriptPath := filepath.Join(binDir, "activate")

	err = ioutil.WriteFile(scriptPath, scriptContent, 0744)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if isHelpRequired() {
		fmt.Println(usage)
		os.Exit(1)
	}
	err := initEnv(os.Args[1], os.Args[2])
	if err != nil {
		if _, err = os.Stat(os.Args[1]); err == nil {
			os.RemoveAll(os.Args[1])
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
