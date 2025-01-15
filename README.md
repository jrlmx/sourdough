# Sourdough 🍞
## A Customizable Laravel Installer/Starter Kit

> [!WARNING]
> This was created primarily for personal use - I won't necessarily make changes based on your feedback unless it's a egregious bug or security issue.
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

Originally designed to scaffold new TALL Stack applications with Livewire/Flux, Sourdough now supports multiple "starters" (starter-kits) - including completely custom kits. To create your own, make a new folder (e.g., "my-starter") in the starters directory, create a config.json to define your composer/npm dependencies including any to be removed from the default Laravel Application structure, and add your boilerplate files to a "stubs" directory as if you were adding them to a new project. **Note**: Sourdough will maintain the "stubs" directory structure when copying files to a fresh Laravel project.

The default "flux" starter will install Livewire, Volt, Folio, Flux and Flux Pro (you'll be prompted for your Flux credentials when you create a new project) and then scaffold an opinionated auth flow using Volt's functional API. Keep in mind it's entirely customizable as long as you're willing to compile your own binary.

### Example starters/my-starter/config.json:

```json
{
    "php": {
        "prod": ["laravel/folio", "livewire/livewire", "livewire/volt", "livewire/flux", "livewire/flux-pro"],
        "dev": ["wire-elements/wire-spy"],
        "remove":[]
    },
    "js": {
        "prod": [],
        "dev": ["prettier", "prettier-plugin-blade", "@tailwindcss/typography"],
        "remove": []
    },
    "repos": [
        {
            "name": "flux-pro",
            "url": "composer.fluxui.dev",
            "auth": true
        }
    ],
    "artisan": ["folio:install", "volt:install"],
    "npx": [],
    "remove_files": ["resources/views/welcome.blade.php"]
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

Note: I plan to introduce lifecycle hooks to allow execution at specific points in the installation process - so the config.json "api" for this feature is likely to change or be expanded on in some form - see the "road to stability" section below.

## Inspect Starter(s) after compiling the executable

You can use the -config and -tree flags to view the config.json and file tree, respectively, of a specified starter within the embedded filesystem.

```bash
sourdough -config {starter_name}
sourdough -tree {starter_name}
```

## Road to somekind of stability...

Move the default starter(s) to a separate unembedded directory or a different repository...
...Develop a compilation or installation script to check the contents of the starters folder. If no starters are detected, prompt the user to move or download the default starters into the empty folder. This approach minimizes the risk of overwriting a user’s existing starter(s) when updating their local repository to the latest version.
...Alternatively, rename the "flux" starter to "default" and update the .gitignore to exclude all other subdirectories. This method is less robust but simpler to implement.

Write tests—and then write more tests. Improve command validation. Define a schema for config.json once its structure is finalized. Etc...

## Planned Features

- [x] Support multiple "starters" (read: starter-kits)
- [x] Handle multiple private repos and auth.json credentials
- [x] Support commandline args for selecting a specific "starter"
- [x] Allow viewing a starter's config and file tree via commandline args
- [x] Support arbitrary commands for things like /vendor script execution during install (in-progress)
- [ ] Expose hooks in config.json to allow commands to be executed at specific stages in the the installation process
- [ ] Build a TUI for inspecting installed starters, their config, and embeded file systems accessible via commandline flag
- [ ] Allow the injection of code snippets into a specifically targeted file, closure, or array

- [ ] ~~Add a Svelte-Inertia starter~~ (cancelled - for now)

## FAQ

**Why use Go?** Because I wanted Sourdough to compile to a single binary while remaining customizable. Yes, PHP apps can be frankensteined into a self-contained binary, which is cool, but that's not as straightforward for my use case - also I can write Go, so I did.

**Why build web apps with PHP?** Because I like my batteries included for most of the projects I currently work on. Go is awesome, and if someone builds a Laravel equivalent with all the polish it's acquired since its inception, I will consider switching.
