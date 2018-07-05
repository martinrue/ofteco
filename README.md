# Frekvenco

Frekvenco builds Esperanto usage stats from a list of YouTube videos.

## Building

```
make build
make build-linux
```

Both targets create binaries in the local `dist` directory.

## Usage

```
→ ./dist/frekvenco --help
Frekvenco (v0.0.1)

Usage:
  frekvenco [config]

Application Config:
  --videos=     path to input file containing video IDs
  --title=      page title in output document
  --header-1=   primary header in output document
  --header-2=   secondary header in output document
  --logo=       URL of logo image in output document
  --logo-link=  URL of logo link in output document
```

## Input File

The input file specified by the `--videos` flag must consist of line-separated YouTube video IDs.

Example:

```
FyMAmlQotvA
moscO9-3KAs
WU1ppQkRLFM
Mf83hCF5Cxg
```

## Output

Frekvenco outputs a standalone HTML rendering of the analysis to `stdout`.
