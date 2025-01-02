# Sourdough 🍞

An opinionated installer and starter kit for Laravel/Livewire projects with Folio, Volt, and Flux UI Pro integration.

**Note**: I do not plan to release precompiled binaries for this tool. If you come across a compiled version, please ensure you trust the source before downloading or using it. You are solely responsible for your security and any consequences resulting from the use or misuse of this tool.

## Overview

Sourdough is a Go-based CLI tool that automates the setup of a fresh Laravel project with the complete TALL stack (Tailwind CSS, Alpine.js, Laravel, Livewire) along with Folio, Volt, and Flux UI Pro. It includes an opinionated starter kit that can be customized to your needs before compilation.

The tool streamlines the process by:
- Creating a fresh Laravel project
- Removing unwanted default Laravel scaffolding
- Installing and configuring the TALL stack components
- Setting up Laravel/Folio and Livewire/Volt
- Integrating Flux UI Pro (requires valid credentials)
- Installing customizable starter files

## Prerequisites

- PHP
- Composer
- Laravel Installer (Composer create-project support coming soon)
- Node.js and npm (yarn/pnpm support planned)
- Go (for compilation)
- GNU Make (optional, for `make install`)
- Valid Flux UI Pro license

For Laravel installation instructions, refer to the [official Laravel documentation](https://laravel.com/docs/installation).

### Installing Go

If you don't have Go installed, you can use [webinstall.dev](https://webinstall.dev/):
```bash
# Using webinstall.dev
curl -sS https://webinstall.dev/golang | bash
```
Or follow the instructions in the [official Go documentation](https://golang.org/doc/install).

## Installation

### Option 1: Using Make (Fastest - if you're using WSL or Ubuntu)
```bash
# 1. Clone the repository
git clone https://github.com/yourusername/sourdough.git

# 2. Navigate to the project
cd sourdough

# 3. Build and install (this will compile the binary and move it to ~/.local/bin)
make install
```

### Option 2: Manual Go Installation (Recommended)
```bash
# 1. Clone the repository
git clone https://github.com/yourusername/sourdough.git

# 2. Navigate to the project
cd sourdough

# 3. Build the binary
go build -o sourdough

# 4. Move to a directory in your PATH (e.g., ~/.local/bin)
# I'm using Ubuntu on WSL, so I moved it to ~/.local/bin
mv sourdough ~/.local/bin/
chmod +x ~/.local/bin/sourdough
```

## Usage

1. Run the installer:
```bash
sourdough
```

2. Follow the interactive prompts to:
   - Specify your project details
   - Enter your Flux UI Pro credentials
   - Configure any custom options

3. The tool will automatically:
   - Create a fresh Laravel project
   - Install and configure all dependencies
   - Set up authentication flows
   - Apply your custom starter kit
   - Remove unwanted packages/files

## Customization

### Starter Kit Templates

The tool uses an embedded `./stubs` directory for custom templates. Any files placed in this directory before compilation will be copied into new projects, maintaining their directory structure.

⚠️ **Warning**: Files from the starter kit will overwrite existing files without prompting.

### Configuration

Project dependencies and file cleanup can be customized through the `config.json` file before compilation. Edit this file to:
- Add/remove packages to be installed
- Specify files to be removed from the default Laravel installation
- Configure other build preferences

Example `config.json`:
```json
{
   // Specify PHP and JS packages to be installed
   "packages": {
       "php": [
           "laravel/folio",
           "livewire/livewire",
           "livewire/volt",
           "livewire/flux",
           "livewire/flux-pro"
       ],
       "js": [
           "prettier",
           "prettier-plugin-blade",
           "@tailwindcss/typography"
       ]
   },
   "cleanup": {
       // Specify files to be removed from the default Laravel installation
       "files": [
           "resources/views/welcome.blade.php"
       ],
       // Specify PHP and JS packages to be removed
       "packages": {
           "php": [],
           "js": []
       }
   }
}
```

## Roadmap

- [ ] Offer to install PHP, Composer, and Laravel installer via Laravel's `php.new` command
- [ ] Support for `composer create-project` as an alternative to the Laravel installer
- [ ] Package manager alternatives (yarn, pnpm)
- [ ] Additional starter kit templates
