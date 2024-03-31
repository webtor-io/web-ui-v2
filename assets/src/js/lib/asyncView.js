import MD5 from "crypto-js/md5";
export function initAsyncView(target, name, init, destroy) {
    const layout = target.getAttribute('async-layout');
    if (layout) {
        name = name + '_' + MD5(layout).toString();
    }
    const onLoad = function(e) {
        console.log(`async:${name} load`);
        if (e && e.detail && e.detail.target) {
            target = e.detail.target;
        }
        target.classList.add('async-view');
        target.classList.add(`async-view-${name}`);
        init.call(target);
    }
    window.addEventListener(`async:${name}`, onLoad);
    if (document.readyState !== 'complete') {
        onLoad();
    }
    const listener = async (e) => {
        console.log(`async:${name} destroy`);
        const event = new CustomEvent(`async:${name}_destroyed`);
        if (destroy) {
            let target = document;
            if (e && e.detail && e.detail.target) {
                target = e.detail.target;
            }
            await destroy.call(target);
        }
        window.dispatchEvent(event);
    } 
    window.addEventListener(`async:${name}_destroy`, listener);
    const event = new CustomEvent(`async:${name}_loaded`);
    window.dispatchEvent(event);
}
