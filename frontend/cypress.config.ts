import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
    baseUrl: 'http://capstone-dre.dns.net.za',
  },
  env: {
    TEST_USER: "gaced90493@rinseart.com",
    TEST_PASSWORD: "Password1234#"
  },
  component: {
    devServer: {
      framework: "next",
      bundler: "webpack",
    },
  },
});
