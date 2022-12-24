import progressLog from './progressLog';
import alert from './alert';
import av from './asyncView';

av('download/post', (target) => {
    alert(target);
    progressLog(target);
});