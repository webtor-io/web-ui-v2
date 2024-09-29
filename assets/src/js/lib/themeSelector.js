const storageKey = 'theme';
function getThemes(themeSelector) {
    return themeSelector.getAttribute('data-toggle-theme').split(',').map((t) => t.trim());
}
function changeFavicons(current) {
    const els = document.querySelectorAll('link[rel~="icon"]')
    for (const el of els) {
        el.href = el.href.replace(/\w+\/favicon/, `${current}/favicon`);
    }
}
export function initTheme(themeSelector) {
    const [darkTheme, lightTheme] = getThemes(themeSelector);
    let currentTheme = window.localStorage.getItem(storageKey);
    if (currentTheme === null) {
        currentTheme = darkTheme;
        if (window.matchMedia && !window.matchMedia('(prefers-color-scheme: dark)')) {
            currentTheme = lightTheme;
        }
    }
    if (currentTheme === lightTheme) {
        themeSelector.checked = true;
        changeFavicons(currentTheme);
    }
    document.querySelector('html').setAttribute('data-theme', currentTheme);
    window.localStorage.setItem(storageKey, currentTheme);
}
export function themeSelector(themeSelector) {
    const [darkTheme, lightTheme] = getThemes(themeSelector);
    let currentTheme = window.localStorage.getItem(storageKey);
    if (currentTheme === lightTheme) themeSelector.checked = true;
    themeSelector.addEventListener('change', (e) => {
        currentTheme = e.target.checked ? lightTheme : darkTheme;
        document.querySelector('html').setAttribute('data-theme', currentTheme);
        changeFavicons(currentTheme);
        window.localStorage.setItem(storageKey, currentTheme);
    });
}