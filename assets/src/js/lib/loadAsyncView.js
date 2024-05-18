import executeScriptElements from "./executeScriptElements";
function loadAsyncView(target, body, template) {
    const nodes = target.querySelectorAll('.async-view');
    const els = [];
    for (const n of nodes) {
        els.push(n);
    }
    let counter = 0;
    for (const el of els) {
        const t = el.getAttribute('async-template');
        if (!t) continue;
        const listener = (e) => {
            counter++;
            if (counter == els.length) {
                renderBody(target, body, template);
            }
        }
        window.addEventListener(`async:${t}_destroyed`, listener, { once: true });
        const detail = {
            target: el,
        };
        const event = new CustomEvent(`async:${t}_destroy`, { detail });
        window.dispatchEvent(event);
    }
    if (els.length == 0) {
        renderBody(target, body, template);
    }
}
function renderBody(target, body, template) {
    target.innerHTML = body;
    target.classList.add('async-loaded');
    target.setAttribute('async-template', template);

    executeScriptElements(target);
    const detail = {
        target,
        template,
    };
    // Update async elements
    const event = new CustomEvent('async', { detail });
    window.dispatchEvent(event);

    // Process async views
    if (template) {
        const event = new CustomEvent('async:' + template, { detail });
        window.dispatchEvent(event);
    }
}

export default loadAsyncView;