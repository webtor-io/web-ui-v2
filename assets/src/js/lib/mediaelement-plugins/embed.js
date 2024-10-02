Object.assign(MediaElementPlayer.prototype, {
    async buildembed(player, controls, layers) {
        player.embedButton = document.createElement('div');
        player.embedButton.className = `${this.options.classPrefix}button ${this.options.classPrefix}embed-button`;
        player.embedButton.innerHTML =
            `<button type="button" role="button" aria-owns="${this.id}" tabindex="0">&lt;<span class="slash">/</span>&gt;</button>`;
        this.addControlElement(player.embedButton, 'embed');
        player.embedLayer = document.createElement('div');
        player.embedLayer.className = `${this.options.classPrefix}layer ${this.options.classPrefix}overlay ${this.options.classPrefix}embed`;
        const embedContainer = document.createElement('div');
        player.embedLayer.appendChild(embedContainer);
        const playLayer = layers.querySelector(`.${this.options.classPrefix}overlay-play`);
        layers.insertBefore(player.embedLayer, playLayer);
        const t = () => {
            const checkbox = document.getElementById('embed-checkbox');
            checkbox.checked = !checkbox.checked;
        }
        player.embedLayer.style.display = 'none';
        player.embedLayer.style.zIndex = 2;
        player.embedButton.addEventListener('click', t);
    },
})