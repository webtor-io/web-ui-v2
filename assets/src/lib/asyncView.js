export default function(name, init, destroy) {
    ['DOMContentLoaded', `async:${name}`].forEach(e => window.addEventListener(e, (e) => {
        let target = document;
        if (e && e.detail && e.detail.target) {
            target = e.detail.target;
        }
        init(target);
    }));
    window.addEventListener(`async:${name}_destroy`, (e) => {
        if (destroy) {
            let target = document;
            if (e && e.detail && e.detail.target) {
                target = e.detail.target;
            }
            destroy(target);
        }
    });

    const event = new CustomEvent(`async:${name}_loaded`);
    window.dispatchEvent(event);
}