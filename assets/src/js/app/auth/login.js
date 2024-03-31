window.submitLoginForm = function(target, e) {
    (async (data) => {
        const initProgressLog = (await import('../../lib/progressLog')).initProgressLog;
        const pl = initProgressLog(document.querySelector('.progress-alert'));
        pl.clear();
        pl.inProgress('sending magic link to ' + data.email, 'login');
        const supertokens = (await import('../../lib/supertokens'));
        const closeable = document.querySelector('.closeable');
        if (closeable) {
            const initCloseable = (await import('../../lib/closeable')).initCloseable;
            initCloseable(closeable);
        }
        try {
            const resp = await supertokens.sendMagicLink(data);
            pl.done('login');
            pl.finish('magic link sent to ' + data.email);
        } catch (err) {
            console.log(err);
            if (err.statusText) {
                pl.error(err.statusText.toLowerCase(), 'login');
            } else if (err.message) {
                pl.error(err.message.toLowerCase(), 'login');
            } else {
                pl.error('unknown error', 'login');
            }
        }
    })({
        email: target.querySelector('input[name=email]').value,
    });
    e.preventDefault();
    return false;
}
