import av from '../lib/av';
av(async function() {
    if (window.umami && window._tier !== 'free') {
        window.umami.identify({
            tier: window._tier,
        });
    }
    const self = this;
    const themeSelector  = (await import('../lib/themeSelector')).themeSelector;
    themeSelector(this.querySelector('[data-toggle-theme]'));
    window.addEventListener('auth', function() {
        self.reload();
    }, { once: true });
});

export {}

