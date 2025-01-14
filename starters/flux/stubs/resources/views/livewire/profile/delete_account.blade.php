<?php

use App\Livewire\Actions\Logout;

use function Livewire\Volt\rules;
use function Livewire\Volt\state;

state([
    'current_password' => '',
    'show_confirmation' => false,
]);

rules([
    'current_password' => 'required|string|current_password',
]);

$deleteAccount = function (Logout $logout) {
    $this->validate();

    tap(Auth::user(), $logout(...))->delete();

    $this->redirect(route('login'));
};

?>


<div x-data class="space-y-6">
    <flux:button variant="danger" type="button" @click="$wire.show_confirmation = true" class="w-full">Delete Account</flux:button>

    <p class="prose dark:prose-invert text-sm">
        Once your account is deleted, all of its resources and data will be permanently deleted. Before deleting your account, please download any data or information that you wish
        to retain.
    </p>

    <flux:modal wire:model="show_confirmation" name="confirm" class="space-y-6">
        <flux:heading size="lg" level="3">Delete Account</flux:heading>

        <p class="prose dark:prose-invert text-sm">
            Are you sure you want to delete your account? Once your account is deleted, all of its resources and data will be permanently deleted. Please enter your password to
            confirm you would like to permanently delete your account.
        </p>

        <form wire:submit="deleteAccount" class="space-y-6">
            <flux:input type="password" label="Password" wire:model="current_password" placeholder="Password" required autofocus />

            <div class="flex gap-2">
                <flux:spacer />
                <flux:button type="button" @click="$wire.show_confirmation = false">Cancel</flux:button>
                <flux:button variant="danger" type="submit" wire:loading.attr="disabled">Delete Account</flux:button>
            </div>
        </form>
    </flux:modal>
</div>
