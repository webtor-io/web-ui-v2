import av from './asyncView';
import 'mediaelement';
import './player.css';
const {MediaElementPlayer} = global;

function initPlayer(target) {
    const el = target.querySelector('.player');
    const player = new MediaElementPlayer(el, {
        autoRewind: false,
        defaultSeekBackwardInterval: (media) => 15,
        defaultSeekForwardInterval: (media) => 15,
        stretching: 'responsive',
        iconSprite: 'assets/mejs-controls.svg',
        hls: {
            autoStartLoad: true,
            startPosition: 0,
            manifestLoadingTimeOut: 1000 * 60 * 10,
            manifestLoadingMaxRetry: 100,
            manifestLoadingMaxRetryTimeout: 1000 * 10,
            capLevelToPlayerSize: true,
            capLevelOnFPSDrop: true,
            // progressive: true,
            testBandwidth: false,
        },
        async success(media) {
            self.hlsPlayer = media.hlsPlayer;
            if (media.hlsPlayer) {
                media.hlsPlayer.on(Hls.Events.MANIFEST_PARSED, function (event, data) {
                    if (media.hlsPlayer.levels.length > 1) {
                        media.hlsPlayer.startLevel = 1;
                    }
                });
                media.hlsPlayer.on(Hls.Events.ERROR, function (event, data) {
                    if (data.fatal) {
                        switch (data.type) {
                            case Hls.ErrorTypes.NETWORK_ERROR:
                                // try to recover network error
                                debug('fatal network error encountered, try to recover');
                                media.hlsPlayer.startLoad();
                                break;
                            case Hls.ErrorTypes.MEDIA_ERROR:
                                debug('fatal media error encountered, try to recover');
                                media.hlsPlayer.recoverMediaError();
                                break;
                            default:
                                // cannot recover
                                media.hlsPlayer.destroy();
                                break;
                        }
                    }
                });
            }
        },
    });
}

av('action/stream_video', (target) => {
    initPlayer(target);
});

av('action/stream_audio', (target) => {
    initPlayer(target);
});
