# copy-basta

[![Build Status](https://travis-ci.com/Spin14/copy-basta.svg?branch=master)](https://travis-ci.com/Spin14/copy-basta)
[![Go Report Card](https://goreportcard.com/badge/github.com/Spin14/copy-basta)](https://goreportcard.com/report/github.com/Spin14/copy-basta)

## Intro

_todo_

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

## Install

### Binary Releases

For Linux, Mac and Windows binary releases are [available here](https://github.com/acciaioli/copy-basta/releases).

### Using go

_todo_

### From source

_todo_

## How to stop copy pasting

To stop copy-pasting we need a `basta` template project.
Once we have one, we can use the `copy-basta` cli tool to create new projects.

### 1) Create a template project

##### What goes into the template project?

1. `static files`.
By default, a file is a `static file`.
When generating a new project, these will just be copied to the new project, not questions asked.
Good example of files that are likely static: `.gitignore`, `a-linter-config.yaml`

1. `spec.yaml`.
This `YAML` file defines the variables that the template needs to generate a new project.

1. `basta files`.
Files with the `.basta` extension are template files.
When generating a new project, these files will be injected with the variables defined in the specification file.
`basta files` must be valid [go text templates](https://golang.org/pkg/text/template).

1. `.bastaignore`.
This text file can be used to let the cli tool which files/directories should be ignored when generating a new project.
Good examples of filers/dir that are likely to go into this file: `.git/`, `.idea/`, the template's `README.md`.


### 2) Generate a new project

##### How do we make use of the template project?

1. [Install](#install) the `copy-basta` cli tool.

1. run the `copy-basta` `generate` command. 
The `--src` (root of the template directory) and `--dest` (new project destination) parameters must be provided.
Ex: `copy-basta generate --src=my-server-template --dest=x-service`.

1. provided the necessary inputs. (You will be asked for them).

That is it all. You should have a new directory with your new project source code.


## Get Started

Once you have installed the cli tool:

```
▶ copy-basta init --name my-template
[INFO]	validating user input
[INFO]	bootstrapping new template project
        location: my-template
[INFO]	done
```

This creates a new directory called `my-template`. This is our template project.

```
▶ tree my-template -a
my-template
├── .bastaignore
├── main.sh.basta
├── readme.md
└── spec.yaml

0 directories, 4 files
```

The files in this directory [should look familiar](#What goes into the template project?)
 
We can now generate a new project from this template.

```
▶ copy-basta generate --src=my-template --dest=new-project
[INFO]	validating user input
[INFO]	loading specification file
[INFO]	parsing template files
[INFO]	getting template variables dynamically

your name so that you can be greeted [string] 
? myName    Chi

your favorite greet expression [string] 
? greet [hello]    

[INFO]	creating new project
        location: new-project
[INFO]	done
```

Notice that template specification may provide defaults. In this case I took the default "greet". 

```
▶ tree new-project -a
new-project
├── main.sh
└── spec.yaml

0 directories, 2 files
```

Our new project is ready!

```
▶ cd new-project
▶ ./main.sh
hello Chi!
```

## Features

_todo_

## RoadMap

_todo_