<?php

namespace App\Livewire\Traits;

use Illuminate\Auth\Events\Lockout;
use Illuminate\Support\Facades\RateLimiter;

trait WithRateLimiting
{
    protected function throttle(string $key, int $limit, callable $callback): void
    {
        if (! RateLimiter::tooManyAttempts($key, $limit)) {
            RateLimiter::hit($key);

            return;
        }

        event(new Lockout(request()));

        $callback(RateLimiter::availableIn($key));
    }
}
