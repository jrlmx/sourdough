# Sourdough

A Customizable Installer/Starter Kit for Laravel Projects

App Flow:

- Ensure PHP, Composer, Laravel Installer and NPM are installed
- Check if the project folder already exists
- IF the folder already exists... EXIT

- Create a new Laravel Application and descend to the newly created project directory
- Prompt the User to select a Starter Kit (Svelte-Inertia, TALL+Flux UI - determined by top-level folders in "kits" embeded FS)
- Store the corresponding path in the project config for future reference
- Ingest the config.json for the selected kit and store the values in the project config
- Inspect the "repos" key of the project config and determine if any need authorization
    - IF authorization is required: prompt the user for credentials and store them temporarily
- IF any repos required authorization, create an auth.json file and add it to the newly created project folder
- Install Composer dependencies (config.json > "PHP")
- Install NPM dependencies (config.json > "NPM")
- Run Artisan commands
- Copy the files in the "stubs" subfolder to the new project

- IF any steps failed, clean up the project folder so a new run can be attempted

- FIN
