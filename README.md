# An Opinionated Installer/Starter Kit for Laravel/Livewire with Folio, Volt, and Flux UI

Yes, it’s written in Go. No I won't compile it for you...

## Overview

This standalone executable simplifies the setup of Livewire, Folio, Volt, and Flux UI Pro—along with a few extras—by creating a fresh Laravel project and configuring it with the necessary dependencies. Just provide your Flux UI credentials, and the tool takes care of the rest.

This is my go-to setup for new TALL stack projects. While it’s tailored to my preferences, you’re welcome to give it a try - or change it to suite your own.

## Prerequisites

Before using this tool, make sure you have the following installed on your system:

- **PHP**
- **Composer**
- **Laravel Installer**

- **Go (Golang)** - needed to compile the executable.

See: The [Laravel Docs](https://laravel.com/docs/11.x/installation) for instalation instructions.

#### Don't have Go installed?

Check out another of my favorite tools - webinstall.dev || [Webi](https://webinstall.dev/) and follow this guide: [Install Go](https://webinstall.dev/golang/)
Or go checkout the Golang Docs...

## Customizing

A key feature of this tool is its use of an embedded `./stubs` directory. Any files placed in this folder (pre-build) will be copied into the new Laravel project, preserving their directory structure. However, **it will MERCILESSLY overwrite existing files** without asking—so handle with care.

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

1. Run the `sourdough` command.
2. Follow the prompts to:
   - Create a fresh Laravel project.
   - Install Livewire, Folio, Volt, Flux UI Pro, and other goodies.
   - Scaffold a typical auth-flow, profile, and layout(s)
3. Enjoy your new setup!

## To-Do

- [ ] Add more prompts & options for package installation.
- [ ] Implement a package list with composer and npm packages to be installed.
- [ ] Include better error handling for unsupported systems.
- [ ] Add fallback to composer create-project when creating the project.
