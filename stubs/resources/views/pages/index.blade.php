<x-layouts.base>
    <div class="min-h-screen content-center">
        <div class="max-w-xl p-6 mx-auto space-y-6">
            <flux:heading size="xl">Sourdough</flux:heading>

            <p class="prose dark:prose-invert">A highly customizable TALL Stack + Flux UI installer/starter-kit for your next Laravel project.</p>

            <div class="flex justify-between gap-2">
                <a href="https://github.com/jrlmx/sourdough" class="hover:underline">GitHub</a>
                <div class="flex items-center gap-2">
                    @auth
                        <a href="{{ route('dashboard') }}" class="hover:underline">Dashboard</a>
                        <a href="{{ route('logout') }}" class="hover:underline">Logout</a>
                    @else
                        <a href="{{ route('login') }}" class="hover:underline">Login</a>
                        <flux:separator vertical />
                        <a href="{{ route('register') }}" class="hover:underline">Register</a>
                    @endauth
                </div>
            </div>
        </div>
    </div>
</x-layouts.base>
