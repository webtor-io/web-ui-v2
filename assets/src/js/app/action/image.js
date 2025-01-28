import av from '../../lib/av';

av(() => {
    const event = new CustomEvent('player_ready');
    window.dispatchEvent(event);
});

export {}