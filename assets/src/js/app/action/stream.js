import av from '../../lib/av';

av(async function() {
    const initPlayer = (await import('../../lib/mediaelement')).initPlayer;
    initPlayer(this);
}, async function() {
    const destroyPlayer = (await import('../../lib/mediaelement')).destroyPlayer;
    destroyPlayer();
});

export {}