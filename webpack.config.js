const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');

module.exports = (env, options) => {
    const devMode = options.mode !== 'production';
    const devEntries = devMode ? {
        "dev/browse": './assets/src/dev/browse.js'
    } : {};
    return {
        entry: {
            index: './assets/src/index.js',
            action: './assets/src/action.js',
            player: './assets/src/player.js',
            layout: './assets/src/layout.js',
            ...devEntries,
        },
        output: {
            filename: '[name].js',
            path: path.resolve(__dirname, 'assets', 'dist'),
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
            minimizer: [
                `...`,
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
                    loader: 'babel-loader',
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
            new MiniCssExtractPlugin(),
            new CopyPlugin({
                patterns: [
                    { from: 'node_modules/mediaelement/build/mejs-controls.svg', to: 'mejs-controls.svg' },
                ],
            }),
        ],
    };
}