# An installer for Laravel Livewire w/ Folio, Volt, and Flux UI

Yes... it's written in Go.

## Overview

This standalone executable will install Livewire, Folio, Volt, Flux UI - Pro and other goodies into a fresh Laravel project - you provide your Flux UI credentials and it will do the rest.

This is my personal starting setup for new TALL stack projects - use it at your own risk.

## Customizing...


## Installation

### Using: "make build" command

1. Clone this repo
2. Navigate to the project root
3. Run `make build` - compiles and installs the `sourdough` binary to `~/.local/bin` then makes it executable.

### Manually: Using Go

1. Clone this repo
2. Run `go build -o sourdough`
3. Move the `sourdough` binary to a location in your `$PATH` like `~/.local/bin`

## Usage

1. Create a fresh Laravel project
2. Navigate to the project root
3. Run the `sourdough` command
4. ... follow the prompts
5. Enjoy!

## Todo

- [ ] Create a fresh Laravel project before installing packages
- [ ] Add a file to make adding or removing packages from the install list easy
