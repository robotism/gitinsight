<template>
  <q-breadcrumbs class="mx-1 md:mx-4 my-3">
    <template v-slot:separator>
      <q-icon size="1.5em" name="chevron_right" color="primary" />
    </template>
    <q-breadcrumbs-el
      v-for="(item, index) in menus"
      :key="index"
      :label="$te(item.label) ? $t(item.label) : item.label"
      :icon="item.icon"
      :class="{
        'bg-gray-100 dark:bg-gray-800': index == menus.length - 1,
        'hover:bg-gray-200 dark:hover:bg-gray-700': index == menus.length - 1,
      }"
      class="cursor-pointer px-2 py-1 rounded"
      @click="$router.push(item.to)"
    />
  </q-breadcrumbs>
</template>

<script lang="ts" setup>
const { getPagePathArray } = useLocalePage();
const { links } = useNavigation();
const localeRoute = useLocaleRoute();

const menus = computed(() => {
  const array = getPagePathArray();
  const items: any[] = array.map((i: any) => {
    const index = array.indexOf(i);
    const navi = links.value.find((it) => {
      return it.label.indexOf(i) >= 0;
    });
    return {
      to: localeRoute("/" + array.slice(0, index + 1).join("/")),
      label: index == 0 ? navi?.label || i : i,
      icon: navi?.icon,
    };
  });
  if (items[0].label != links.value?.[0]?.label) {
    items.unshift({
      to: localeRoute("/"),
      label: links.value?.[0]?.label,
      icon: links.value?.[0]?.icon,
    });
  }
  return items;
});
</script>
