import MD5 from "crypto-js/md5";
import {makeDebug} from './debug';
const debug = await makeDebug('webtor:embed:message');
export function initAsyncView(target, name, init, destroy) {
    const layout = target.getAttribute('data-async-layout');
    if (layout) {
        name = name + '_' + MD5(layout).toString();
    }
    const onLoad = function(e) {
        debug(`webtor:async view loaded name=%o`, name);
        if (e && e.detail && e.detail.target) {
            target = e.detail.target;
        }
        target.classList.add('async-view');
        target.classList.add(`async-view-${name}`);
        // In case if async binding has not been invoked
        if (!target.reload) {
            target.reload = function() {
                return new Promise(function(resolve, _) {
                    target.reloadResolve = resolve;
                })
            }
        }
        init.call(target);
    }
    const onDestroy = async (e) => {
        debug(`webtor:async view destroyed name=%o`, name);
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
    let keySuffix = '';
    if (init) keySuffix += init.toString();
    if (destroy) keySuffix += destroy.toString();
    let key = `__async${name}_loaded_${MD5(keySuffix)}`;
    if (!window[key]) {
        window.addEventListener(`async:${name}`, onLoad);
        window.addEventListener(`async:${name}_destroy`, onDestroy);
        window[key] = true;
        onLoad();
    }

}
