# Sourdough 🍞
## A Customizable Laravel Installer/Starter Kit

> [!WARNING]
> This was created primarily for personal use - I won't necessarily make changes based on your feedback.
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

## Planned Features

- [x] Support multiple "starters" (read: starter-kits)
- [x] Handle multiple private repos and auth.json credentials
- [ ] Support commandline args for selecting a specific "starter"
- [ ] Allow the injection of code snippets into a specifically targeted file, closure, or array

- [ ] ~~Add a Svelte-Inertia starter~~ (cancelled - for now)

## FAQ

**Why use Go?** Because I wanted Sourdough to compile to a single binary while remaining customizable. Yes, PHP apps can be frankensteined into a self-contained binary, which is cool, but that's not as straightforward for my use case - also I can write Go, so I did.

**Why build web apps with PHP?** Because I like my batteries included for most of the projects I currently work on. Go is awesome, and if someone builds a Laravel equivalent with all the polish it's acquired since its inception, I will consider switching.
