import {init, refresh} from '../lib/supertokens';
const av = (await import('../lib/asyncView')).initAsyncView;

av(document.querySelector('nav'), 'index', async function() {
    const self = this;
    window.addEventListener('auth', function() {
        self.reload();
    }, { once: true });
});

try {
    await init(window._CSRF);
} catch (err) {
    console.log(err);
}
if (window._sessionExpired) {
    try {
        await refresh(window._CSRF);
        window.dispatchEvent(new CustomEvent('auth'));
    } catch (err) {
        console.log(err);
    }
}

export {}
