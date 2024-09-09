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
import {themeChange} from 'theme-change';
themeChange();

function onLoad() {
    let currentTheme = window.localStorage.getItem('theme');
    const themeSelector = document.querySelector('[data-toggle-theme]');
    const [darkTheme, lightTheme] = themeSelector.getAttribute('data-toggle-theme').split(',').map((t) => t.trim());
    if (currentTheme === null) {
        currentTheme = darkTheme;
        if (window.matchMedia && !window.matchMedia('(prefers-color-scheme: dark)')) {
            currentTheme = lightTheme;
        }
    }
    if (currentTheme === lightTheme) themeSelector.checked = true;
    document.querySelector('html').setAttribute('data-theme', currentTheme);
    document.body.style.display = 'block';
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
