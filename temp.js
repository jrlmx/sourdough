const config = {
    // @action|mode:command args
    config: ['foo', 'quiet:bar', 'interact:baz', '@create_project:foo', '@composer_install:quiet:bar'],
    /**
     * "@action": [
     *     "mode:action",
     * ]
     */
    config: {
        '@end': ['foo', 'quiet:bar', 'interact:baz'],
        '@create_project': ['foo'],
        '@composer_install': ['quiet:bar']
    }
}
