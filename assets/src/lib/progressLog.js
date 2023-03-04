import loadAsyncView from "./loadAsyncView";
export default function(el, func) {
    const url = el.getAttribute('async-progress-log');
    const lt = el.querySelector('#log-target');
    const src = new EventSource(url, {
        withCredentials: true,
    });
    let hidden = true;
    src.onmessage = (ev) => {
        if (hidden) {
            hidden = false;
            el.classList.remove('hidden');
        }
        const data = JSON.parse(ev.data);
        if (data.level == 'close') {
            src.close();
        }
        if (data.level == 'finish') {
            src.close();
            const pre = document.createElement('pre');
            pre.innerText = data.message;
            pre.classList.add('done-summary')
            lt.appendChild(pre);
            el.querySelector('#alert-close-wrapper').classList.remove('hidden');
        }
        if (data.level == 'redirect') {
            src.close();
            const pre = document.createElement('pre');
            pre.classList.add('in-progress')
            pre.innerText = 'success! redirecting'
            pre.classList.add('done-summary')
            lt.appendChild(pre);
            el.setAttribute('download', data.location);
            el.requestSubmit();
        }
        if (data.level == 'rendertemplate') {
            data.render = (el) => {
                loadAsyncView(el, data.body, data.template);
            }
        }
        if (data.level == 'download') {
            src.close();
            const pre = document.createElement('pre');
            pre.innerText = 'success! download should start right now!'
            pre.classList.add('done-summary')
            lt.appendChild(pre);
            window.location = data.location;
            el.querySelector('#alert-close-wrapper').classList.remove('hidden');
        }
        if (data.level == 'info') {
            const pre = document.createElement('pre');
            pre.innerText = data.message;
            lt.appendChild(pre);
        }
        if (data.level == 'inprogress') {
            const pre = document.createElement('pre');
            pre.innerText = data.message;
            pre.classList.add('in-progress')
            pre.setAttribute('task-tag', data.tag);
            lt.appendChild(pre);
        }
        if (data.level == 'done') {
            const el = lt.querySelector(`*[task-tag='${data.tag}']`);
            el.classList.remove('in-progress')
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
        if (data.level == 'error') {
            src.close();
            const tel = lt.querySelector(`*[task-tag='${data.tag}']`);
            tel.classList.remove('in-progress')
            let span = tel.querySelector('span');
            if (!span) {
                span = document.createElement('span');
                tel.appendChild(span);
            }
            span.innerText = '...[error]';
            const pre = document.createElement('pre');
            pre.innerText = data.message;
            pre.classList.add('error-summary')
            lt.appendChild(pre);
            el.querySelector('#alert-close-wrapper').classList.remove('hidden');
        }
        if (func !== undefined) {
            func(data);
        }
    };
}