import av from './lib/asyncView';
import 'mediaelement';
import './mediaelement-plugins/availableprogress';
import './styles/player.css';

const {MediaElementPlayer} = global;

function ready() {
    const event = new CustomEvent('player_ready');
    window.dispatchEvent(event);
}

function initPlayer(target) {
    const el = target.querySelector('.player');
    const duration = el.getAttribute('data-duration') ? parseFloat(el.getAttribute('data-duration')) : -1;
    const features = [
        'playpause',
        'current',
        'progress',
        'duration',
        'tracks',
        'volume',
        'fullscreen',
    ];
    if (duration > 0) {
        features.push('availableprogress');
    }
    const player = new MediaElementPlayer(el, {
        autoRewind: false,
        defaultSeekBackwardInterval: (media) => 15,
        defaultSeekForwardInterval: (media) => 15,
        stretching: 'responsive',
        iconSprite: 'assets/mejs-controls.svg',
        features,
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
            if (duration > 0) {
                const oldGetDuration = media.getDuration;
                media.oldGetDuration = function() {
                    return oldGetDuration.call(media);
                }
                media.getDuration = function() {
                    if (duration > 0) return duration;
                    return this.oldGetDuration();
                }
                const oldSetCurrentTime = player.setCurrentTime;
                player.setCurrentTime = function(time, userInteraction = false) {
                    if (time > media.oldGetDuration()) {
                        return;
                    }
                    return oldSetCurrentTime.call(player, time, userInteraction);
                }
            }
            media.addEventListener('canplay', () => {
                ready();
            });
            if (media.hlsPlayer) {
                media.addEventListener('seeking', () => {
                    if (media.hlsPlayer.loadLevel > 1) {
                        media.hlsPlayer.loadLevel = 1;
                    }
                });
                media.addEventListener('seeked', () => {
                    media.hlsPlayer.loadLevel = -1;
                });
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
    target.player = player;
}

function destroyPlayer(target) {
    if (target.player) {
        target.player.remove();
    }
}

av('action/preview_image_async', (target) => {
    ready();
});

for (const format of ['audio', 'video']) {
    av('action/stream_'+format + '_async', (target) => {
        initPlayer(target);
    }, (target) => {
        destroyPlayer(target);
    });
}