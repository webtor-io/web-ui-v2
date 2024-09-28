import {init, refresh} from '../lib/supertokens';
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
