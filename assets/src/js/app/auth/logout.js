const target = document.currentScript.parentElement;
const av = (await import('../../lib/asyncView')).initAsyncView;
av(target, 'auth/logout', async function() {
    const initProgressLog = (await import('../../lib/progressLog')).initProgressLog;
    const pl = initProgressLog(target.querySelector('.progress-alert'));
    pl.clear();
    pl.inProgress('logging out', 'logout' );
    const supertokens = (await import('../../lib/supertokens'));
    try {
        await supertokens.logout();
        pl.done('logout');
        pl.finish('logout successful');
        window.dispatchEvent(new CustomEvent('auth'));
    } catch (err) {
        console.log(err);
        if (err.statusText) {
            pl.error(err.statusText.toLowerCase(), 'logout');
        } else if (err.message) {
            pl.error(err.message.toLowerCase(), 'logout');
        } else {
            pl.error('unknown error', 'logout');
        }
    }
});

export {}
