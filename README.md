## BatNoter Frontend Module
This is the frontend module of batnoter application.

### Build the frontend module

The following commands are used to build this frontend module.
```shell
npm install
npm run build
```

### Running the frontend module locally

Before starting the frontend app, please make sure that you have either started the backend module locally or pointed the frontend to staging api-server using proxy configuration.

If you are working only on the frontend changes and do not want to start the backend module locally, you can point the frontend app to staging server for API access.

To point frontend app to the staging api-server, Simply create `.env.development.local` file at the root of the frontend module with below contents.
```shell
REACT_APP_PROXY_API_URL=https://batnoter-staging.herokuapp.com
```

Then start the frontend react app with below command.
```shell
npm start
```
This will start the react app in the development mode.

Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

### Analyzing the bundle size
Source map explorer analyzes JavaScript bundles using the source maps. This helps you understand where code bloat is coming from.

Run below command to generate the report.
```shell
npm run analyze
```

### Create production build
```shell
npm run build
```

This will create the production build of application in the `build` folder.
It correctly bundles application in production mode and optimizes the build for the best performance.

### Run the production build of frontend module
When you do `npm start` the app is started in development mode.

If you have built the app and you want to run the app in production mode (i.e.from build directory), 
Install `serve` and start the app using below commands
```shell
npm install serve -g
npm run build
serve -s build
```

Open [http://localhost:3000](http://localhost:3000) to view it in the browser.
