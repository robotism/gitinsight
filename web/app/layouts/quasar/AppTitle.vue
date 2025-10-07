<template>
  <div class="flex row flex-nowrap items-center">
    <img
      v-if="$env.site.logo"
      :src="$env.site.logo"
      :alt="$env.site.logo"
      class="hidden md:flex h-[24px] w-[24px] mr-2 rounded-full outline-2 outline-offset-1 outline-cyan-500"
    />
    <NuxtLinkLocale
      to="/"
      class="cursor-pointer text-primary hover:text-secondary text-xl font-bold"
      style="text-shadow: #fc0 1px 0 10px"
    >
      {{ $t("app.name") }}
    </NuxtLinkLocale>
  </div>
</template>

<script lang="ts" setup></script>
