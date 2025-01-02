<?php

use Flux\Flux;
use Illuminate\Validation\Rule;

use function Livewire\Volt\rules;
use function Livewire\Volt\state;

state([
    'name' => fn () => auth()->user()->name,
    'email' => fn () => auth()->user()->email,
]);

rules([
    'name' => ['required', 'string', 'min:3', 'max:255'],
    'email' => ['required', 'string', 'email', 'lowercase', 'max:255', Rule::unique('users')->ignore(auth()->user()->id)],
]);

$updateAccount = function (): void {
    $this->validate();

    $user = auth()->user();

    $user->fill($this->only('name', 'email'));

    if ($user->isDirty('email')) {
        $user->email_verified_at = null;
    }

    $user->save();

    Flux::toast('Profile updated.');
};

?>


<div>
    <form wire:submit="updateAccount" class="space-y-6">
        <flux:input type="text" label="Name" wire:model="name" placeholder="Name" required autofocus />

        <flux:input type="email" label="Email" wire:model="email" placeholder="Email" required />

        <flux:button type="submit" variant="primary" wire:loading.attr="disabled" class="w-full">Save</flux:button>
    </form>
</div>
