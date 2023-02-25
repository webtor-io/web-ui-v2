import progressLog from './progressLog';
import alert from './alert';
import av from './asyncView';

av('action/post', (target) => {
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