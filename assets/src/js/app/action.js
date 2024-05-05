const target = document.currentScript.parentElement;
const av = (await import('../lib/asyncView')).initAsyncView;
av(target, 'action/post', async function() {
    const self = this;
    const progress = self.querySelector('form');
    const el = document.createElement('div');
    const closeable = target.querySelector('.closeable');
    if (closeable) {
        const initCloseable = (await import('../lib/closeable')).initCloseable;
        initCloseable(closeable);
    }
    const initProgressLog = (await import('../lib/progressLog')).initProgressLog;
    initProgressLog(progress, function(ev) {
        if (ev.level != 'rendertemplate') return;
        window.addEventListener('player_ready', function(e) {
            progress.classList.add('hidden');
            el.classList.remove('hidden');
        }, {once: true});
        el.classList.add('hidden');
        el.classList.add('mb-5')
        self.appendChild(el);
        ev.render(el);
    });
});

export {}