<?php

namespace Tests\Feature;

use App\Models\User;

it('can render the user dashboard', function () {
    $this->actingAs(User::factory()->create())
        ->get(route('dashboard'))
        ->assertOk();
});
