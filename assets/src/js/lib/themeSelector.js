export default function(themeSelector) {
    const [darkTheme, lightTheme] = themeSelector.getAttribute('data-toggle-theme').split(',').map((t) => t.trim());
    let currentTheme = window.localStorage.getItem('theme');
    if (currentTheme === null) {
        currentTheme = darkTheme;
        if (window.matchMedia && !window.matchMedia('(prefers-color-scheme: dark)')) {
            currentTheme = lightTheme;
        }
    }
    if (currentTheme === lightTheme) themeSelector.checked = true;
    document.querySelector('html').setAttribute('data-theme', currentTheme);
    window.localStorage.setItem('theme', currentTheme);
    themeSelector.addEventListener('change', (e) => {
        currentTheme = e.target.checked ? lightTheme : darkTheme;
        document.querySelector('html').setAttribute('data-theme', currentTheme);
        window.localStorage.setItem('theme', currentTheme);
    });
}