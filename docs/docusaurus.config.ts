import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";

const config: Config = {
  title: "Dispute Resolution Engine",
  tagline: "Streamlining the ADR Process using AI and Natural Language Processing",
  favicon: "img/dre_logo.svg",

  // Set the production url of your site here
  url: "https://capstone-dre.dns.net.za",
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: "/Dispute-Resolution-Engine/",

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: "cos301-se-2024", // Usually your GitHub org/user name.
  projectName: "Dispute-Resolution-Engine", // Usually your repo name.
  trailingSlash: false,

  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      {
        docs: {
          sidebarPath: "./sidebars.ts",
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            "https://github.com/COS301-SE-2024/Dispute-Resolution-Engine/tree/dev/dre-docs/",
        },
        theme: {
          customCss: "./src/css/custom.css",
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    image: "img/docusaurus-social-card.jpg",
    navbar: {
      title: "DRE",
      logo: {
        alt: "DRE Logo",
        src: "img/dre_logo.svg",
      },
      items: [
        {
          type: "docSidebar",
          sidebarId: "tutorialSidebar",
          position: "left",
          label: "Documentation",
        },
        // { to: "/blog", label: "Blog", position: "left" },
        {
          href: "https://github.com/COS301-SE-2024/Dispute-Resolution-Engine",
          label: "GitHub",
          position: "right",
        },
        {
          href: "https://capstone-dre.dns.net.za",
          label: "Go to Site",
          position: "left",
        },
      ],
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
