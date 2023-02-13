import progressLog from './progressLog';
import alert from './alert';
import av from './asyncView';

av('action/post', (target) => {
    alert(target);
    const progress = target.querySelector('form');
    progressLog(progress, function(ev) {
        if (ev.level == 'rendertag') {
            const p = ev.payload;
            const el = document.createElement(p.tag);
            if (p.tag == 'img') {
                el.setAttribute('src', p.src);
                el.setAttribute('alt', p.alt);
            }
            el.classList.add('hidden');
            el.classList.add('mb-5')
            target.appendChild(el);
            el.onload = function() {
                progress.classList.add('hidden');
                el.classList.remove('hidden');
            }
        }
    });
});