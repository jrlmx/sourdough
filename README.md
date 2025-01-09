# Sourdough
## A Customizable Laravel Project Starter Kit

> [!WARNING]
> This was created primarily for personal use - I won't necessarily make changes based on your feedback.
> The installer and the default "starter" are works-in-progress and may contain bugs.

Want to try it out? Clone the repo, compile the executable with Go, and move it to a folder in your $PATH.

> [!TIP]
> **Don't have Go installed?** It's really easy with [webinstall.dev](https://webinstall.dev/golang/)

When run, Sourdough will use the Laravel Installer to create a new project, remove/uninstall any default files or dependencies listed in the "remove_files" array (see: config.json), inspect any composer repositories, and prompt for credentials if required to create an auth.json file. This allows the repos to be installed/accessed in the next step. It will then install a set of predefined composer/npm dependencies (see: config.json) and copy the contents of the starter's stubs folder to the new project, mercilessly replacing default files - so be careful with this step.

Originally designed to scaffold new TALL Stack applications with Livewire/Flux, Sourdough now supports multiple "starters" (starter-kits) - including completely custom kits. To create your own, make a new folder (e.g., "my-starter") in the starters directory, create a config.json to define your composer/npm dependencies including any to be removed from the default Laravel Application structure, and add your boilerplate files to a "stubs" directory as if you were adding them to a new project. **Note**: Sourdough will maintain the "stubs" directory structure when copying files to a fresh Laravel project.

The default "flux" starter will install Livewire, Volt, Folio, Flux and Flux Pro (you'll be prompted for your Flux credentials when you create a new project) and then scaffold an opinionated auth flow using Volt's functional API. Keep in mind it's entirely customizable as long as you're willing to compile your own binary.

## Example starters/my-starter/config.json:

```json
{
    "php": {
        "prod": ["laravel/folio", "livewire/livewire", "livewire/volt", "livewire/flux", "livewire/flux-pro"],
        "dev": ["wire-elements/wire-spy"],
        "remove": []
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
    "remove_files": ["resources/views/welcome.blade.php"]
}
```

## Example starters folder structure

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

## Installer Dependencies

PHP
Composer
Laravel Installer (currently)
Node/NPM
Go - to compile the binary

## Default Starter Dependencies

A [Flux Pro](https://fluxui.dev/pricing) license.

**Why use Go?** Because I wanted Sourdough to compile to a single binary while remaining customizable. Yes, PHP apps can be frankensteined into a self-contained binary, which is cool, but that's not as straightforward for my use case - also I can write Go, so I did.

**Why build web apps with PHP?** Because I like my batteries included for most of the projects I currently work on. Go is awesome, and if someone builds a Laravel equivalent with all the polish it's acquired since its inception, I will consider switching - because Go is awesome.
