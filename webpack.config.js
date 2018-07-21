//install yarn
//yarn add siema
//yarn add webpack
//yarn add babel-loader
//yarn add babel-core@6
//yarn add babel-preset-env
//yarn add babel-preset-es2015
//yarn add sw-precache-webpack-plugin
//yarn add babel-plugin-transform-react-jsx
const webpack = require('webpack'); //to access built-in plugins
const path = require('path');

var SWPrecacheWebpackPlugin = require('sw-precache-webpack-plugin');

module.exports = {
  entry: {
    web: "./src/#routes.js",
    mobile: "./src/mobile/#routes.js",
  },
  output: {
    filename: "[name].min.js",
    path: path.resolve(__dirname, "./assets/bin")
  },
  module: {
      rules: [{
          test: /\.js$/,
          exclude: /node_modules/,
          loader: 'babel-loader',
      }]
  },
  plugins: [
    new SWPrecacheWebpackPlugin(
      {
        cacheId: 'uniportv1',
        dontCacheBustUrlsMatching: /\.\w{8}\./,
        filename: 'service-worker.js',
        minify: true,
        navigateFallback: '/',
        staticFileGlobsIgnorePatterns: [/\.map$/, /manifest\.json$/],
      }
    ),
  ],
}
