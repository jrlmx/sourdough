<?php

namespace Tests\Feature;

use App\Models\User;
use Livewire\Volt\Volt;

it('can render profile page', function (): void {
    $this->actingAs(User::factory()->create());

    $this->get(route('profile'))
        ->assertOk()
        ->assertSeeVolt('profile.update_account')
        ->assertSeeVolt('profile.update_password')
        ->assertSeeVolt('profile.delete_account');
});

it('can update account', function (): void {
    $user = User::factory()->create();

    $this->actingAs($user);

    $component = Volt::test('profile.update_account')
        ->set('name', 'Test User')
        ->set('email', 'test@example.com')
        ->call('updateAccount')
        ->assertHasNoErrors()
        ->assertNoRedirect();

    $user->refresh();

    $this->assertEquals('Test User', $user->name);
    $this->assertEquals('test@example.com', $user->email);
    $this->assertNull($user->email_verified_at);
});

it('does not change email verification status if email is not changed', function (): void {
    $user = User::factory()->create([
        'email' => 'test@example.com',
    ]);

    $this->actingAs($user);

    $component = Volt::test('profile.update_account')
        ->set('name', 'Test User')
        ->set('email', 'test@example.com')
        ->call('updateAccount')
        ->assertHasNoErrors()
        ->assertNoRedirect();

    $this->assertNotNull($user->refresh()->email_verified_at);
});

it('can delete account', function (): void {
    $user = User::factory()->create();

    $this->actingAs($user);

    $component = Volt::test('profile.delete_account')
        ->set('current_password', 'password')
        ->call('deleteAccount')
        ->assertHasNoErrors()
        ->assertRedirect(route('login'));

    $this->assertNull(User::find($user->id));
});

it('cannot delete account with invalid password', function (): void {
    $user = User::factory()->create();

    $this->actingAs($user);

    $component = Volt::test('profile.delete_account')
        ->set('current_password', 'wrong-password')
        ->call('deleteAccount')
        ->assertHasErrors(['current_password'])
        ->assertNoRedirect();

    $this->assertNotNull($user->fresh());
});
