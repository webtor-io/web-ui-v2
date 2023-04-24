import drop from './lib/drop';
import alert from './lib/alert';
import progressLog from './lib/progressLog';
import av from './lib/asyncView';

av('index_async', (target) => {
    drop(target);
    alert(target);
    const progress = target.querySelector('.progress-alert');
    if (progress != null) {
        progressLog(progress);
    }
});
