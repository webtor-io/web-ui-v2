import av from './lib/asyncView';
import 'mediaelement';
import './mediaelement-plugins/availableprogress';
import './mediaelement-plugins/advancedtracks';
import './styles/player.css';

window.toggleOpenSubtitles = function(e) {
    const el = document.getElementById('opensubtitles');
    const ele = document.getElementById('embedded');
    const hidden = el.classList.contains('hidden');
    if (hidden) {
        e.classList.remove('btn-outline');
        el.classList.remove('hidden');
        ele.classList.add('hidden');
    } else {
        e.classList.add('btn-outline');
        el.classList.add('hidden');
        ele.classList.remove('hidden');
    }
}

async function markTrack(e, type) {
    if (e.getAttribute('data-default') == 'true') {
        return;
    }
    e.classList.add('text-primary', 'underline');
    e.setAttribute('data-default', 'true');
    const s = document.getElementById('subtitles');
    const es = s.querySelectorAll(`.${type}`);
    for (const ee of es) {
        if (ee == e) continue;
        ee.classList.remove('text-primary', 'underline');
        ee.removeAttribute('data-default');
    }
    const csrf = s.getAttribute('data-csrf');
    await fetch(`/stream-video/${type}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'X-CSRF-TOKEN': s.getAttribute('data-csrf'),
        },
        body:   JSON.stringify({
            id:         e.getAttribute('data-id'),
            resourceID: s.getAttribute('data-resource-id'),
            itemID:     s.getAttribute('data-item-id'),
        }),
    });
}

window.setAudio = function(e) {
    markTrack(e, 'audio');
    const provider = e.getAttribute('data-provider');
    if (hlsPlayer && provider == 'MediaProbe') {
        hlsPlayer.audioTrack = e.getAttribute('data-mp-id');
    }
}

window.setSubtitle = function(e) {
    markTrack(e, 'subtitle');
    const videos = document.querySelectorAll('video.player');

    const provider = e.getAttribute('data-provider');
    if (hlsPlayer && provider == 'MediaProbe') {
        const id = parseInt(e.getAttribute('data-mp-id'));
        hlsPlayer.subtitleTrack = id;
    } else if (video) {
        const id = e.getAttribute('data-id');
        for (const p of videos) {
            for (const t of p.textTracks) {
                if (t.id ==  id) {
                    t.mode = 'showing';
                } else {
                    t.mode = 'hidden';
                }
            }
        }
    }
}

const {MediaElementPlayer} = global;

function ready() {
    const event = new CustomEvent('player_ready');
    window.dispatchEvent(event);
}

let player;
let hlsPlayer;
let video;

function initPlayer(target) {
    video = target.querySelector('.player');
    const duration = video.getAttribute('data-duration') ? parseFloat(video.getAttribute('data-duration')) : -1;
    const features = [
        'playpause',
        'current',
        'progress',
        'duration',
        'advancedtracks',
        'volume',
        'fullscreen',
    ];
    if (duration > 0) {
        features.push('availableprogress');
    }
    player = new MediaElementPlayer(video, {
        autoRewind: false,
        defaultSeekBackwardInterval: (media) => 15,
        defaultSeekForwardInterval: (media) => 15,
        stretching: 'responsive',
        iconSprite: 'assets/mejs-controls.svg',
        features,
        hls: {
            // debug: true,
            autoStartLoad: true,
            startPosition: 0,
            manifestLoadingTimeOut: 1000 * 60 * 10,
            manifestLoadingMaxRetry: 100,
            manifestLoadingMaxRetryTimeout: 1000 * 10,
            capLevelToPlayerSize: true,
            capLevelOnFPSDrop: true,
            // progressive: true,
            testBandwidth: false,
            path: 'https://cdn.jsdelivr.net/npm/hls.js@1.4.0',
        },
        error: function(e) {
            console.log(e);
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
                if (hlsPlayer) {
                    const audioId = document.querySelector('.audio[data-default=true]').getAttribute('data-mp-id');
                    const subId = document.querySelector('.subtitle[data-default=true]').getAttribute('data-mp-id');
                    if (audioId) hlsPlayer.audioTrack = audioId;
                    if (subId) hlsPlayer.subtitleTrack = subId;
                }
                ready();
            });
            if (media.hlsPlayer) {
                hlsPlayer = media.hlsPlayer;
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
}

function destroyPlayer() {
    if (player) {
        player.remove();
        player = null;
    }
    if (hlsPlayer) {
        hlsPlayer.destroy();
        hlsPlayer = null;
    }
}

av('action/preview_image_async', (target) => {
    ready();
});

for (const format of ['audio', 'video']) {
    av('action/stream_'+format + '_async', (target) => {
        initPlayer(target);
    }, (target) => {
        destroyPlayer();
    });
}