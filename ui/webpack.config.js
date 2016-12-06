var webpack     = require('webpack');
var webpackhtml = require('html-webpack-plugin');

//const PATH = require('path');
//const ROOT = PATH.resolve(__dirname);

module.exports = {
    entry: {
        'polyfills': './polyfills.ts',
        'vendor':    './vendor.ts',
        'app':       './src/main.ts'
    },
    output: {
        filename: './dist/[name].js'
    },
    devtool: 'source-map',
    resolve: {
        extensions: ['', '.webpack.js', '.web.js', '.ts', '.js']
    },
    module: {
        loaders: [
            {
                test:   /\.ts$/,
                loaders: ['ts-loader', 'angular2-template-loader']
            },
            {
                test: /\.html$/,
                loader: 'html'
            },
            {
               test: /\.css$/,
               loader: 'file?name=dist/assets/css/[name].[ext]'
            },
            {
                test: /\.(woff|woff2|ttf|eot|ico)$/,
                loader: 'file?name=dist/assets/fonts/[name].[hash].[ext]'
            }
        ]
    },
    plugins: [
        new webpack.optimize.CommonsChunkPlugin({
            name: ['app', 'vendor', 'polyfills']
        })
    ]
}
