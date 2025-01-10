<?php

use Illuminate\Support\Facades\Password;

use function Laravel\Folio\middleware;
use function Laravel\Folio\name;
use function Livewire\Volt\rules;
use function Livewire\Volt\state;

name('password.request');

middleware('guest');

state(['email' => '']);

rules(['email' => ['required', 'string', 'email']]);

$sendResetLink = function (): void {
    $this->validate();

    $status = Password::sendResetLink($this->only('email'));

    if ($status != Password::RESET_LINK_SENT) {
        $this->addError('email', 'Invalid email address.');

        return;
    }

    $this->reset('email');

    session()->flash('status', $status);
};
?>


<x-layouts.auth pg_title="Reset Password">
    <p>Please enter your email address to request a password reset link.</p>
    @volt('password.request')
        <div>
            @if (session('status'))
                <small>{{ session('status') }}</small>
            @endif

            <form wire:submit.prevent="sendResetLink" class="space-y-4">
                <flux:input type="email" label="Email" wire:model="email" placeholder="Email" autofocus required />
                <flux:button type="submit" class="outline" wire:loading.attr="disabled" class="w-full">Send Password Reset Link</flux:button>
            </form>
        </div>
    @endvolt
</x-layouts.auth>
