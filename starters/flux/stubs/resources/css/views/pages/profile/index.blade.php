<?php

use function Laravel\Folio\middleware;
use function Laravel\Folio\name;

name('profile');

middleware(['auth']);

?>


<x-layouts.app pg_title="Profile">
    <div class="space-y-6">
        <div class="grid md:grid-cols-2">
            <flux:heading size="lg" level="2">User Details</flux:heading>
            <livewire:profile.update_account />
        </div>

        <flux:separator variant="subtle" />

        <div class="grid md:grid-cols-2">
            <flux:heading size="lg" level="2">Change Password</flux:heading>
            <livewire:profile.update_password />
        </div>

        <flux:separator variant="subtle" />

        <div class="grid md:grid-cols-2">
            <flux:heading size="lg" level="2">Delete Account</flux:heading>
            <livewire:profile.delete_account />
        </div>
    </div>
</x-layouts.app>
