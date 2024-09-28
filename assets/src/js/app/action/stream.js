import av from '../../lib/av';

function ready() {
    const event = new CustomEvent('player_ready');
    window.dispatchEvent(event);
}

av(async function() {
    const initPlayer = (await import('../../lib/mediaelement')).initPlayer;
    initPlayer(this, ready);
}, async function() {
    const destroyPlayer = (await import('../../lib/mediaelement')).destroyPlayer;
    destroyPlayer();
});

export {}