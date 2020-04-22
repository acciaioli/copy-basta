# copy-basta

[![Build Status](https://travis-ci.com/Spin14/copy-basta.svg?branch=master)](https://travis-ci.com/Spin14/copy-basta)
[![Go Report Card](https://goreportcard.com/badge/github.com/Spin14/copy-basta)](https://goreportcard.com/report/github.com/Spin14/copy-basta)


- [Intro](#intro)
- [Install](#install)
- [How to stop copy pasting](#how-to-stop-copy-pasting)
- [Quick Start](#quick-start)
- [RoadMap](#roadmap)


---
## Intro

`copy-basta` is templating command line interface (cli) tool.
It aims to support the usage of template code bases with well defined inputs and outputs.

The cli is written in `go` and only [golang's text/template](https://golang.org/pkg/text/template) are supported.
This doesn't mean that you need to know `go` to use this tool, it just means that you need to use golang's templating language
in your template files.

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

For Linux, Mac and Windows binary releases are [vailable here](https://github.com/acciaioli/copy-basta/releases).

#### Using go (compile source code)

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


## Quick Start

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

The files in this directory [should look familiar](#what-goes-into-the-template-project)
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
└── main.sh

0 directories, 1 files
```

Our new project is ready!

```
▶ cd new-project
▶ ./main.sh
hello Chi!
```

---
## RoadMap

- proper support for `array` and `object` types (including input prompts)
- generate from remote location (--src=https://github.com/acciaioli/example-template)
- deal with variables with dashes (this needs to be forbidden or worked around, no good for go templates)
- code documentation
- cli documentation
- heathcheck command (quickly tests the template, for a dev friendly experience)
- support install via snap, apt, brew, etc...


---
## Template Project Documentation

Every file in your template is either a `static file` or a `template file`.

### Template files

A file that ends with `.basta` is a template file. For template files, 
both file name and file content must be valid go text/template templates.

When generating a new project, 
these files are copied over to the new dir, 
the `.basta` extension is dropped 
and the variable placeholders get replaced with the variables provided by the user.
File permissions are kept.

(the above is true unless a file [is ignored](#ignoring-files)

### Static files

A file is static by default. All files that are not `basta templates` are `static files`. 

When generating a new project static files are copied over to the new dir.
File permissions are kept.

(the above is true unless a file [is ignored](#ignoring-files)

There are to "special" static files, the [specification file](#specification) and the [ignore file](#ignoring-files).

### Specification

The specification file (`spec.yaml`) is where we define all the 
variables required to generate a new project from our template.

These variables will have to be provided when generating a new project.

Example `spec.yaml`

```yaml
---
variables:
    - name: color
      type: string
      description: your favorite color
      default: blue
    - name: luckyNumber
      type: integer
      description: your lucky number
```
---
#### `variable.name`
__required__, __string__

The name of the variable. These names are the variables of our templates.

---
#### `variable.type`
__optional__, __string__

The type of the variable. 
If provided, it should be an [open API 3.0 type](https://swagger.io/docs/specification/data-models/data-types)
Both default & user provided values are checked against the type.

When not specified, the type of the variable becomes `any` and type checks are skipped.

---
#### `variable.description`
__optional__, __string__

The description of the variable. This helps the user when they are generating a new project.

---
#### `variable.default`
__optional__, __any__

The default value for the variable. Can be taken by the user when generating a new project.

The provided default must be consistent with the variable type.


### Ignoring Files

The ignore file (`.bastaignore`) is where we define all the 
directories, files and file patterns that we wish to ignore 
when generating a new project.

Ignore files will not be processed, therefore the generated project 
will not contain any of them.

Besides project specific files, 
we probably want to ignore the 
`.git` directory, 
the `spec.yaml` file 
and the `.bastaignore` file.

__Note that these files are not alike `gitignore` files.__

Example `.bastaignore`

```
# this is a comment

# this is how we can ignore everything under a directory
ignored_dir/

# this is how we can ignore a single file 
ignored_file.txt

# this is how we can ignore file via pattern
ignored_pattern*
*.
```
