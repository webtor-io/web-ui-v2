import {initProgressLog} from "../../lib/progressLog";

initProgressLog(document.querySelector('.test-info'))
    .info('some info message')
    .info('some another info message')
    .close();

initProgressLog(document.querySelector('.test-progress'))
    .inProgress('done', 'something very very very long is done').done()
    .inProgress('progress', 'something very very very long in progress')
    .close();

initProgressLog(document.querySelector('.test-status'))
    .inProgress('done', 'something very very very long is done')
    .updateStatus('test')
    .done()
    .close();

let counter = 1;
const e = initProgressLog(document.querySelector('.test-dynamic-status'))
    .inProgress('progress','something with dynamic status in progress');

setInterval(() => {
    e.updateStatus(counter++);
}, 1000);

initProgressLog(document.querySelector('.test-done-message'))
    .inProgress('done','something very very very long is done')
    .done('some done message')
    .close();

initProgressLog(document.querySelector('.test-warn-message'))
    .inProgress('warn', 'something very very very long has warn')
    .warn('some warning')
    .close();

initProgressLog(document.querySelector('.test-error-message'))
    .inProgress('error', 'something very very very long has error')
    .error('some very very really very long error message')
    .close();

initProgressLog(document.querySelector('.test-error-message-with-extra'))
    .inProgress('error', 'something very very very long has error')
    .error('some very very really very long error message')
    .close();

initProgressLog(document.querySelector('.test-done-with-extra'))
    .inProgress('done', 'something very very very long')
    .done()
    .close();

