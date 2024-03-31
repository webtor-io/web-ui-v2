import {init} from '../lib/supertokens';
try {
    const res = await init();
} catch (err) {
    console.log(err);
}
const av = (await import('../lib/asyncView')).initAsyncView;
av(document.querySelector('nav'), 'index', async function() {
    const self = this;
    window.addEventListener('auth', function() {
        self.load();
    }, { once: true });
});

export {}
