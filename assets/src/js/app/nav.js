import av from '../lib/av';
av(async function() {
    const self = this;
    const themeSelector  = (await import('../lib/themeSelector')).themeSelector;
    themeSelector(this.querySelector('[data-toggle-theme]'));
    window.addEventListener('auth', function() {
        self.reload();
    }, { once: true });
});

export {}

