import executeScriptElements from "./executeScriptElements";
const loadedViews = {};
function loadAsyncView(target, body, template, layout = 'async') {
    target.innerHTML = body;
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
