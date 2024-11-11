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
import initAsyncView from '../lib/asyncView';

const initTheme = (await import('../lib/themeSelector')).initTheme;
window.loaded = false;
function onLoad() {
    if (window.loaded) return;
    window.loaded = true;
    initTheme(document.querySelector('[data-toggle-theme]'));
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
        update(key, val) {
            if (key === 'title') document.querySelector('title').innerText = val;
        },
        fallback: {
            selector: 'main',
            layout: 'async',
        },
    });
    initAsyncView();
}

if (document.readyState !== 'loading') {
    onLoad();
}
window.addEventListener('DOMContentLoaded', onLoad);
