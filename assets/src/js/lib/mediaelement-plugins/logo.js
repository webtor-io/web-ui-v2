Object.assign(MediaElementPlayer.prototype, {
    watchControlsVisible(controls, callback) {
        const observer = new MutationObserver((mutations) => {
            mutations.forEach((mutation) => {
                if (mutation.attributeName === 'class') {
                    callback(!mutation.target.classList.contains(`${this.options.classPrefix}offscreen`));
                }
            });
        });
        observer.observe(controls, {
            attributes: true,
        });
    },
    async buildlogo(player, controls, layers) {
        const logoLayer = document.createElement('div');
        layers.logoLayer = logoLayer;
        logoLayer.className = `${this.options.classPrefix}layer ${this.options.classPrefix}webtor-logo`;
        const logo = document.getElementById('logo').cloneNode(true);
        logo.classList.remove('hidden');
        logoLayer.appendChild(logo);
        layers.appendChild(logoLayer);
        this.watchControlsVisible(controls, (visible) => {
            if (visible) {
                logoLayer.classList.remove('hidden');
            } else {
                logoLayer.classList.add('hidden');
            }
        });
    }
});