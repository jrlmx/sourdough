<?php

namespace App\Livewire\Traits;

use Illuminate\Auth\Events\Lockout;
use Illuminate\Support\Facades\RateLimiter;

trait WithRateLimiting
{
    protected function getThrottleKey(): string
    {
        return request()->ip();
    }

    protected function throttle(int $limit, callable $callback): void
    {
        $throttleKey = $this->getThrottleKey();

        if (! RateLimiter::tooManyAttempts($throttleKey, $limit)) {
            RateLimiter::hit($throttleKey);

            return;
        }

        event(new Lockout(request()));

        $callback(RateLimiter::availableIn($throttleKey));
    }
}
