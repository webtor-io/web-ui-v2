import message from './message';
import initAsyncView from '../../lib/asyncView';

if (window._umami) {
    (await import('../../lib/umami')).init(window, Object.assign(window._umami, {
        tag: 'embed',
    }));
}

function setHeight() {
    const width = document.body.offsetWidth;
    const height = width/16*9;
    document.body.style.height = height + 'px';
}

window.addEventListener('load', async () => {
    const progress = document.querySelector('.progress-alert');
    const el = document.createElement('div');
    const initProgressLog = (await import('../../lib/progressLog')).initProgressLog;
    initProgressLog(progress, function(ev) {
        if (ev.level !== 'rendertemplate') return;
        window.addEventListener('player_ready', function() {
            window.removeEventListener('resize', setHeight);
            document.body.style.height = 'auto';
            progress.classList.add('hidden');
            el.classList.remove('hidden');
            if (window._domainSettings.ads && window._injectAds) {
                console.log('injecting ads');
                message.send('inject', window._injectAds);
            }
        }, {once: true});
        el.classList.add('hidden');
        el.classList.add('mb-5')
        document.body.appendChild(el);
        ev.render(el);
    });
    if (!window._embedSettings.height) {
        const s = document.createElement('script');
        s.src = 'assets/lib/iframeResizer.contentWindow.min.js';
        document.body.appendChild(s);
    }
    setHeight();
    window.addEventListener('resize', setHeight);
    if (window._embedSettings.poster) {
        document.body.style.backgroundImage = 'url(' + window._embedSettings.poster + ')';
        document.body.style.backgroundSize = 'cover';
    }
    initAsyncView();
});