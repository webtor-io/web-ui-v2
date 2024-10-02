import av from '../../lib/av';
av( function() {
    const query = window.location.hash.replace('#', '');
    const urlParams = new URLSearchParams(query);
    const action = urlParams.get('action');
    const modal = urlParams.get('modal');
    if (!action) return;
    const form = document.querySelector('form.' + action);
    const purgeInput = document.createElement('input');
    purgeInput.setAttribute('type', 'hidden');
    purgeInput.setAttribute('name', 'purge');
    purgeInput.setAttribute('value', 'true');
    form.appendChild(purgeInput);
    form.requestSubmit();
    window.addEventListener('player_ready', function () {
        if (!modal) return;
        const checkbox = document.getElementById(modal + '-checkbox');
        checkbox.checked = true;
    });
});
