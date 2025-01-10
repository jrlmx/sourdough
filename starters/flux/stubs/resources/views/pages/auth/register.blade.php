<?php

use App\Models\User;
use Illuminate\Auth\Events\Registered;
use Illuminate\Support\Facades\Hash;
use Illuminate\Validation\Rules;

use function Laravel\Folio\middleware;
use function Laravel\Folio\name;
use function Livewire\Volt\rules;
use function Livewire\Volt\state;

name('register');
middleware('guest');

state([
    'name' => '',
    'email' => '',
    'password' => '',
    'password_confirmation' => '',
]);

rules([
    'name' => ['required', 'string', 'max:255', 'min:3'],
    'email' => ['required', 'string', 'email', 'max:255', 'min:3', 'unique:users'],
    'password' => ['required', 'confirmed', 'string', Rules\Password::defaults()],
    'password_confirmation' => ['required', 'string'],
]);

$register = function () {
    $this->validate();

    $user = User::create([
        'name' => $this->name,
        'email' => $this->email,
        'password' => Hash::make($this->password),
    ]);

    event(new Registered($user));

    auth()->login($user);

    $this->redirectIntended(route('dashboard', absolute: false), navigate: true);
};

?>


<x-layouts.auth pg_title="Sign up">
    @volt('auth.register')
        <div>
            <form wire:submit="register" class="space-y-3">
                <flux:input label="Name" wire:model="name" type="text" autofocus />
                <flux:input label="Email" wire:model="email" type="email" placeholder="email@domain.ext" />
                <flux:input label="Password" wire:model="password" type="password" />
                <flux:input label="Confirm Password" wire:model="password_confirmation" type="password" />
                <flux:button type="submit" variant="primary" class="w-full">Sign up</flux:button>
            </form>
        </div>
    @endvolt

    <flux:subheading class="text-center">
        Already have an account?
        <flux:link href="{{ route('login') }}">Sign in</flux:link>
    </flux:subheading>
</x-layouts.auth>
