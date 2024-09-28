import av from '../lib/av';
av(async function() {
    const self = this;
    const progress = self.querySelector('form');
    const el = document.createElement('div');
    const initProgressLog = (await import('../lib/progressLog')).initProgressLog;
    initProgressLog(progress, function(ev) {
        if (ev.level !== 'rendertemplate') return;
        const obj = this;
        window.addEventListener('player_ready', function() {
            obj.skip = true;
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