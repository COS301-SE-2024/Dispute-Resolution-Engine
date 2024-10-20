import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
    baseUrl: 'https://capstone-dre.dns.net.za',
  },
  env: {
    TEST_USER: "johndoe@example.com",
    TEST_PASSWORD: "qwerty12345!"
  },
  component: {
    devServer: {
      framework: "next",
      bundler: "webpack",
    },
  },
});
