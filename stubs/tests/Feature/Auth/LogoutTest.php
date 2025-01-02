<?php

namespace Tests\Feature\Auth;

use App\Models\User;
use Livewire\Volt\Volt;

it('can render the logout component', function (): void {
    $this->actingAs(User::factory()->create());

    $this->get(route('dashboard', absolute: false))
        ->assertOk()
        ->assertSeeVolt('auth.logout');
});

it('can logout a user', function (): void {
    $this->actingAs(User::factory()->create());

    $component = Volt::test('auth.logout')
        ->call('logout');

    $component->assertRedirect(route('login'));

    $this->assertGuest();
});
