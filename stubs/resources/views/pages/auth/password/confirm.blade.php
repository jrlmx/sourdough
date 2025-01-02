<?php

use Illuminate\Validation\ValidationException;

use function Laravel\Folio\middleware;
use function Laravel\Folio\name;
use function Livewire\Volt\rules;
use function Livewire\Volt\state;

name('password.confirm');

middleware(['auth']);

state([
    'password' => '',
]);

rules([
    'password' => 'required|string|current_password',
]);

$confirm = function (): void {
    $this->validate();

    if (
        ! auth()
            ->guard('web')
            ->validate([
                'email' => auth()->user()->email,
                'password' => $this->password,
            ])
    ) {
        throw ValidationException::withMessages([
            'password' => __('auth.password'),
        ]);
    }

    session(['auth.password_confirmed_at' => time()]);

    $this->redirectIntended(default: route('dashboard', absolute: false), navigate: true);
};

?>


<x-layouts.auth pg_title="Confirm Password">
    <p>This is a secure area of the application. Please confirm your password before continuing.</p>
    @volt('password.confirm')
        <div>
            <form wire:submit="confirm" class="space-y-3">
                <flux:input type="password" label="password" wire:model="password" placeholder="Password" required />
                <flux:button type="submit" variant="primary" class="outline" wire:loading.attr="disabled">Confirm Password</flux:button>
            </form>
        </div>
    @endvolt
</x-layouts.auth>
