export async function makeDebug(name) {
    if (localStorage.debug) {
        const makeDebug = (await import('debug')).default;
        return makeDebug(name);
    } else {
        return function() {};
    }
}