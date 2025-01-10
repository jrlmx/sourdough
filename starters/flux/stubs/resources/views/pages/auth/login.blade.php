<?php

use App\Livewire\Traits\WithRateLimiting;
use Illuminate\Validation\ValidationException;

use function Laravel\Folio\middleware;
use function Laravel\Folio\name;
use function Livewire\Volt\protect;
use function Livewire\Volt\rules;
use function Livewire\Volt\state;
use function Livewire\Volt\uses;

name('login');
middleware('guest');

uses([WithRateLimiting::class]);

state([
    'email' => '',
    'password' => '',
    'remember' => false,
]);

rules([
    'email' => 'required|string|email',
    'password' => 'required|string',
    'remember' => 'boolean',
]);

$throttleKey = protect(function () {
    return str($this->email)
        ->lower()
        ->append('|'.request()->ip)
        ->transliterate();
});

$login = function () {
    $this->validate();

    $this->throttle(5, function ($seconds) {
        throw ValidationException::withMessages([
            'email' => trans('auth.throttle', [
                'seconds' => $seconds,
                'minutes' => ceil($seconds / 60),
            ]),
        ]);
    });

    if (! auth()->attempt($this->only('email', 'password'), $this->remember)) {
        $this->addError('email', 'The provided credentials do not match our records.');

        return;
    }

    $this->redirectIntended(route('dashboard', absolute: false), navigate: true);
};

?>


<x-layouts.auth pg_title="Sign in">
    @volt('auth.login')
        <div>
            <form wire:submit="login" class="space-y-3">
                <flux:input label="Email" wire:model="email" type="email" placeholder="email@domain.ext" autofocus />
                <flux:input label="Password" wire:model="password" type="password" viewable />
                <flux:checkbox label="Remember Me" wire:model="remember" />
                <flux:button type="submit" variant="primary" class="w-full">Sign in</flux:button>
            </form>
        </div>
    @endvolt

    <flux:subheading class="text-center">
        Don't have an account?
        <flux:link href="{{ route('register') }}">Register</flux:link>
    </flux:subheading>
</x-layouts.auth>
