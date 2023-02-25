import av from './asyncView';

av('action/stream_audio', (target) => {
    if (window.videojs) {
        const el = target.querySelector('.video-js');
        videojs(el, {
            liveui: true,
            fluid: true,
        });
    }
});
