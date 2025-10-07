<template>
  <Suspense>
    <template #default>
      <div class="">
        <NuxtRouteAnnouncer />
        <NuxtLoadingIndicator />
        <NuxtLayout name="quasar" class="">
          <NuxtPage />
        </NuxtLayout>
      </div>
    </template>
    <template #fallback>
      <div class="h-screen w-screen flex justify-center items-center">
        <q-spinner-bars color="primary" size="4em" />
      </div>
    </template>
  </Suspense>
</template>

<script lang="ts" setup>
const router = useRouter();
const { computedTitle } = useWindow();

const { start, finish } = useLoadingIndicator({
  duration: 300,
  throttle: 0,
});

const beforeRoute = (to: any, from: any, next: any) => {
  start();
  next();
};

const afterRoute = (to: any, from: any) => {
  finish();
};

onMounted(() => {
  router.beforeEach(beforeRoute);
  router.afterEach(afterRoute);
});

onUnmounted(() => {
  router.beforeEach(() => {});
  router.afterEach(() => {});
});

useHead(() => ({
  titleTemplate: computedTitle(),
}));
</script>
