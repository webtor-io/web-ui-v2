import av from '../../lib/av';

function ready() {
    const event = new CustomEvent('player_ready');
    window.dispatchEvent(event);
}


av(() => {
    ready();
});

export {}