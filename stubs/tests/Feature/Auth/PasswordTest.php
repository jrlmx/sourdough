<?php

namespace Tests\Feature\Auth;

use App\Models\User;
use Illuminate\Auth\Notifications\ResetPassword;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\Notification;
use Livewire\Volt\Volt;

it('can render confirm password screen', function (): void {
    $this->actingAs(User::factory()->create())
        ->get(route('password.confirm'))
        ->assertOk()
        ->assertSeeVolt('password.confirm');
});

it('can confirm password', function (): void {
    $this->actingAs(User::factory()->create());

    $component = Volt::test('password.confirm')
        ->set('password', 'password')
        ->call('confirm');

    $component->assertRedirect(route('dashboard', absolute: false))->assertHasNoErrors();
});

it('cannot confirm password with invalid password', function (): void {
    $this->actingAs(User::factory()->create());

    $component = Volt::test('password.confirm')
        ->set('password', 'wrong-password')
        ->call('confirm');

    $component
        ->assertNoRedirect()
        ->assertHasErrors(['password']);
});

it('can render password reset link screen', function (): void {
    $this->get(route('password.request'))
        ->assertOk()
        ->assertSeeVolt('password.request');
});

it('can send password reset link', function (): void {
    Notification::fake();

    $user = User::factory()->create();

    $component = Volt::test('password.request')
        ->set('email', $user->email)
        ->call('sendResetLink');

    Notification::assertSentTo($user, ResetPassword::class);
});

it('can render password reset screen', function (): void {
    Notification::fake();

    $user = User::factory()->create();

    Volt::test('password.request')
        ->set('email', $user->email)
        ->call('sendResetLink');

    Notification::assertSentTo($user, ResetPassword::class, function ($notification) {
        $this->get(route('password.reset', ['token' => $notification->token]))
            ->assertOk()
            ->assertSeeVolt('password.reset');

        return true;
    });
});

it('can reset password with valid token', function (): void {
    Notification::fake();

    $user = User::factory()->create();

    Volt::test('password.request')
        ->set('email', $user->email)
        ->call('sendResetLink');

    Notification::assertSentTo($user, ResetPassword::class, function ($notification) use ($user) {
        $component = Volt::test('password.reset', ['token' => $notification->token])
            ->set('email', $user->email)
            ->set('password', 'password')
            ->set('password_confirmation', 'password')
            ->call('resetPassword')
            ->assertHasNoErrors()
            ->assertRedirect(route('login', absolute: false));

        return true;
    });
});

it('can update password', function (): void {
    $user = User::factory()->create();

    $this->actingAs($user);

    $component = Volt::test('profile.update_password')
        ->set('current_password', 'password')
        ->set('new_password', 'new-password')
        ->set('new_password_confirmation', 'new-password')
        ->call('updatePassword')
        ->assertHasNoErrors()
        ->assertNoRedirect();

    $this->assertTrue(Hash::check('new-password', $user->refresh()->password));
});

it('cannot update password with invalid current password', function (): void {
    $user = User::factory()->create();

    $this->actingAs($user);

    $component = Volt::test('profile.update_password')
        ->set('current_password', 'wrong-password')
        ->set('new_password', 'new-password')
        ->set('new_password_confirmation', 'new-password')
        ->call('updatePassword')
        ->assertHasErrors(['current_password'])
        ->assertNoRedirect();
});
