import progressLog from './progressLog';
import alert from './alert';
import av from './asyncView';

av('action/post', (target) => {
    alert(target);
    const progress = target.querySelector('form');
    progressLog(progress, function(ev) {
        if (ev.level == 'rendertemplate') {
            const el = document.createElement('div');
            el.classList.add('hidden');
            el.classList.add('mb-5')
            target.appendChild(el);
            ev.render(el);
            progress.classList.add('hidden');
            el.classList.remove('hidden');
        }
    });
});