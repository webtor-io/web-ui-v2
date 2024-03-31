export function initDrop(dropzone) {
    const dropzoneInput = dropzone.querySelector('.dropzone-input');

    dropzone.addEventListener('click', function(e) {
        dropzoneInput.click();
    });

    ['drag', 'dragstart', 'dragend', 'dragover', 'dragenter', 'dragleave', 'drop'].forEach(function(event) {
        document.addEventListener(event, function(e) {
            e.preventDefault();
            e.stopPropagation();
        });
    });

    document.addEventListener('drop', function(e) {
        const dt = e.dataTransfer;
        dropzoneInput.files = dt.files;
        dropzone.requestSubmit();
    });

    dropzoneInput.addEventListener('change', function(e) {
        dropzone.requestSubmit();
    });
}