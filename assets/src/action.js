import progressLog from './lib/progressLog';
import alert from './lib/alert';
import av from './lib/asyncView';

av('action/post_async', (target) => {
    alert(target);
    const progress = target.querySelector('form');
    const el = document.createElement('div');
    window.addEventListener('player_ready', function(e) {
        progress.classList.add('hidden');
        el.classList.remove('hidden');
    }, {once: true});
    progressLog(progress, function(ev) {
        if (ev.level == 'rendertemplate') {
            el.classList.add('hidden');
            el.classList.add('mb-5')
            target.appendChild(el);
            ev.render(el);
        }
    });
});