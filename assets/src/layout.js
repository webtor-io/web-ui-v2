async function replacePage(url) {
    const res = await fetch(url);
    const text = await res.text();
    document.open();
    document.write(text);
    document.close();
    window.replacing = false;
}
if (window.replaced == null) {
    window.replaced = false;
    window.replacing = false;
}
if (!window.replaced && window.location.hash.startsWith('#/')) {
    window.replaced = true;
    window.replacing = true;
    const url = window.location.hash.replace(/^\#/, '') 
    replacePage(url);
}
if (window.location.pathname != '/') {
    window.history.replaceState({
    }, '', '/#' + window.location.pathname);
}

function showProgress() {
    const progress = document.getElementById('progress');
    progress.classList.remove('hidden');
}
function hideProgress() {
    const progress = document.getElementById('progress');
    progress.classList.add('hidden');
}

function hideWrapper() {
    var w = document.getElementById('wrapper');
    w.classList.add('hidden');
}
function showWrapper() {
    if (window.replacing) return;
    var w = document.getElementById('wrapper');
    w.classList.remove('hidden');
}
window.webtor = {
    progress: {
        show: showProgress,
        hide: hideProgress,
    },
    wrapper: {
        show: showWrapper,
        hide: hideWrapper,
    }
}

import './style.css';
import async from './async';

window.addEventListener('DOMContentLoaded', () => {
    showWrapper();
    hideProgress();
    async(['a', 'form'], {
        before: showProgress,
        after: hideProgress,
        fallback: {
            selector: '#content',
            layout: 'async',
        },
    });
});