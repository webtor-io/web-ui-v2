const query = window.location.hash.split('?')[1];
const urlParams = new URLSearchParams(query);
const action = urlParams.get('action');
const modal = urlParams.get('modal');
if (action) {
    window.addEventListener('DOMContentLoaded', function() {
        const el = document.querySelector('form.'+ action);
        el.requestSubmit();
        window.addEventListener('player_ready', function() {
            if (modal == 'subtitles') {
                const checkbox = document.getElementById('subtitles-checkbox');
                checkbox.checked = true;
            }
        })
    });
}
