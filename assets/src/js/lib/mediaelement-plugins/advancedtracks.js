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
    if (window.hlsPlayer && provider == 'MediaProbe') {
        window.hlsPlayer.audioTrack = e.getAttribute('data-mp-id');
    }
}

window.setSubtitle = function(e) {
    markTrack(e, 'subtitle');
    const videos = document.querySelectorAll('video.player');

    const provider = e.getAttribute('data-provider');
    if (window.hlsPlayer && provider == 'MediaProbe') {
        const id = parseInt(e.getAttribute('data-mp-id'));
        window.hlsPlayer.subtitleTrack = id;
    } else {
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

Object.assign(MediaElementPlayer.prototype, {
    async buildadvancedtracks(player, controls, layers) {
        player.tracksButton = document.createElement('div');
        player.tracksButton.className = `${this.options.classPrefix}button ${this.options.classPrefix}captions-button`;
        player.tracksButton.innerHTML =
            `<button type="button" role="button" aria-owns="${this.id}" tabindex="0">
                <svg xmlns="http://www.w3.org/2000/svg" id="mep_0-icon-captions" class="mejs__icon-captions" aria-hidden="true" focusable="false">
                    <use xlink:href="assets/mejs-controls.svg#icon-captions"></use>
                </svg>
            </button>`;
        this.addControlElement(player.tracksButton, 'tracks');
        player.tracksLayer = document.createElement('div');
        player.tracksLayer.className = `${this.options.classPrefix}layer ${this.options.classPrefix}overlay ${this.options.classPrefix}tracks`;
        const tracksContainer = document.getElementById('subtitles');
        const checkbox = document.getElementById('subtitles-checkbox');
        player.tracksLayer.appendChild(checkbox);
        player.tracksLayer.appendChild(tracksContainer);
        const playLayer = layers.querySelector(`.${this.options.classPrefix}overlay-play`);

        layers.insertBefore(player.tracksLayer, playLayer);
        const t = (e) => {
            e.preventDefault();
            checkbox.checked = !checkbox.checked;
            if (checkbox.checked) {
                player.tracksLayer.style.display = 'block';
            } else {
                player.tracksLayer.style.display = 'none';
            }
            return false;
        }
        const close = tracksContainer.querySelector('label[for=subtitles-checkbox]');
        player.tracksLayer.style.display = 'none';
        player.tracksLayer.style.zIndex = 1000;
        player.tracksButton.addEventListener('click', t);
        close.addEventListener('click', t);

    },
})