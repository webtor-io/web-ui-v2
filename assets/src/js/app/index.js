const target = document.currentScript.parentElement;
const av = (await import('../lib/asyncView')).initAsyncView;
av(target, 'index', async function() {
    const dropzone = target.querySelector('.dropzone');
    if (dropzone) {
        const initDrop = (await import('../lib/drop')).initDrop;
        initDrop(dropzone);
    }

    const closeable = target.querySelector('.closeable');
    if (closeable) {
        const initCloseable = (await import('../lib/closeable')).initCloseable;
        initCloseable(closeable);
    }
    const progress = this.querySelector('.progress-alert');
    if (progress != null) {
        const initProgressLog = (await import('../lib/progressLog')).initProgressLog;
        initProgressLog(progress);
    }
});

export {}