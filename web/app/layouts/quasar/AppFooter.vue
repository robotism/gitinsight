<template>
  <div class="flex h-full w-full items-center justify-center">
    <span v-if="$env.site.host">
      Copyright © {{ year }} {{ $env.site.host }} All rights reserved.
    </span>
  </div>
</template>

<script lang="ts" setup>
const year = computed(() => new Date().getFullYear());
</script>
