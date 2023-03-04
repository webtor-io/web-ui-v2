export default function(target) {
    const alert = target.querySelector('#alert');
    const closeButton = target.querySelector('#alert-close');
    if (!closeButton) {
        return;
    }
    closeButton.addEventListener('click', function(e) {
        alert.classList.remove('popin');
        // alert.classList.add('opacity-0');
        setTimeout(() => { alert.remove( )}, 200);
    });
}