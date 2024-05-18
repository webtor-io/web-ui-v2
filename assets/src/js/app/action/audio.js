const target = document.currentScript.parentElement;

function ready() {
    const event = new CustomEvent('player_ready');
    window.dispatchEvent(event);
}

const av = (await import('../../lib/asyncView')).initAsyncView;
av(target, 'action/stream_audio', async function() {
    const initPlayer = (await import('../../lib/mediaelement')).initPlayer;
    initPlayer(this, ready);
}, async function() {
    const destroyPlayer = (await import('../../lib/mediaelement')).destroyPlayer;
    destroyPlayer();
});

export {}