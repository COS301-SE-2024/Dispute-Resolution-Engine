import type { Config } from "tailwindcss";

const config = {
  //darkMode: ["class"],
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  prefix: "",
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {
      keyframes: {
        "accordion-down": {
          from: { height: "0" },
          to: { height: "var(--radix-accordion-content-height)" },
        },
        "accordion-up": {
          from: { height: "var(--radix-accordion-content-height)" },
          to: { height: "0" },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
      },
      colors: {
        primary: {
          50: "#DEE9FF",
          100: "#C9DCFF",
          200: "#A0C2FF",
          300: "#78A8FF",
          400: "#4F8DFF",
          500: "#2673FF",
          600: "#0054ED",
          700: "#0040B5",
          800: "#002C7D",
          900: "#001845",
          950: "#000E29",
        },
        secondary: {
          50: "#FFFFFF",
          100: "#FFFFFF",
          200: "#FFFFFF",
          300: "#FFFFFF",
          400: "#DEFCFF",
          500: "#B5F8FF",
          600: "#7DF3FF",
          700: "#45EDFF",
          800: "#0DE8FF",
          900: "#00C0D4",
          950: "#00A6B8",
        },
        "surface-light": {
          50: "#FFFFFF",
          100: "#F0F5F6",
          200: "#D0DEE3",
          300: "#B0C5D0",
          400: "#90A9BD",
          500: "#708BAA",
          600: "#576C90",
          700: "#435071",
          800: "#303751",
          900: "#1D1F31",
          950: "#141421",
        },
        "surface-dark": {
          50: "#7EA4DE",
          100: "#729ADB",
          200: "#5A86D4",
          300: "#4170CE",
          400: "#315CC0",
          500: "#2B4DA7",
          600: "#253F8F",
          700: "#1F3277",
          800: "#18255E",
          900: "#161D42",
          950: "#0F132D",
        },

        dre: {
          bg: {
            light: "#f0f5f6",
            dark: "#111827",
          },
          "100": "#9beaf2",
          "200": "#336cd4",
          "300": "#24479d",
          "400": "#080965",
          "500": "#040653",
        },
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
} satisfies Config;

export default config;
