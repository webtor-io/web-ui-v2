import loadAsyncView from "./loadAsyncView";
export function initProgressLog(el, func) {
    const lt = el.querySelector('.log-target');
    const url = el.getAttribute('async-progress-log');
    const obj = {
        inited: false,
    };
    if (url) {
        const src = new EventSource(url, {
            withCredentials: true,
        });
        src.onmessage = (ev) => {
            const data = JSON.parse(ev.data);
            onMessage.call(obj, el, data, lt, func);
            if (['finish', 'close', 'redirect', 'download', 'error'].includes(data.level)) {
                src.close();
            }
        };
    }
    return {
        clear() {
            this.push({level:   'clear'});
        },
        inProgress(message, tag) {
            this.push({level: 'inprogress', message, tag});
        },
        done(tag) {
            this.push({level: 'done', tag});
        },
        finish(message) {
            this.push({level: 'finish', message});
        },
        error(message, tag) {
            this.push({level: 'error', message, tag});
        },
        push(data) {
            onMessage.call(obj, el, data, lt, func);
        }
    };
}

function onMessage(el, data, lt, func) {
    if (data.level == 'clear') {
        lt.innerText = '';
        el.classList.add('hidden');
        el.querySelector('.alert-close-wrapper').classList.add('hidden');
        this.inited = false;
    } else if (!this.inited) {
        this.inited = true;
        el.classList.remove('hidden');
    }
    if (data.level == 'finish') {
        const pre = document.createElement('pre');
        pre.innerText = data.message;
        pre.classList.add('done-summary');
        lt.appendChild(pre);
        el.querySelector('.alert-close-wrapper').classList.remove('hidden');
    }
    if (data.level == 'redirect') {
        const pre = document.createElement('pre');
        pre.classList.add('in-progress');
        pre.innerText = 'success! redirecting';
        pre.classList.add('done-summary');
        lt.appendChild(pre);
        el.setAttribute('action', data.location);
        el.requestSubmit();
    }
    if (data.level == 'rendertemplate') {
        data.render = (el) => {
            loadAsyncView(el, data.body, data.template);
        };
    }
    if (data.level == 'download') {
        const pre = document.createElement('pre');
        pre.innerText = 'success! download should start right now!';
        pre.classList.add('done-summary');
        lt.appendChild(pre);
        window.location = data.location;
        el.querySelector('.alert-close-wrapper').classList.remove('hidden');
    }
    if (data.level == 'info') {
        const pre = document.createElement('pre');
        pre.innerText = data.message;
        lt.appendChild(pre);
    }
    if (data.level == 'inprogress') {
        const pre = document.createElement('pre');
        pre.innerText = data.message;
        pre.classList.add('in-progress');
        pre.setAttribute('task-tag', data.tag);
        lt.appendChild(pre);
    }
    if (data.level == 'done') {
        const el = lt.querySelector(`*[task-tag='${data.tag}']`);
        el.classList.remove('in-progress');
        let span = el.querySelector('span');
        if (!span) {
            span = document.createElement('span');
            el.appendChild(span);
        }
        span.innerText = '...[done]';
    }
    if (data.level == 'statusupdate') {
        const el = lt.querySelector(`*[task-tag='${data.tag}']`);
        let span = el.querySelector('span');
        if (!span) {
            span = document.createElement('span');
            el.appendChild(span);
        }
        span.innerText = ' (' + data.message + ')';
    }
    if (data.level == 'warn') {
        const tel = lt.querySelector(`*[task-tag='${data.tag}']`);
        tel.classList.remove('in-progress');
        let span = tel.querySelector('span');
        if (!span) {
            span = document.createElement('span');
            tel.appendChild(span);
        }
        span.innerText = '...[error]';
        const pre = document.createElement('pre');
        pre.innerText = data.message;
        pre.classList.add('error-summary');
        lt.appendChild(pre);
    }
    if (data.level == 'error') {
        const tel = lt.querySelector(`*[task-tag='${data.tag}']`);
        tel.classList.remove('in-progress');
        let span = tel.querySelector('span');
        if (!span) {
            span = document.createElement('span');
            tel.appendChild(span);
        }
        span.innerText = '...[error]';
        const pre = document.createElement('pre');
        pre.innerText = data.message;
        pre.classList.add('error-summary');
        lt.appendChild(pre);
        el.querySelector('.alert-close-wrapper').classList.remove('hidden');
        const errorEl = el.querySelector('.alert-close-wrapper .error');
        if (errorEl) errorEl.classList.remove('hidden');
    }
    if (func !== undefined) {
        func(data);
    }
}
