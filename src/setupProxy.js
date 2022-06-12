/* This file does not support ES6 format */
/* https://create-react-app.dev/docs/proxying-api-requests-in-development/ */

/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable no-undef */
const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function (app) {
    app.use(
        '/api',
        createProxyMiddleware({
            target: process.env.REACT_APP_PROXY_API_URL,
            logLevel: "debug",
            changeOrigin: true,
        })
    );
};
