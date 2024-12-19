# An Opinionated Starter Kit for Laravel/Livewire with Folio, Volt, and Flux UI

Yes, it’s written in Go.

## Overview

This standalone executable streamlines the setup of Livewire, Folio, Volt, and Flux UI Pro—along with a few extras—into a fresh Laravel project. Simply provide your Flux UI credentials, and the tool takes care of the rest.

This is my go-to setup for new TALL stack projects, and while it’s tailored to my preferences, you’re welcome to try it out. Just remember: **use at your own risk.**

## Customizing

A key advantage of this tool is its use of a `./stubs` directory. Any files placed in this folder will be copied to the project root, maintaining their directory structure. However, **it will MERCILESSLY overwrite existing files** without asking—so handle with care.

## Installation

### Option 1: Using `make build`

1. Clone this repository.
2. Navigate to the project root.
3. Run `make build` to compile the `sourdough` binary, install it to `~/.local/bin`, and make it executable.

### Option 2: Manually with Go

1. Clone this repository.
2. Run `go build -o sourdough`.
3. Move the `sourdough` binary to a directory in your `$PATH` (e.g., `~/.local/bin`).

## Usage

1. Create a fresh Laravel project.
2. Navigate to the project root.
3. Run the `sourdough` command.
4. Follow the prompts.
5. Enjoy your new setup!

## To-Do

- [ ] Automate the creation of a fresh Laravel project before installation.
- [ ] Add a configuration file to simplify adding or removing packages from the install list.
