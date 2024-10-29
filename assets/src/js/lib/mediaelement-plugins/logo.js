Object.assign(MediaElementPlayer.prototype, {
    async buildlogo(player, controls, layers) {
        const playLayer = layers.querySelector(`.${this.options.classPrefix}overlay-play`);
        const logoButton = document.createElement('div');
        logoButton.className = `${this.options.classPrefix}button ${this.options.classPrefix}logo-button`;
        logoButton.innerHTML = `
            <div class="font-baskerville sm:text-5xl text-4xl absolute right-0 top-0 p-3 sm:p-4 z-50 text-center">
                <a href="#" target="_blank">
                    <span>web</span><span class="text-accent">tor</span>
                </a>
            </div>`;
        logoButton.querySelector('a').setAttribute('href', window._embedSettings.baseUrl);
        playLayer.appendChild(logoButton);
    }
});