const target = document.currentScript.parentElement;

function ready() {
    const event = new CustomEvent('player_ready');
    window.dispatchEvent(event);
}


const av = (await import('../../lib/asyncView')).initAsyncView;
av(target, 'action/preview_image', () => {
    ready();
});

export {}