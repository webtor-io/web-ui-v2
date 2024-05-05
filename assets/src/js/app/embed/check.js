import message from './message';
const sha1 = require('sha1');

function setWidth() {
    const width = document.body.offsetWidth;
    const height = width/16*9;
    document.body.style.height = height + 'px';
}
setWidth();
window.addEventListener('resize', setWidth);
window.addEventListener('click', async () => {
    message.send('init');
    const init = await message.receiveOnce('init');
    const c = await check();
    if (c) {
        initEmbed(init);
    }
});

async function check() {
    message.send('inject', window._checkScript);
    const check = await message.receiveOnce('check');
    return sha1(window._id + check) == _checkHash;
}

function initEmbed(init) {
    const form = document.createElement('form');
    form.setAttribute('method', 'post');
    form.setAttribute('enctype', 'multipart/form-data');
    const csrf = document.createElement('input');
    csrf.setAttribute('name', '_csrf');
    csrf.setAttribute('value', window._CSRF);
    form.append(csrf);
    const i = document.createElement('input');
    i.setAttribute('name', 'settings');
    i.setAttribute('value', JSON.stringify(init));
    form.append(i);
    document.body.append(form);
    // form.setAttribute('action', '/');
    form.submit();
}