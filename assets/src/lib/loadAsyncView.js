import executeScriptElements from "./executeScriptElements";
const loadedViews = {};
function loadAsyncView(target, body, template, layout = 'async') {
    const nodes = target.querySelectorAll('.async-loaded');
    const els = [];
    for (const n of nodes) {
        els.push(n);
    }
    els.push(target);
    for (const el of els) {
        const template = el.getAttribute('data-async-template');
        if (!template) continue;
        const detail = {
            target: el,
            layout: el.getAttribute('data-async-layout'),
        };
        const event = new CustomEvent(`async:${template}_destroy`, { detail });
        window.dispatchEvent(event);
    }
    target.innerHTML = body;
    target.classList.add('async-loaded');
    target.setAttribute('data-async-template', template);
    target.setAttribute('data-async-layout', layout);
    executeScriptElements(target);
    const detail = {
        target,
        layout,
    };
    const event = new CustomEvent('async', { detail });
    window.dispatchEvent(event);
    if (template) {
        const event = new CustomEvent('async:'+template, { detail });
        window.dispatchEvent(event);
        if (!loadedViews[template]) {
            window.addEventListener(`async:${template}_loaded`, async function(e) {
                loadedViews[template] = true;
                window.dispatchEvent(event);
            });
        }
    }
}
export default loadAsyncView;
