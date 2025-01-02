<?php

use App\Livewire\Actions\Logout;

use function Livewire\Volt\state;

state([
    'component' => null,
    'icon' => null,
    'text' => __('Logout'),
])->locked();

$logout = function (Logout $action) {
    $action();

    return redirect()->route('login');
};

?>


@if ($component)
    <x-dynamic-component :$component :$icon wire:click="logout">
        {{ $text }}
    </x-dynamic-component>
@else
    <flux:button :icon="$icon" variant="ghost" wire:click="logout">
        {{ $text }}
    </flux:button>
@endif
