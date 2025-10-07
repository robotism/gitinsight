<template>
  <q-toolbar class="flex row w-full h-full items-center justify-between">
    <div class="flex row flex-nowrap items-center">
      <AppTitle class="ml-2" />
      <q-btn
        class="text-nowrap flex md:hidden font-bold"
        flat
        dense
        padding="8px"
        color="primary"
        icon="las la-bars"
        @click="isNaviDrawerOpen = !isNaviDrawerOpen"
        :label="$te(breadcrumbRoot) ? $t(breadcrumbRoot) : ''"
      />
    </div>

    <AppNavi class="w-auto items-center hidden md:flex" />

    <q-drawer overlay v-model="isNaviDrawerOpen" :breakpoint="99999">
      <q-card
        class="flex column item-center justify-between m-2 p-2"
        :style="{
          minHeight: $env.style.header.height,
        }"
        :bordered="false"
        flat
      >
        <AppTitle class="mt-8" />
        <AppNavi class="mt-8" vertical />
      </q-card>
    </q-drawer>

    <div class="flex row items-center">
      <q-btn
        dense
        flat
        padding="8px"
        color="primary"
        icon="bi-github"
        v-if="$env.site.github"
        @click="$win.open($env.site.github)"
      >
        <q-tooltip>
          {{ $env.site.github }}
        </q-tooltip>
      </q-btn>
      <LocaleSelectButton />
      <ThemeToggleButton />
    </div>
  </q-toolbar>
</template>

<script lang="ts" setup>
import AppTitle from "./AppTitle.vue";
import AppNavi from "./AppNavi.vue";
import LocaleSelectButton from "~/components/tool/LocaleSelectButton.vue";
import ThemeToggleButton from "~/components/tool/ThemeToggleButton.vue";

const isNaviDrawerOpen = ref(false);

const router = useRouter();
const localeRoute = useLocaleRoute();
const $win = useWindow();

const breadcrumbRoot = computed(() => {
  const cur = router.currentRoute.value.path;
  const loc = localeRoute("/")?.path || "";
  const path = cur.replace(loc, "");
  const items = path.split("/").filter((i) => !!i);
  return "nav." + items[0];
});
</script>
