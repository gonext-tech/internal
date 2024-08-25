module.exports = {
  apps: [
    {
      name: "internal", // Replace with your desired app name
      script: "./bin/internal", // Replace with the path to your Go application build
      watch: false, // Set to true if you want PM2 to watch for changes and restart
      env: {
        NODE_ENV: "production"
      },
      env_production: {
        NODE_ENV: "production"
      }
    }
  ]
};
