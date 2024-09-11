const target = document.currentScript.parentElement;
const av = (await import('../../lib/asyncView')).initAsyncView;
av(target, 'auth/verify', async function() {
    const initProgressLog = (await import('../../lib/progressLog')).initProgressLog;
    const pl = initProgressLog(target.querySelector('.progress-alert'));
    pl.clear();
    const e = pl.inProgress('verify', 'checking magic link');
    const supertokens = (await import('../../lib/supertokens'));
    try {
        const res = await supertokens.handleMagicLinkClicked(window._CSRF);
        if (res.status === 'OK') {
            e.done('login successful');
            window.dispatchEvent(new CustomEvent('auth'));
        } else if (res.status === 'RESTART_FLOW_ERROR') {
            e.error('magic link expired, try to login again');
        } else {
            e.error('login failed, try to login again');
        }
    } catch (err) {
        if (err.statusText) {
            e.error(err.statusText.toLowerCase());
        } else if (err.message) {
            e.error(err.message.toLowerCase());
        } else {
            e.error('unknown error');
        }
    }
    e.close();
});

export {}
