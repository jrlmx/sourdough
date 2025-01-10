@props(['app_title' => null, 'pg_title' => null])

<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">

<head>
    <!-- Meta -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="csrf-token" content="{{ csrf_token() }}">
    @stack('head.meta')

    <title>
        {{ str($app_title ?? config('app.name'))->when($pg_title, fn($app_title) => $app_title->append(' | ' . $pg_title)) }}
    </title>

    <!-- Fonts -->
    <link rel="preconnect" href="https://fonts.bunny.net">
    <link href="https://fonts.bunny.net/css?family=inter:400,500,600&display=swap" rel="stylesheet" />

    <!-- Styles -->
    @vite('resources/css/app.css')
    @fluxStyles
    @stack('head.styles')

    <!-- Scripts -->
    @vite('resources/js/app.js')
    @stack('head.scripts')
</head>

<body {{ $attributes->merge(['class' => 'min-h-screen bg-white text-black dark:bg-zinc-800 dark:text-white']) }}>
    {{ $slot }}

    @fluxScripts
    @stack('body.scripts')

    @persist('toast')
        <flux:toast />
    @endpersist
</body>

</html>
