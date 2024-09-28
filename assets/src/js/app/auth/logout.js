import av from '../../lib/av';
av( async function() {
    const initProgressLog = (await import('../../lib/progressLog')).initProgressLog;
    const pl = initProgressLog(this.querySelector('.progress-alert'));
    pl.clear();
    const e = pl.inProgress('logout', 'logging out');
    const supertokens = (await import('../../lib/supertokens'));
    try {
        await supertokens.logout(window._CSRF);
        e.done('logout successful');
        window.dispatchEvent(new CustomEvent('auth'));
    } catch (err) {
        console.log(err);
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
