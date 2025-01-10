<?php

use function Laravel\Folio\middleware;
use function Laravel\Folio\name;

name('dashboard');

middleware(['auth', 'verified']);

?>


<x-layouts.app>
    <h1>Dashboard</h1>
</x-layouts.app>
