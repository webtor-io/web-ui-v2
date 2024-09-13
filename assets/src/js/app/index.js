const target = document.currentScript.parentElement;
const av = (await import('../lib/asyncView')).initAsyncView;
av(target, 'index', async function() {
    const dropzone = target.querySelector('.dropzone');
    if (dropzone) {
        const initDrop = (await import('../lib/drop')).initDrop;
        initDrop(dropzone);
    }
    const progress = this.querySelector('.progress-alert');
    if (progress != null) {
        const initProgressLog = (await import('../lib/progressLog')).initProgressLog;
        initProgressLog(progress);
    }
});

export {}