const path = require('path');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');
const TerserPlugin = require('terser-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const FaviconsWebpackPlugin = require('favicons-webpack-plugin');

const fs = require('fs');

function getEntries(path, ext, prefix = '') {
    return new Promise((resolve) => {
        fs.readdir(path, { recursive: true }, (err, files) => {
            const entries = {};
            for (const f of files) {
                if (f.endsWith(ext)) entries[prefix + f.replace(ext, '')] = path + '/' + f;
            }
            resolve(entries);
        });
    })
}

function makeFavicons(theme) {
    return new FaviconsWebpackPlugin({
        logo: `./assets/src/images/logo-${theme}.svg`,
        prefix: `${theme}/`,
        favicons: {
            icons: {
                android: false,
                appleIcon: false,
                appleStartup: false,
                favicons: true,
                windows: false,
                yandex: false,
            },
        },
    });
}

module.exports = async (env, options) => {
    const jsEntries = await getEntries('./assets/src/js/app', '.js');
    const styleEntries = await getEntries('./assets/src/styles', '.css');
    const devMode = options.mode !== 'production';
    const devEntries = devMode ? await getEntries('./assets/src/js/dev', '.js', 'dev/') : {};
    return {
        entry: {
            ...jsEntries,
            ...styleEntries,
            ...devEntries,
        },
        devtool: 'source-map',
        output: {
            filename: '[name].js',
            chunkFilename: '[name].[chunkhash].js',
            path: path.resolve(__dirname, 'assets', 'dist'),
            clean: true,
        },
        devServer: {
            port: 8082,
            client: {
                webSocketURL: 'auto://0.0.0.0:0/ws',
            },
            static: './assets/dist',
            allowedHosts: ['all'],
            devMiddleware: {
                publicPath: '/assets',
                index: false,
            },
            proxy: [
                {
                    context: () => true,
                    target: 'http://localhost:8080',
                },
            ],
            watchFiles: ['templates/*.html', 'assets/src/*'],
        },
        optimization: {
            // splitChunks: {},
            minimize: true,
            minimizer: [
                new TerserPlugin({ parallel: true }),
                new CssMinimizerPlugin({
                    minimizerOptions: {
                        preset: [
                            "default",
                            {
                                discardComments: { removeAll: true },
                            },
                        ],
                    },
                }),
            ],
        },
        module: {
            rules: [
                {
                    test: /\.js$/,
                    include: path.resolve(__dirname, 'assets', 'src'),
                    loader: 'babel-loader'
                },
                {
                    test: /\.css$/i,
                    include: path.resolve(__dirname, 'assets', 'src'),
                    use: [
                        devMode ? 'style-loader' : MiniCssExtractPlugin.loader,
                        'css-loader',
                        'postcss-loader'
                    ],
                },
            ]
        },
        plugins: [
            new MiniCssExtractPlugin({
                filename: '[name].css',
            }),
            new CopyPlugin({
                patterns: [
                    { from: 'node_modules/mediaelement/build/mejs-controls.svg', to: 'mejs-controls.svg' },
                    { from: 'node_modules/hls.js/dist/hls.min.js', to: 'lib/hls.min.js'},
                    { from: 'node_modules/iframe-resizer/js/iframeResizer.contentWindow.min.js', to: 'lib/iframeResizer.contentWindow.min.js'},
                ],
            }),
            makeFavicons('lofi'),
            makeFavicons('night'),
        ],
    };
}