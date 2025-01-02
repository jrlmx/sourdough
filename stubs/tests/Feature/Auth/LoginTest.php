<?php

namespace Tests\Feature\Auth;

use App\Models\User;
use Livewire\Volt\Volt;

it('can render the login page', function (): void {
    $response = $this->get(route('login'))
        ->assertOk()
        ->assertSeeVolt('auth.login');
});

it('can login a user', function (): void {
    User::create([
        'name' => 'Test User',
        'email' => 'test@example.com',
        'password' => 'password',
    ]);

    $component = Volt::test('auth.login')
        ->set('email', 'test@example.com')
        ->set('password', 'password');

    $component->call('login');

    $component->assertHasNoErrors(['email', 'password']);

    $component->assertRedirect(route('dashboard', absolute: false));

    $this->assertAuthenticated();
});

it('cannot login a user with invalid credentials', function (): void {
    User::create([
        'name' => 'Test User',
        'email' => 'test@example.com',
        'password' => 'password',
    ]);

    $component = Volt::test('auth.login')
        ->set('email', 'test@example.com')
        ->set('password', 'wrong-password');

    $component->call('login');

    $component->assertHasErrors(['email'])
        ->assertNoRedirect();

    $this->assertGuest();
});

it('can render the dashboard', function (): void {
    $this->actingAs(User::factory()->create())
        ->get(route('dashboard'))
        ->assertOk();
});
