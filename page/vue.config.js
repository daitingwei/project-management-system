    const path = require('path');

    function resolve(dir) {
        return path.join(__dirname, dir)
    }

    module.exports = {
        publicPath: process.env.NODE_ENV === 'production' ? process.env.VUE_APP_BUILD_PATH: '/',
        productionSourceMap: false,
        css: {
            loaderOptions: {
                less: {
                    modifyVars: {
                        blue: '#3a82f8',
                        'text-color': '#333'
                    },
                    javascriptEnabled: true
                }
            }
        },
        devServer: {
            host: process.env.VUE_APP_DEV_HOST || '0.0.0.0',
            port: process.env.VUE_APP_DEV_PORT || '8045',
            https: false,
            hot: true,
            open: false,
            clientLogLevel: 'none',
            hotOnly: false,
            historyApiFallback: true,
            allowedHosts: ['.ngrok.io', '.tunnel'],
            disableHostCheck: true,
            before(app) {
                app.disable('etag')
            },
            proxy: {
                '/api': {
                    target: `${process.env.VUE_APP_API_URL}/`,
                    ws: true,
                    changOrigin: true,
                    secure: false,
                    pathRewrite: {
                        '^/api': ''
                    }
                }
            },
        },
        configureWebpack: config => {
            config.resolve = {
                extensions: ['.js', '.vue', '.json', ".css"],
                alias: {
                    'vue$': 'vue/dist/vue.esm.js',
                    '@': resolve('src'),
                    'assets': resolve('src/assets'),
                    'components': resolve('src/components')
                }
            }
        },
        lintOnSave: false
    };
