export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.config.globalProperties.$env = useRuntimeConfig().public;
});
