export default function(name, init) {
    ['DOMContentLoaded', `async:${name}_async`].forEach(e => window.addEventListener(e, (e) => {
        let target = document;
        if (e && e.detail && e.detail.target) {
            target = e.detail.target;
        }
        init(target);
    }));

    const event = new CustomEvent(`async:${name}_async_loaded`);
    window.dispatchEvent(event);
}