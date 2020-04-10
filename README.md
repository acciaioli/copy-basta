# copy-basta

```bash
$ copy-basta 
Basta! Stop copying.

This CLI can be used to bootstrap go projects in seconds, and stop the copy paste madness

Usage:
  copy-basta [command]

Available Commands:
  generate    generates new project based on the template and provided parameters
  help        Help about any command

Flags:
  -h, --help   help for copy-basta

Use "copy-basta [command] --help" for more information about a command.
```

## Example 

The following command generates a new go project from a template, compiles it and runs it. 

```bash
$ make demo
```

## Todo

- spec file support (templates should be the ones defining what is required)
- default parameters support
- type parameter support (templates should be the ones defining what is required)
- generate from remote repo