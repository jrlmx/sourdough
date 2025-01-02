<?php

namespace Tests\Feature\Auth;

use Livewire\Volt\Volt;

it('can render the registration page', function (): void {
    $response = $this->get(route('register'))
        ->assertOk()
        ->assertSeeVolt('auth.register');
});

it('can register a new user', function (): void {
    $component = Volt::test('auth.register')
        ->set('name', 'Test User')
        ->set('email', 'test@example.com')
        ->set('password', 'password')
        ->set('password_confirmation', 'password');

    $component->call('register');

    $component->assertHasNoErrors(['name', 'email', 'password', 'password_confirmation']);

    $component->assertRedirect(route('dashboard', absolute: false));

    $this->assertAuthenticated();
});
