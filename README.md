# goenv

Isolated development environments for Go

## Installation

```{bash}
$ go get github.com/vbogretsov/goenv
```

## Usage

```{bash}
$ goenv ENVDIRECTORY IMPORTPATH
$ source ENVDIRECTORY/bin/activate
```

#### Example

```{bash}
$ goenv .env github.com/vbogretsov/goenv
$ source .env/bin/activate
```

#### Deactivate

```{bash}
$ deactivate
```

## Licence

See the LICENCE file.