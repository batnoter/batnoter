import { createProxyMiddleware } from 'http-proxy-middleware';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
module.exports = function (app: any) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: process.env.REACT_APP_PROXY_API_URL,
      logLevel: "debug",
      changeOrigin: true,
    })
  );
};
