<template>
  <div class="hidden dark:flex">
    <q-btn
      dense
      flat
      padding="8px"
      color="primary"
      icon="bi-moon"
      @click="isDark = !isDark"
    />
  </div>
  <div class="flex dark:hidden">
    <q-btn
      dense
      flat
      padding="8px"
      color="primary"
      icon="bi-sun"
      @click="isDark = !isDark"
    />
  </div>
</template>

<script lang="ts" setup>
import { Dark } from "quasar";

const colorMode = useColorMode();

const isDark = computed({
  get() {
    return colorMode.value === "dark";
  },
  set(value: boolean) {
    Dark.set(value);
    colorMode.preference = value ? "dark" : "light";
  },
});

onMounted(() => {
  Dark.set(isDark.value);
});
</script>
