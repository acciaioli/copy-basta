# copy-basta

[![Build Status](https://travis-ci.com/acciaioli/copy-basta.svg?branch=master)](https://travis-ci.com/acciaioli/copy-basta)
[![Go Report Card](https://goreportcard.com/badge/github.com/acciaioli/copy-basta)](https://goreportcard.com/report/github.com/acciaioli/copy-basta)


- [Intro](#intro)
- [Install](#install)
- [Template specification](#template-specification)
- [Quick Start](#quick-start)
- [Examples](#examples)

---
## Intro

`copy-basta` is templating command line interface (cli) tool.

It helps managing and using usage template projects with well defined inputs and outputs.

The cli is written in `go` and only [golang's text/template](https://golang.org/pkg/text/template) are supported.
This doesn't mean that you need to know `go` to use this tool, it just means that you need to use golang's templating language
in your template files.

The only opinionated thing you need to in order to use this tool is a `basta.yaml` file alongside your template files.

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

---
## Install

### Binary Releases

For Linux, Mac and Windows binary releases are [available here](https://github.com/acciaioli/copy-basta/releases).

#### Linux
```
sudo /bin/sh -c \
	'wget https://github.com/acciaioli/copy-basta/releases/latest/download/copy-basta.linux-amd64 -O /usr/local/bin/copy-basta \
	&& chmod +x /usr/local/bin/copy-basta'
```



#### Source code (using `go`)

```
▶ git clone git@github.com:acciaioli/copy-basta.git

▶ cd copy-basta

▶ make install
>> installing cli (dev)
>> done

▶ copy-basta --version
copy-basta version snapshot-user-4334710
```

This will run `go install` and the binary will be available from your go path

---
## Template specification
The template specification lives in `basta.yaml` file.

This file contains all the information required by the cli tool to make use of your template project.

Check the example bellow!

```yaml
---
# ignored files will not be copied to generated projects.
# they are for template development only
ignore:
 # ignore this file 
 - ignore-me.txt
 # ignore everything inside this dir 
 - .git/

# pass-through files will be copied untouched to generated projects.
passed-through:
 # passed-through files that start with just-copy
 - just-copy-*
 # passed-through everything inside this dir
 - html/

# in the variables section we declare the inputs we need to generate
# new projects from our template.
variables:
  # simple variable
  - name: recipe
  # with description
  - name: chef
    description: Name of the chef that authored this recipe
  # with type
  - name: estimatedCost
    type: number
  # with default
  - name: isVegan
    type: boolean
    default: false
  # with everything
  - name: igredients
    type: array
    description: Ingredients 
    default: water,salt,love
```
#### More on Variables

##### `variable.name`
The `name` identifies the variable. It should be featured in at least on template file.

The `name` is only required field of a variable.

##### `variable.type`
If provided, it should be an [open API 3.0 type](https://swagger.io/docs/specification/data-models/data-types).

Both default & user provided values are checked against the type.

When the variable type is not specified type checks are skipped.

##### `variable.description`

The description of the variable.

This helps users when they are generating a new project.

##### `variable.default`

The default value for the variable. 

This can be used by users when they are generating a new project.

The provided default must be consistent with the variable type.

___
## Quick Start

Once you have installed the cli tool:

```
▶ copy-basta init --name my-template
[INFO]	validating user input
[INFO]	bootstrapping new template project
        location: my-template
[INFO]	done
```

A new directory called `my-template` was created. 
This is our template project.

```
▶ tree my-template -a
my-template
├── main.sh
├── readme.md
└── basta.yaml

0 directories, 4 files
```

Now you can generate a new project from this template

```
▶ copy-basta generate --src=my-template --dest=new-project
[INFO]	validating user input
[INFO]	loading specification file
[INFO]	parsing template files
[INFO]	getting template variables dynamically

your name so that you can be greeted [string] 
? name    Chi

your favorite greet expression [string] 
? greet [hello]    

[INFO]	creating new project
        location: new-project
[INFO]	done
```

Notice that template specification may provide defaults.
In this case I took the default "greet". 

```
▶ tree new-project -a
new-project
└── main.sh

0 directories, 1 files
```

Our new project is ready!

```
▶ cd new-project
▶ ./main.sh
hello Chi!
```


## Examples

### Gorilla Mux Hello World

```
▶ copy-basta generate \
    --src=https://github.com/acciaioli/gorilla-mux-hello-world-basta-template \
    --dest=my-service
...
▶ cd my-service
▶ go run main.go
Starting up on 8000
```
