window.submitLoginForm = function(target, e) {
    (async (data) => {
        const initProgressLog = (await import('../../lib/progressLog')).initProgressLog;
        const pl = initProgressLog(document.querySelector('.progress-alert'));
        pl.clear();
        const e = pl.inProgress('login','sending magic link to ' + data.email);
        const supertokens = (await import('../../lib/supertokens'));
        try {
            await supertokens.sendMagicLink(data, window._CSRF);
            e.done('magic link sent to ' + data.email);
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
    })({
        email: target.querySelector('input[name=email]').value,
    });
    e.preventDefault();
    return false;
}
