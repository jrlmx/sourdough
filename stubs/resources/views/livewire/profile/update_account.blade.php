<?php

use Flux\Flux;
use Illuminate\Validation\Rule;

use function Livewire\Volt\{state, mount, rules};

state([
    'name' => '',
    'email' => '',
]);

rules([
    'name' => ['required', 'string', 'min:3', 'max:255'],
    'email' => ['required', 'string', 'email', 'lowercase', 'max:255', Rule::unique('users')->ignore(auth()->user()->id)],
]);

mount(function (): void {
    $this->name = auth()->user()->name;
    $this->email = auth()->user()->email;
});

$update = function (): void {
    $this->validate();

    $user = auth()->user();

    if ($this->email !== auth()->user()->email) {
        $user->email_verified_at = null;
    }

    auth()->user()->update($this->only('name', 'email'));

    Flux::toast('Profile updated.');
};

?>

<div>
    <form wire:submit="update" class="space-y-6">
        <flux:input type="text" label="Name" wire:model="name" placeholder="Name" required autofocus />

        <flux:input type="email" label="Email" wire:model="email" placeholder="Email" required />

        <flux:button type="submit" variant="primary" wire:loading.attr="disabled" class="w-full">
            Save
        </flux:button>
    </form>
</div>
