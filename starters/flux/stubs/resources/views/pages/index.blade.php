<x-layouts.base>
    <div class="min-h-screen content-center">
        <div class="max-w-xl p-6 mx-auto space-y-6">
            <flux:heading size="xl">Sourdough</flux:heading>

            <p class="prose dark:prose-invert">A highly customizable TALL Stack + Flux UI installer/starter-kit for your next Laravel project.</p>

            <div class="flex justify-between items-center gap-6">
                <a href="https://github.com/jrlmx/sourdough" class="hover:underline">GitHub</a>
                <div class="flex items-center gap-6">
                    @auth
                        <flux:button variant="ghost" icon="home" href="{{ route('dashboard') }}" wire:navigate class="hover:underline">Dashboard</flux:button>
                        <livewire:auth.logout icon="arrow-right-start-on-rectangle" />
                    @else
                        <flux:button variant="ghost" href="{{ route('login') }}" icon="arrow-right-start-on-rectangle" wire:navigate>Login</flux:button>
                        <flux:button variant="ghost" href="{{ route('register') }}" icon="arrow-right-start-on-rectangle" wire:navigate>Register</flux:button>
                    @endauth
                </div>
            </div>
        </div>
    </div>
</x-layouts.base>
