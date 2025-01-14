<?php

use App\Livewire\Actions\Logout;

$logout = function (Logout $action) {
    $action();

    return $this->redirect('/');
};

?>


<x-layouts.base>
    <div class="min-h-screen content-center">
        <flux:card class="max-w-xl p-6 mx-auto space-y-6 shadow">
            <div class="flex justify-between items-center">
                <flux:heading size="xl">Sourdough</flux:heading>
                <flux:button x-data x-on:click="$flux.dark = ! $flux.dark" icon="moon" variant="subtle" aria-label="Toggle dark mode" />
            </div>

            <p class="prose dark:prose-invert">A highly customizable TALL Stack + Flux UI installer/starter-kit for your next Laravel project.</p>

            <div class="flex justify-between items-center gap-6">
                <a href="https://github.com/jrlmx/sourdough" class="hover:underline">GitHub</a>
                @volt('welcome.nav')
                    <div class="flex items-center gap-2">
                        @auth
                            <flux:button variant="ghost" href="{{ route('dashboard') }}" icon="home" wire:navigate class="hover:underline">Dashboard</flux:button>
                            <flux:button variant="ghost" icon="arrow-right-start-on-rectangle" wire:click="logout">Logout</flux:button>
                        @else
                            <flux:button variant="ghost" href="{{ route('login') }}" icon="arrow-right-start-on-rectangle" wire:navigate>Login</flux:button>
                            <flux:button variant="ghost" href="{{ route('register') }}" icon="arrow-right-start-on-rectangle" wire:navigate>Register</flux:button>
                        @endauth
                    </div>
                @endvolt
            </div>
        </flux:card>
    </div>
</x-layouts.base>
