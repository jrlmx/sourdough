# Sourdough

A customizable tool for scaffolding new Laravel apps.

## Contents

1) [Intro](#Intro)
2) [Requirements](#Requirements)
3) [Installation](#Installation)
4) [Usage](#Usage)

## Intro

Sourdough is a tool for scaffolding new Laravel projects - not just with the usual auth, but with all the YOUR favorite dependencies, files, and opinions all baked in - configure it once and use it again and again.
If you're like me, then you have a specific set of dependencies, files, or opinions you bring with you whenever you start a new project... and you COULD fork Laravel Breeze and customize it to your liking... but that just feels like overkill for a couple of extra dependencies...?
Or, maybe, also like me you want to create your own personal starter kit(s), and like me you don't want to have a whole Laravel application worth of files cluttering up the repository...
If so, you might want to try Sourdough.
Sourdough uses customizable "Starters": A git repo consisting of a "config.json" file and a "stubs" directory.

The config file let's you define:
- flags to pass to the Laravel installer
- the template: a url pointing to a git repo - sourdough will clone the repo to the new project directory instead of using the Laravel installer
- composer repositories to add before installing dependencies
- default files to remove (e.g. "/resources/views/welcome.blade.php")
- composer/npm dependencies to install (dev or prod) or remove
- artisan/npx commands to run after dependencies have been installed (respectively)
- commands to be run at specific points (hooks) in the build process

The stubs folder contains files to be copied into the newly created laravel project after it's been created. Want to include a specific js file in all your new projects? Now you can... place it in "stubs/resources/js/my-script.js".

>[!Warning]
> This process currently overwrites existing files without prompting... so use with caution

## Requirements

- PHP, Composer, and the Laravel Installer
- NPM (support for other package managers is on the roadmap)
- GIT (sourdough uses the git cli to create a shallow clone of the starter repository)
- Go (only if you're compiling from source)

>[!Tip]
> If you want to install [Go](https://webinstall.dev/go/) you can do it easily with [webi](https://webinstall.dev/)

## Installation

### Compile from Source

Clone the repo
```sh
git clone https://github.com/jrlmx/sourdough.git
```
 Compile the sourdough binary by using "go build ." in the cloned project directory.
```sh
go build .
```
Move the compiled "sourdough" binary to a folder on your $PATH and you're good to go.
Alternatively, if you used webi to install go, it will have added ~/go/bin to your $PATH. You can simply run "go install ." and it will install sourdough to that folder automatically.
```sh
go install .
```

### Download a precompiled binary...

Previous versions of sourdough (like a week ago) used an embeded filesystem to store starters and required a recompilation step to add or modify those starters - this made validating the binary essentially impossible. The current version uses a locally stored sqlite database created when you first fun Sourdough, the database is stored in $HOME/.sourdough (Ubuntu - Linux) or your Operating System's equivalent.

This means that I will likely provide a precompiled binary you can simply download to a folder in your path in the near future and a checksum to validate the executable. Coming... Soon... probably, maybe...

## Usage

By default, there are no starter options in the database, but adding them is easy - just use the "--add" flag.

```sh
sourdough --add
```

You'll be prompted for a name and a url (points to the git repository where the starter - either local or remote). Want to see an example?
Add my example-starter using the add command:

name: jrlmx/example-starter
url: https://github.com/jrlmx/example-starter

Then run sourdough again:

```sh
sourdough myapp
```

Follow the prompts, when the process is complete you should have a new project scaffolded and ready to go.

### Optional Arguments & Flags

Usage: sourdough {project-name} --starter {starter-name}

--starter {starter-name} (select the starter to use)
--add (will display a prompt to add a new Starter)
--remove {starter-name} (will remove starter with name)
--export (exports a json file of stored starter options)
--import (imports a json file to store starter options)
--hooks (displays a list of available command hooks)

Use --help to display all usage information

### Starter Config

An example starter -> config.json
```json
{
    "template": null,
    "flags": ["--no-interaction", "--database=mysql"],
    "remove_files": [
        "resources/views/welcome.blade.php",
        "resources/js/bootstrap.js",
        "tailwind.config.js"
    ],
    "composer": {
        "production": ["livewire/livewire", "laravel/folio", "livewire/volt"],
        "development": ["wire-elements/wire-spy", "soloterm/solo"],
        "remove": [],
        "repositories": []
    },
    "npm": {
        "production":[],
        "development": ["tailwindcss@^4.0.0", "@tailwindcss/vite"],
        "remove": []
    },
    "artisan_commands": ["folio:install", "volt:install", "solo:install"],
    "npx_commands": [],
    "commands": ["quiet:git init", "quiet:git add .", "quiet:git commit -m \"initial\""]
}
```

#### An Explanation of the config.json file:

"template" - A git url to be cloned into the new project directory instead of running "laravel new"

"flags" - Arguments to be passed to the Laravel installer. Note: you should use the "--flag=value" format if possible to avoid unexpected behavior.

"remove_files" - A list of files to remove after the project has been created, but before dependencies are installed. You must provide the relative path from the root of the "stubs directory" (e.g. "resources/views/welcome.blade.php")

"composer" - A manifest of dependencies to be installed or removed. Note: the string is passed as is to the composer require command after being validated, this means you can select a specific version just as you would with the composer cli.

"npm" - A manifest of dependencies to be installed or removed. Works pretty much the same as the Composer manifest, but currently private repositories are not supported for NPM in Sourdough. This is mostly because I've had no need for them.

"artisan_commands" - A list of "php artisan" commands to be run immeadiately after "composer install" (e.g. "folio:install", "volt:install"). If you need more granular control define your command in the commands array and use an appropriate hook - the default amounts to "@composer_install:command".

"npx_commands" - A list of "npx" commands to be run immeadiately after "npm install" (e.g. "npx shadcn@latest init"). If you need more granular control define your command in the commands array and use an appropriate hook - the default amounts to "@npm_install:command".

"commands" - A list of more flexible commands, white-listed so that only git, php, laravel, composer, npm, npx, and node_module/ or vendor/ scripts can be used. To run the command quietly, you can add the "quiet:" prefix and no output will be displayed. Or you can run the command interactively with the "interact:" prefix. Normal command execution displays output, but runs doesn't allow prompts to streamline the build process... *this could potentially cause issues when running commands that need interaction from the user*. Finally, you can use a "@{hook_name}:" prefix to target a specific step in the build process.

> [!Note]
> Hooked commands execute AFTER the selected action has completed - this functionality might be expanded in the future.

#### Adding repositories:

To streamline the installation process, Sourdough collects basic auth credentials for any repositories where "auth" is true, and creates an auth.json file. This isn't strictly required, since the repository_install command runs interactively, it should prompt you for credentials if the repository supports that functionality.
Example Repository:
```json
{
    "name": "flux-pro",
    "url": "composer.fluxui.dev",
    "auth": true
}
```
If the repo you're accessing doesn't suppport basic auth (but still some kind of authentication), you should in theory be able to interactively supply your credentials - just set "auth" to false for that private repository.

## Additional Starters

This project started from a desire to create my own Flux UI starter kit. My [flux-starter](https://github.com/jrlmx/flux-starter) is available on GitHub however it requires that you have a Flux-Pro license to use it - you'll be prompted for your Livewire/Flux credentials as part of the build process. Sourdough and the Flux UI starter are works in progress... so expect bugs and breaking changes - although the former would only affect new projects not existing ones.

A Svelte/Inertia starter is planned - with ShadCn-Svelte installed by default. But, it's not a priority for me at the moment...

## Road to somekind of stablilty...

Automated Testing. Caching of starter files to improve performance. An cross platform installation script for Sourdough to streamline installation.

## Other ideas and plans...

A web UI or TUI for managing installed starter options.
An "optional,feature:" flag to optionally omit commands or dependencies during the build process.
The ability to inject code snippets into a specific file, class, or clossure - useful if overwritting an existing file is not desirable.

## Why write it in GO?

Well, it's not speed, Sourdough utilizes a number of other cli-tools to do what it does, which creates a performance bottleneck - so that complicates the whole "go is faster" argument.

Go applications compile to a single binary, which makes it very portable, and while PHP applications can be Frankensteined into an executable - it's just not as straight-forward for my use case.

Why not write web applications in Go? I like to have my batteries included.
