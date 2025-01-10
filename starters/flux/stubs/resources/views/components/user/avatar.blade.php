@props([
    'round' => false,
    'size' => 64,
    'text' => 'AB',
    'font_family' => 'sans-serif',
    'font_weight' => 'bold',
    'font_size' => 0.5,
    'bg_color' => 'transparent',
])

<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="{{ $size . 'px' }}"
    height="{{ $size . 'px' }}" viewBox="0 0 {{ $size }} {{ $size }}" version="1.1"
    {{ $attributes->merge(['class' => 'w-full h-full']) }}>
    @if ($round)
        <circle fill="{{ $bg_color }}" cx="{{ $size / 2 }}" cy="{{ $size / 2 }}" r="{{ $size / 2 }}" />
    @else
        <rect fill="{{ $bg_color }}" width="{{ $size }}" height="{{ $size }}" />
    @endif
    <text x="50%" y="50%" style="line-height: 1;font-family: {{ $font_family }}" alignment-baseline="middle"
        text-anchor="middle" font-size="{{ round($size * $font_size) }}" font-weight="{{ $font_weight }}"
        dy=".1em" dominant-baseline="middle" fill="currentColor">
        {{ $text }}
    </text>
</svg>
