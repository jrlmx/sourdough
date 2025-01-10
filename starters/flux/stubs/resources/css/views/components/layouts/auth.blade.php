@props(['app_title' => null, 'pg_title' => null])

<x-layouts.base :$app_title :$pg_title>
    <div class="min-h-screen content-center">
        <flux:card class="mx-auto max-w-md space-y-3">
            {{ $slot }}
        </flux:card>
    </div>
</x-layouts.base>
