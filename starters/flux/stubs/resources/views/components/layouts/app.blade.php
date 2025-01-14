<?php

use App\Livewire\Actions\Logout;

$logout = function (Logout $action) {
    $action();

    return redirect()->route('login');
};

?>


@props(['app_title' => null, 'pg_title' => null])

<x-layouts.base :$app_title :$pg_title>
    <flux:header container class="border-b border-zinc-200 bg-zinc-50 dark:border-zinc-700 dark:bg-zinc-900">
        <flux:sidebar.toggle class="lg:hidden" icon="bars-2" inset="left" />

        <flux:brand href="{{ route('dashboard') }}" logo="https://fluxui.dev/img/demo/logo.png" name="Acme Inc." class="dark:hidden max-lg:hidden" />
        <flux:brand href="{{ route('dashboard') }}" logo="https://fluxui.dev/img/demo/dark-mode-logo.png" name="Acme Inc." class="hidden dark:flex max-lg:!hidden" />

        <flux:navbar class="-mb-px max-lg:hidden">
            <flux:navbar.item icon="home" href="{{ route('dashboard', absolute: false) }}">Dashboard</flux:navbar.item>
        </flux:navbar>

        <flux:spacer />

        <flux:button x-data x-on:click="$flux.dark = ! $flux.dark" icon="moon" variant="subtle" aria-label="Toggle dark mode" />

        @volt('user.menu')
            <flux:dropdown position="top" align="start">
                <flux:profile avatar="https://fluxui.dev/img/demo/user.png">
                    <x-slot:avatar>
                        <x-user.avatar :text="auth()->user()->initials" />
                    </x-slot>
                </flux:profile>

                <flux:menu>
                    <flux:menu.item icon="user" href="{{ route('profile') }}">Profile</flux:menu.item>
                    <flux:menu.separator />
                    <flux:menu.item icon="arrow-right-start-on-rectangle" wire:click="logout">Logout</flux:menu.item>
                </flux:menu>
            </flux:dropdown>
        @endvolt
    </flux:header>

    <flux:sidebar stashable sticky class="border-r border-zinc-200 bg-zinc-50 dark:border-zinc-700 dark:bg-zinc-900 lg:hidden">
        <flux:sidebar.toggle class="lg:hidden" icon="x-mark" />

        <flux:brand href="{{ route('dashboard') }}" logo="https://fluxui.dev/img/demo/logo.png" name="Acme Inc." class="px-2 dark:hidden" />
        <flux:brand href="{{ route('dashboard') }}" logo="https://fluxui.dev/img/demo/dark-mode-logo.png" name="Acme Inc." class="hidden px-2 dark:flex" />

        <flux:navlist variant="outline">
            <flux:navlist.item icon="home" href="{{ route('dashboard', absolute: false) }}">Dashboard</flux:navlist.item>
        </flux:navlist>
    </flux:sidebar>

    <flux:main container>
        {{ $slot }}
    </flux:main>
</x-layouts.base>
