import ChromecastPlayer from './chromecast/player';
Object.assign(MediaElementPlayer.prototype, {
    async buildchromecast(player) {
        const t = this;
        window['__onGCastApiAvailable'] = function(isAvailable) {
            if (isAvailable && window.cast) {
                t._initializeCastPlayer();
            }
        };
        if (window.cast) {
            t._initializeCastPlayer();
        } else {
            const s = document.createElement('script');
            s.type = 'text/javascript';
            s.src = 'https://www.gstatic.com/cv/js/sender/v1/cast_sender.js?loadCastFramework=1';
            document.body.appendChild(s);
        }
        player.castButton = document.createElement('div');
        player.castButton.className = `${t.options.classPrefix}button ${t.options.classPrefix}cast-button`;
        player.castButton.innerHTML =
            `<google-cast-launcher id="castbutton"></google-cast-launcher>`;
        player.container.appendChild(player.castButton);
    },
    /**
     *
     * @private
     */
    _initializeCastPlayer () {
        const t = this;

        const
            context = cast.framework.CastContext.getInstance(),
            session = context.getCurrentSession()
        ;

        context.setOptions({
            receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
            autoJoinPolicy: chrome.cast.AutoJoinPolicy.PAGE_SCOPED,
            androidReceiverCompatible: true,
        });

        t.remotePlayer = new cast.framework.RemotePlayer();
        t.remotePlayerController = new cast.framework.RemotePlayerController(t.remotePlayer);
        t.remotePlayerController.addEventListener(cast.framework.RemotePlayerEventType.IS_CONNECTED_CHANGED, t._switchToCastPlayer.bind(this));

        if (session) {
            session.endSession(true);
            // t._switchToCastPlayer();
        }
    },

    /**
     *
     * @private
     */
    _switchToCastPlayer () {
        const t = this;

        if (t.proxy) {
            t.proxy.pause();
        }
        if (cast && cast.framework) {
            if (t.remotePlayer.isConnected) {
                t._setupCastPlayer();
                return;
            }
        }
        t._setDefaultPlayer();
    },
    /**
     *
     * @private
     */
    _setupCastPlayer () {
        const
            t = this,
            context = cast.framework.CastContext.getInstance(),
            castSession = context.getCurrentSession()
        ;

        if (t.loadedChromecast === true) {
            return;
        }

        t.loadedChromecast = true;

        t.proxy = new ChromecastPlayer(t.remotePlayer, t.remotePlayerController, t.media, t.options);

        // If no Cast was setup correctly, make sure it is
        t.media.addEventListener('loadedmetadata', () => {
            if (['SESSION_ENDING', 'SESSION_ENDED', 'NO_SESSION'].indexOf(castSession.getSessionState()) === -1 &&
                t.proxy instanceof DefaultPlayer) {
                t.proxy.pause();
                t.proxy = new ChromecastPlayer(t.remotePlayer, t.remotePlayerController, t.media, t.options);
            }
        });

        t.media.addEventListener('timeupdate', () => {
            t.currentMediaTime = t.getCurrentTime();
        });
    }
})