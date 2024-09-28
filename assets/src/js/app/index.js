import av from '../lib/av';
av(async function() {
    const dropzone = this.querySelector('.dropzone');
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