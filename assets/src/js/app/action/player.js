import '../../../styles/player.css';
const target = document.currentScript.parentElement;
function ready() {
    const event = new CustomEvent('player_ready');
    window.dispatchEvent(event);
}


const av = (await import('../../lib/asyncView')).initAsyncView;
av(target, 'action/preview_image', () => {
    ready();
});

for (const format of ['audio', 'video']) {
    av(target, 'action/stream_'+format, async function() {
        const initPlayer = (await import('../../lib/mediaelement')).initPlayer;
        initPlayer(this, ready);
    }, async function() {
        const destroyPlayer = (await import('../../lib/mediaelement')).destroyPlayer;
        destroyPlayer();
    });
}

export {}