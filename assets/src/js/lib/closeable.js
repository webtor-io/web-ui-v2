export function initCloseable(target) {
    const closeButton = target.querySelector('.closeable-close');
    if (!closeButton) {
        return;
    }
    closeButton.addEventListener('click', function(e) {
        target.classList.remove('popin');
        // alert.classList.add('opacity-0');
        setTimeout(() => { target.classList.add('hidden')}, 200);
    });
}