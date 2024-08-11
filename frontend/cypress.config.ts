import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
    baseUrl: 'http://capstone-dre.dns.net.za',
  },

  component: {
    devServer: {
      framework: "next",
      bundler: "webpack",
    },
  },
});
