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

    Flux::toast('Profile updated.', variant: 'success');
};

$sendVerification = function () {
    if (auth()->user()->hasVerifiedEmail()) {
        return $this->redirectIntended(route('dashboard', absolute: false));
    }

    auth()->user()->sendEmailVerificationNotification();

    Flux::toast('Email verification sent.');
};

?>


<div>
    <form wire:submit="updateAccount" class="space-y-6">
        <flux:input type="text" label="Name" wire:model="name" placeholder="Name" required autofocus />

        <flux:input type="email" label="Email" wire:model="email" placeholder="Email" required />

        @if (auth()->user() instanceof Illuminate\Contracts\Auth\MustVerifyEmail &&! auth()->user()->hasVerifiedEmail())
            <div class="flex items-center gap-6 justify-between text-red-600">
                <p class="text-sm">Your email address is unverified.</p>
                <flux:button variant="ghost" wire:click="sendVerification">Resend verification link</flux:button>
            </div>
        @endif

        <flux:button type="submit" variant="primary" wire:loading.attr="disabled" class="w-full">Save</flux:button>
    </form>
</div>
