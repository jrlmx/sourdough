<?php

namespace Tests\Feature\Auth;

use App\Models\User;
use Livewire\Volt\Volt;

it('can render the user menu', function (): void {
    $this->actingAs(User::factory()->create());

    $this->get(route('dashboard', absolute: false))
        ->assertOk()
        ->assertSeeVolt('user.menu');
});

it('can logout a user', function (): void {
    $this->actingAs(User::factory()->create());

    $component = Volt::test('user.menu')
        ->call('logout');

    $component->assertRedirect(route('login'));

    $this->assertGuest();
});
