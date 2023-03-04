const path = require('path');
const glob = require('glob')
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');
const { PurgeCSSPlugin } = require('purgecss-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');
// const CompressionPlugin = require('compression-webpack-plugin');

module.exports = (env, options) => {
    const devMode = options.mode !== 'production';
    return {
        // experiments: {
        //     topLevelAwait: true,
        // },
        entry: {
            index: './assets/src/index.js',
            action: './assets/src/action.js',
            player: './assets/src/player.js',
            layout: './assets/src/layout.js',
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
            // proxy: {
            //     '*': {
            //         target: 'http://localhost:8080'
            //     }
            // },
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
            // new CompressionPlugin(),
            new MiniCssExtractPlugin(),
            // new PurgeCSSPlugin({
            //     paths: glob.sync('{templates/**/*,assets/src/**/*.js}',  { nodir: true }),
            // }),
            new CopyPlugin({
                patterns: [
                    { from: 'node_modules/mediaelement/build/mejs-controls.svg', to: 'mejs-controls.svg' },
                ],
            }),
        ],
    };
}