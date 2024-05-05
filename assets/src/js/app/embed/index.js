import '../../../styles/embed.css';

function setWidth() {
    const width = document.body.offsetWidth;
    const height = width/16*9;
    document.body.style.height = height + 'px';
}
setWidth();
window.addEventListener('resize', setWidth);

window.addEventListener('load', async () => {
    const progress = document.querySelector('.progress-alert');
    const el = document.createElement('div');
    const initProgressLog = (await import('../../lib/progressLog')).initProgressLog;
    initProgressLog(progress, function(ev) {
        if (ev.level != 'rendertemplate') return;
        window.addEventListener('player_ready', function(e) {
            window.removeEventListener('resize', setWidth);
            document.body.style.height = 'auto';
            progress.classList.add('hidden');
            el.classList.remove('hidden');
        }, {once: true});
        el.classList.add('hidden');
        el.classList.add('mb-5')
        document.body.appendChild(el);
        ev.render(el);
    });
});