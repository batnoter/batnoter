/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable no-undef */
const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function (app) {
    app.use(
        '/api',
        createProxyMiddleware({
            target: process.env.REACT_APP_PROXY_API_URL,
            secure: false,
            logLevel: "debug",
            changeOrigin: true,
        })
    );
};
