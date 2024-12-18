<?php

namespace App\Livewire\Actions;

class Logout
{
    public function __invoke(): void
    {
        auth()->guard('web')->logout();

        session()->invalidate();
        session()->regenerateToken();
    }
}
