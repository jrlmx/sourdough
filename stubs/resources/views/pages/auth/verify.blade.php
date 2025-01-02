<?php

use Flux\Flux;

use function Laravel\Folio\middleware;
use function Laravel\Folio\name;

name('verification.notice');
middleware('auth');

$sendVerificationEmail = function (): void {
    if (auth()->user()->hasVerifiedEmail()) {
        $this->redirectIntended(route('dashboard'), absolute: false);

        return;
    }

    auth()->user()->sendEmailVerificationNotification();

    Flux::toast('Verification link sent! Please check your email.');
};

?>

<x-layouts.auth pg_title="Verify Email">
    <p>
        Thanks for signing up! Before getting started, could you verify your email address by clicking on the link we just emailed to you? If you didn't receive the email, we will
        gladly send you another.
    </p>

    @volt('email.verify')
        <div>
            @if (session('status'))
                <small>{{ session('status') }}</small>
            @endif

            <div class="flex items-center justify-between gap-2">
                <flux:button variant="primary" wire:click="sendVerificationEmail">Resend Verification Email</flux:button>

                <flux:button variant="ghost" class="w-full" wire:click="logout">Logout</flux:button>
            </div>
        </div>
    @endvolt
</x-layouts.auth>
