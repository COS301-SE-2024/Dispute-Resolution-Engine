import type { Config } from "tailwindcss";

const config = {
  content: [
    "./pages/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./app/**/*.{ts,tsx}",
    "./src/**/*.{ts,tsx}",
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
    },
  },
  plugins: [require("tailwindcss-animate")],
} satisfies Config;

export default config;
