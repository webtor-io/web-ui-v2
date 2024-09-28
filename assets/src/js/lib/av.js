export default function (init, destroy = null) {
    const target = document.currentScript.parentElement;
    window.av = window.av || [];
    window.av.push([target, init, destroy]);
}