import message from './message';
const sha1 = require('sha1');
message.send('init');
const data = await message.receiveOnce('init');
console.log(data);
const c = await check();
if (c) {
    initPlaceholder(data);
    window.addEventListener('click', async () => {
        initEmbed(data);
    }, { once: true });
}

function initPlaceholder(data) {
    if (!data.height) {
        function setHeight() {
            const width = document.body.offsetWidth;
            const height = width / 16 * 9;
            document.body.style.height = height + 'px';
        }
        window.addEventListener('resize', setHeight);
        const s = document.createElement('script');
        s.src = 'assets/lib/iframeResizer.contentWindow.min.js';
        document.body.appendChild(s);
        setHeight();
    } else {
        document.body.style.height = data.height;
    }
    if (data.poster) {
        document.body.style.backgroundImage = 'url(' + data.poster + ')';
        document.body.style.backgroundSize = 'cover';
    }
}

async function check() {
    message.send('inject', window._checkScript);
    const check = await message.receiveOnce('check');
    return sha1(window._id + check) == _checkHash;
}

function initEmbed(data) {
    const form = document.createElement('form');
    form.setAttribute('method', 'post');
    form.setAttribute('enctype', 'multipart/form-data');
    const csrf = document.createElement('input');
    csrf.setAttribute('name', '_csrf');
    csrf.setAttribute('value', window._CSRF);
    form.append(csrf);
    const i = document.createElement('input');
    i.setAttribute('name', 'settings');
    i.setAttribute('value', JSON.stringify(data));
    form.append(i);
    document.body.append(form);
    // form.setAttribute('action', '/');
    form.submit();
}