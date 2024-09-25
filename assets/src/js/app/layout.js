function showProgress() {
    const progress = document.getElementById('progress');
    progress.classList.remove('hidden');
}
function hideProgress() {
    const progress = document.getElementById('progress');
    progress.classList.add('hidden');
}

window.progress = {
    show: showProgress,
    hide: hideProgress,
};

import {bindAsync} from '../lib/async';
const av = (await import('../lib/asyncView')).initAsyncView;
import themeSelector from "../lib/themeSelector";

function onLoad() {
    av(document.querySelector('nav'), 'index', async function() {
        themeSelector(document.querySelector('[data-toggle-theme]'));
    });
    document.body.style.display = 'flex';
    hideProgress();
    bindAsync({
        async fetch(f, url, fetchParams) {
            showProgress();
            fetchParams.headers['X-CSRF-TOKEN'] = window._CSRF;
            fetchParams.headers['X-SESSION-ID'] = window._sessionID;
            const res = await fetch(url, fetchParams);
            hideProgress();
            return res;
        },
        fallback: {
            selector: 'main',
            layout: 'async',
        },
    });
}

if (document.readyState !== 'loading') {
    onLoad();
}
window.addEventListener('DOMContentLoaded', onLoad);
