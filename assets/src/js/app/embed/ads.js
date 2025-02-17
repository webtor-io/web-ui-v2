import av from '../../lib/av';
import message from './message';


av(async function() {
    if (window._injectAds) {
        console.log('injecting ads');
        message.send('inject', window._injectAds);
    }
    if (window._videoAds) {
        renderVideoAds(this, window._videoAds);
    }
});

function renderVideoAds(el, ads = []) {
    const ad = Object.assign({}, {
        skipDelay: 5,
    }, ads[0]);
    const event = new CustomEvent('ads_play');
    window.dispatchEvent(event);
    const videoEl = document.createElement('video');
    videoEl.src = ad.src;
    videoEl.autoplay = true;
    videoEl.playsInline = true;
    const aEl = document.createElement('a');
    aEl.classList.add('absolute', 'top-0', 'left-0', 'z-50');
    aEl.href = ad.url;
    aEl.target = '_blank';
    aEl.setAttribute('data-umami-event', `embed-${ad.name}-click`);
    aEl.appendChild(videoEl);
    el.appendChild(aEl);
    const skipDelay = ad.skipDelay;
    const closeEl = document.createElement('button');
    closeEl.classList.add('absolute', 'top-2', 'right-2', 'btn', 'btn-accent', 'btn-sm', 'z-50');
    closeEl.textContent = 'Close (' + skipDelay + ')';
    closeEl.disabled = true;
    videoEl.addEventListener('ended', function() {
        aEl.remove();
        closeEl.remove();
        const event = new CustomEvent('ads_close');
        window.dispatchEvent(event);
    });
    videoEl.addEventListener('play', function() {
        if (window.umami) {
            window.umami.track(`embed-${ad.name}-play`)
        }
        el.appendChild(closeEl);
        let cnt = 0;
        const ip = setInterval(function () {
            cnt++;
            if (cnt >= skipDelay) {
                closeEl.disabled = false;
                clearInterval(ip);
                closeEl.textContent = 'Close (X)';
                closeEl.addEventListener('click', function () {
                    aEl.remove();
                    closeEl.remove();
                    const event = new CustomEvent('ads_close');
                    window.dispatchEvent(event);
                });
                return;
            }
            closeEl.textContent = 'Close (' + (skipDelay - cnt) + ')';
        }, 1000);
    });
}
