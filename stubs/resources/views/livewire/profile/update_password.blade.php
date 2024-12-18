<?php

use Flux\Flux;
use Illuminate\Validation\Rules\Password;

use function Livewire\Volt\{state, rules};

state([
    'password' => '',
    'new_password' => '',
    'new_password_confirmation' => '',
]);

rules([
    'password' => 'required|string|current_password',
    'new_password' => ['required', 'confirmed', 'string', Password::defaults()],
]);

$update = function (): void {
    $this->validate();

    auth()
        ->user()
        ->update([
            'password' => $this->password,
        ]);

    $this->reset(['password', 'new_password', 'password_confirmation']);

    Flux::toast('Password updated.');
};

?>

<div>
    <form wire:submit="update" class="space-y-6">
        <flux:input label="Current Password" type="password" wire:model="password" placeholder="Current Password" viewable
            required autofocus />

        <flux:input label="New Password" type="password" wire:model="new_password" placeholder="New Password" viewable
            required />

        <flux:input label="Confirm New Password" type="password" wire:model="new_password_confirmation"
            placeholder="Confirm New Password" viewable required />

        <flux:button type="submit" variant="primary" wire:loading.attr="disabled" class="w-full">
            Save
        </flux:button>
    </form>
</div>
