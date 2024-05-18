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

import '../../styles/style.css';


import {bindAsync} from '../lib/async';

function onLoad() {
    document.body.style.display = 'block';
    hideProgress();
    bindAsync({
        async fetch(f, url, fetchParams) {
            showProgress();
            fetchParams.headers['X-CSRF-TOKEN'] = window._CSRF;
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
