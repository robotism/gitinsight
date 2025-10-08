// https://nuxt.com/docs/api/configuration/nuxt-config

import appconfig from "./site.config";

import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  ssr: false,
  srcDir: "app/",
  routeRules: {
    "/": { prerender: true },
  },
  devtools: {
    enabled: true,
    timeline: {
      enabled: true,
    },
  },
  vite: {
    plugins: [tailwindcss()],
  },
  nitro: {
    output: {
      publicDir: "dist",
    },
    preset: "static",
    prerender: {
      crawlLinks: true,
      failOnError: false,
    },
  },
  runtimeConfig: {
    public: {
      ...appconfig,
    },
  },
  modules: [
    "@nuxt/devtools",
    "@nuxtjs/sitemap",
    "@nuxtjs/robots",
    "@nuxtjs/google-fonts",
    "@nuxtjs/i18n",
    "@nuxt/eslint",
    "@nuxt/content",
    "@nuxt/image",
    "@vueuse/nuxt",
    "@nuxt/ui",
    "nuxt-quasar-ui",
    "@nuxtjs/color-mode",
    "@nuxt/icon",
    "@nuxt/fonts",
    "nuxt-echarts",
  ],
  experimental: {
    payloadExtraction: false,
  },
  site: {
    url: appconfig.site.url,
  },
  colorMode: {
    preference: "system", // default value of $colorMode.preference
    fallback: "dark", // fallback value if not system preference found
    classPrefix: "",
    classSuffix: "-mode",
    storage: "localStorage", // or 'sessionStorage' or 'cookie'
    storageKey: "nuxt-color-mode",
  },
  eslint: {
    config: {
      stylistic: {
        indent: "tab",
        semi: true,
      },
    },
  },
  googleFonts: {
    families: {
      Roboto: true,
      Inter: [400, 700],
      "Josefin+Sans": true,
      Lato: [100, 300],
      Raleway: {
        wght: [100, 400],
        ital: [100],
      },
      "Crimson Pro": {
        wght: "200..900",
        ital: "200..700",
      },
    },
  },
  i18n: {
    baseUrl: appconfig.site.url,
    strategy: "prefix_and_default",
    locales: [
      {
        code: "en",
        name: "English",
        files: ["en-US.ts"],
      },
      {
        code: "zh",
        name: "中文",
        files: ["zh-CN.ts"],
      },
    ],
    detectBrowserLanguage: {
      useCookie: true,
      cookieCrossOrigin: true,
      redirectOn: "root", // all, root, no prefix
    },
    defaultLocale: `en`,
  },
  css: ["~/assets/css/app.css"],
  quasar: {
    lang: "en-US",
    plugins: [
      "Dark",
      "Notify",
      "Dialog",
      "Loading",
      "LoadingBar",
      "BottomSheet",
      "AppVisibility",
      "AppFullscreen",
      "LocalStorage",
      "SessionStorage",
    ],
    extras: {
      animations: "all",
      font: "roboto-font",
      fontIcons: ["bootstrap-icons", "line-awesome"],
    },
  },
  content: {
    watch: {
      enabled: true,
      port: 4000,
      showURL: false,
    },
    build: {
      markdown: {
        highlight: {
          theme: {
            // Default theme (same as single string)
            default: "github-light",
            // Theme used if `html.dark`
            dark: "github-dark",
            // Theme used if `html.sepia`
            sepia: "monokai",
          },
        },
      },
    },
  },
  echarts: {
    features: ["LabelLayout", "UniversalTransition"],
    charts: ["BarChart", "LineChart", "PieChart", "HeatmapChart"],
    components: [
      "DatasetComponent",
      "GridComponent",
      "TooltipComponent",
      "TitleComponent",
      "LegendComponent",
      'DataZoomComponent',
      'VisualMapComponent',
      'CalendarComponent',
    ],
  },
});
