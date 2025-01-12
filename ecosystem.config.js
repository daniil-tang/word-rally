module.exports = {
  apps: [
    {
      name: "word-rally-client",
      script: "./client/build/index.js",
      env_production: {
        NODE_ENV: "production",
        // BASE_URL: "http://localhost:8080",
      },
      env_development: {
        NODE_ENV: "development",
        // BASE_URL: "http://localhost:8080",
      },
    },
    // {
    //   name: "word-rally-server",
    //   script: "./server/word-rally-server",
    // },
  ],
};
