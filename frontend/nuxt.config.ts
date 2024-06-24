// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: false },
  css: ["~/assets/css/main.css"],
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  modules: ["@vueuse/nuxt", "@pinia/nuxt", "@pinia-plugin-persistedstate/nuxt", "radix-vue/nuxt"],
  vite: {
    optimizeDeps: {
      exclude: ["pdfjs-dist/build/pdf"],
    },
  },
});
