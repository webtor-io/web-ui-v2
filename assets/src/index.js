import drop from './index/drop';
import alert from './alert';
import progressLog from './progressLog';
import av from './asyncView';

av('index', (target) => {
    drop(target);
    alert(target);
    progressLog(target);
});