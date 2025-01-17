# Sourdough 🍞
## A Customizable Laravel Installer/Starter Kit

If you're like me, you have a set of dependencies, configuration files, helpers, or whatever else you just HAVE to have in your new Laravel projects. This tool aims to make it easy to create a new Laravel project with the right dependencies, configuration, and other files in place - and it's completely customizable.

Already have a favorite starter-kit? You can install and extend it with this tool.

Configure it once, use it over and over again.

> [!WARNING]
> This was created primarily for personal use - I won't necessarily make changes based on your feedback unless it's an egregious bug or security issue.
> The installer and the default "starter" are works-in-progress and may contain bugs.

## Installation

1. Clone the repo
2. Compile the executable with Go
3. Move it to a folder in your $PATH

### Prerequisites

- PHP
- Composer
- Laravel Installer (currently)
- Node/NPM
- Go (for compilation)
- A [Flux Pro](https://fluxui.dev/pricing) license (for default starter)

> [!TIP]
> **Don't have Go installed?** It's really easy with [webinstall.dev](https://webinstall.dev/golang/)

## Usage

When run, Sourdough will:

1. Use the Laravel Installer to create a new project
2. Remove any default Laravel files listed in the "remove_files" array and uninstall any PHP/JS dependencies specified in the "remove" arrays of their respective sections (see the example config.json below)
3. Inspect any composer repositories and prompt for credentials if required to create an auth.json file
4. Install and access the repos using the generated auth.json - also ensures auth.json is included in .gitignore
5. Install a set of predefined composer/npm dependencies (see: config.json)
6. Copy the contents of the starter's stubs folder to the new project, mercilessly replacing default files - so be careful with this step

### Creating Custom Starters

Originally designed to scaffold new TALL Stack applications with Livewire/Flux, Sourdough now supports multiple "starters" (starter-kits) - including completely custom ones. To create your own, make a new folder (e.g., "my-starter") in the starters directory, create a config.json to define your composer/npm dependencies including any to be removed from the default Laravel Application structure, and add your boilerplate files to a "stubs" directory as if you were adding them to a new project. **Note**: Sourdough will maintain the "stubs" directory structure when copying files to a fresh Laravel project.

The default "flux" starter will install Livewire, Volt, Folio, Flux and Flux Pro (you'll be prompted for your Flux credentials when you create a new project) and then scaffold an opinionated auth flow using Volt's functional API. Keep in mind it's entirely customizable as long as you're willing to compile your own binary.

### Example starters/my-starter/config.json:

```json
{
    "$schema": "../../config-schema.json",
    "remove_files": ["resources/views/welcome.blade.php", "resources/js/bootstrap.js"],
    "repos": [],
    "php": {
        "prod": ["livewire/livewire", "livewire/volt", "laravel/folio"],
        "dev": ["wire-elements/wire-spy"],
        "remove": []
    },
    "js": {
        "prod": [],
        "dev": ["prettier", "prettier-plugin-blade", "@tailwindcss/typography", "@tailwindcss/forms"],
        "remove": ["axios"]
    },
    "artisan": ["folio:install", "volt:install"],
    "npx": [],
    "commands": ["git init", "git add .", "quiet:git commit -m \"initial\"", "quiet:./vendor/bin/pint"]
}
```

### Example starters folder structure

```
starters/
├── my-starter/
│   ├── config.json
│   └── stubs/
│       ├── app/
│       ├── resources/
│       ├── routes/
│       └── ...
└── flux/
    ├── config.json
    └── stubs/
        ├── app/
        ├── resources/
        ├── routes/
        └── ...
```

## Run a number of user defined cli commands using the "commands" array in your config.json

Currently the commands api supports the following white-list:
php, composer, npm, npx, git - or any script/executable in the /vendors or /node_modules folders... use with caution.

```json
{
    "commands": [
        "interact:git init",
        "git add .",
        "quiet:./vendor/bin/pint"
    ]
}
```

The "interact:" and "quiet:" prefixes change the way Sourdough internally handles the input/output of a given command.

Quiet mode ("quiet:") will not produce any output.
Interactive ("interact:") mode will allow you to interact with any cli as you normally would.
Normal mode (the default - no prefix) simply echos the output without any prompts - this may cause issues with some CLIs that require user interaction.

The -shout flag will override the "quiet:" and force output to be displayed (useful for debugging).

## Hooks

You can use the -hooks flag to view the available hooks for user defined commands.

```bash
sourdough -hooks
```

Hooks execute after the action referenced in the hook name. Usage: `@hook_name:command_name args` in the commands array of the config.json.

## Inspect Starter(s) after compiling the executable

You can use the -config and -tree flags to view the config.json and file tree, respectively, of a specified starter within the embedded filesystem.

```bash
sourdough -config {starter_name}
sourdough -tree {starter_name}
```

## Road to somekind of stability...

Ensure user derived starters are not overwritten when updating the repository. Considering moving the flux starter into it's own repo and leaving the example starter as boilerplate.

Write tests—and then write more tests. Improve command validation. Define a schema for config.json once its structure is finalized. Etc...

## Planned Features

- [x] Support multiple "starters" (read: starter-kits)
- [ ] Support pulling "starters" from a git repo as well or in place of the embedded filesystem (?)
- [x] Handle multiple private repos and auth.json credentials
- [x] Support commandline args for selecting a specific "starter"
- [x] Allow viewing a starter's config and file tree via commandline args
- [x] Support arbitrary commands for things like /vendor script execution during install
- [x] Expose hooks in config.json to allow commands to be executed at specific stages in the the installation process
- [x] Add a "git_url" field to the config.json to optionally pull from a git repo instead of using the laravel installer (in progress)
- [ ] Allow the injection of code snippets into a specifically targeted file, closure, or array
- [ ] Build a TUI for inspecting installed starters, their config, and embeded file systems accessible via commandline flag

### Stretch Goals:
- [ ] Make commands "optional:" by showing a prompt before execution & parse multiple prefix-options like "optional:interact:command"
- [ ] ~~Add a Svelte-Inertia starter~~ (on hold)

## FAQ

**Why use Go?** Because I wanted Sourdough to compile to a single binary while remaining customizable. Yes, PHP apps can be frankensteined into a self-contained binary, which is cool, but that's not as straightforward for my use case

**Why build web apps with PHP?** Because I like my batteries included for most of the projects I currently work on. Go is awesome, and if someone builds a Laravel equivalent with all the polish it's acquired since its inception - I will consider switching
