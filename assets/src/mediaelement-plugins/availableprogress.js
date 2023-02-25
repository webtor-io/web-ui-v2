Object.assign(MediaElementPlayer.prototype, {
    buildavailableprogress(player, controls, layers, media) {
        const slider = player.slider;
        const el = document.createElement('span');
        el.classList.add(this.options.classPrefix + 'available-progress');
        el.classList.add('bg-accent');
        slider.appendChild(el);
        const cb = function() {
            const a = media.oldGetDuration();
            const t = media.getDuration();
            if (a > 0 && t > 0) {
                const progress = a/t;
                if (progress >= 1) {
                    el.style.display = 'none';
                } else {
                    el.style.transform = `scaleX(${progress})`;
                }
            }
        }
        media.addEventListener('timeupdate', cb);
    },
})