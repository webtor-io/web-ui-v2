const target = document.currentScript.parentElement;
const av = (await import('../../lib/asyncView')).initAsyncView;
av(target, 'auth/verify', async function() {
    const initProgressLog = (await import('../../lib/progressLog')).initProgressLog;
    const pl = initProgressLog(target.querySelector('.progress-alert'));
    pl.clear();
    pl.inProgress('checking magic link', 'verify');
    const supertokens = (await import('../../lib/supertokens'));
    try {
        const res = await supertokens.handleMagicLinkClicked();
        if (res.status === 'OK') {
            pl.done('verify');
            pl.finish('login successful');
            window.dispatchEvent(new CustomEvent('auth'));
        } else if (res.status === 'RESTART_FLOW_ERROR') {
            pl.error('magic link expired, try to login again', 'verify');
        } else {
            pl.error('login failed, try to login again', 'verify');
        }
    } catch (err) {
        if (err.statusText) {
            pl.error(err.statusText.toLowerCase(), 'verify');
        } else if (err.message) {
            pl.error(err.message.toLowerCase(), 'verify');
        } else {
            pl.error('unknown error', 'verify');
        }
    }
});

export {}
