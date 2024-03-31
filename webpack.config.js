const path = require('path');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');
const TerserPlugin = require('terser-webpack-plugin');

const fs = require('fs');

function getEntries(path, prefix = '') {
    return new Promise((resolve) => {
        fs.readdir(path, { recursive: true }, (err, files) => {
            const entries = {};
            for (const f of files) {
                if (f.endsWith('.js')) entries[prefix + f.replace('.js', '')] = path + '/' + f;
            }
            resolve(entries);
        });
    })
}

module.exports = async (env, options) => {
    const entries = await getEntries('./assets/src/js/app');
    const devMode = options.mode !== 'production';
    const devEntries = devMode ? await getEntries('./assets/src/js/dev', 'dev/') : {};
    return {
        entry: {
            ...entries,
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
            static: './assets/dist',
            devMiddleware: {
                publicPath: '/assets',
            },
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
                        'style-loader',
                        'css-loader',
                        'postcss-loader'
                    ],
                },
            ]
        },
        plugins: [
            new CopyPlugin({
                patterns: [
                    { from: 'node_modules/mediaelement/build/mejs-controls.svg', to: 'mejs-controls.svg' },
                ],
            }),
        ],
    };
}