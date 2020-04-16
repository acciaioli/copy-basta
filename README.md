# copy-basta

[![Build Status](https://travis-ci.com/Spin14/copy-basta.svg?branch=master)](https://travis-ci.com/Spin14/copy-basta)
[![Go Report Card](https://goreportcard.com/badge/github.com/Spin14/copy-basta)](https://goreportcard.com/report/github.com/Spin14/copy-basta)
```
Basta! Stop copying.

This CLI can be used to bootstrap go projects in seconds, and stop the copy paste madness

Usage:
  copy-basta [command]

Available Commands:
  generate    generates new project based on the template and provided parameters
  help        Help about any command
  init        bootstraps a new copy-basta template project

Flags:
  -h, --help               help for copy-basta
      --log-level string   Used to set the logging level.
                           Available options: [debug, info, warn, error, fatal] (default "info")
  -v, --version            version for copy-basta

Use "copy-basta [command] --help" for more information about a command.
```

## Example 

The following command generates a new go project from a template, compiles it and runs it. 

```bash
$ make demo
```
