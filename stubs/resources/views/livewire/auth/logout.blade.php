<?php

use App\Livewire\Actions\Logout;

use function Livewire\Volt\{state};

state([
    'component' => 'flux:button',
    'icon' => null,
    'text' => __('Logout'),
])->locked();

$logout = function (Logout $action) {
    $action();

    return redirect()->route('login');
};

?>

<x-dynamic-component :$component :$icon wire:click="logout">
    {{ $text }}
</x-dynamic-component>
