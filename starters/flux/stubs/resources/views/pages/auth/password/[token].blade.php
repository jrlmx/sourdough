<?php

use Illuminate\Auth\Events\PasswordReset;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\Password;
use Illuminate\Validation\Rules;

use function Laravel\Folio\middleware;
use function Laravel\Folio\name;
use function Livewire\Volt\rules;
use function Livewire\Volt\state;

name('password.reset');
middleware('guest');

state(['token'])->locked();

state([
    'email' => fn () => request()->string('email')->value(),
    'password' => '',
    'password_confirmation' => '',
]);

rules([
    'token' => ['required'],
    'email' => ['required', 'string', 'email'],
    'password' => ['required', 'string', 'confirmed', Rules\Password::defaults()],
]);

$resetPassword = function (): void {
    $this->validate();

    $status = Password::reset($this->only('email', 'password', 'password_confirmation', 'token'), function ($user, $password) {
        $user
            ->forceFill([
                'password' => Hash::make($password),
                'remember_token' => Str::random(60),
            ])
            ->save();

        event(new PasswordReset($user));
    });

    if ($status != Password::PASSWORD_RESET) {
        $this->addError('email', $status);

        return;
    }

    Flux::toast($status);

    $this->redirect(route('login'), navigate: true);
};

?>


<x-layouts.auth pg_title="Reset Password">
    @volt('password.reset')
        <div>
            <form wire:submit="resetPassword" class="space-y-3">
                <flux:input type="email" label="email" wire:model="email" autofocus placeholder="email@domain.ext" required />
                <flux:input type="password" label="password" wire:model="password" placeholder="Password" required />
                <flux:input type="password" label="password_confirmation" wire:model="password_confirmation" placeholder="Password" required />
                <flux:button type="submit" variant="primary" class="outline" wire:loading.attr="disabled">Reset Password</flux:button>
            </form>
        </div>
    @endvolt
</x-layouts.auth>
