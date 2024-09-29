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
initAsyncView();

const initTheme = (await import('../lib/themeSelector')).initTheme;
function onLoad() {
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
}

if (document.readyState !== 'loading') {
    onLoad();
}
window.addEventListener('DOMContentLoaded', onLoad);
