<template>
  <q-tabs
    class="mx-1"
    style="color: var(--ui-primary)"
    dense
    shrink
    no-caps
    stretch
    :model-value="tab"
    :align="align"
    :vertical="vertical"
    :inline-label="inlineLabel==true || inlineLabel==false? inlineLabel : !$q.screen.lt.md"
    :switch-indicator="switchIndicator"
  >
    <q-tab
      v-for="(item, index) in links"
      :key="index"
      :name="item.label"
      :icon="item.icon"
      :label="$t(item.label)"
      style="justify-content: start;"
      @click="item.click"
    />
  </q-tabs>
</template>

<script lang="ts" setup>
const { getPagePathArray } = useLocalePage();
const { links } = useNavigation();

const props = defineProps({
  vertical: {
    type: Boolean,
    default: undefined,
  },
  inlineLabel: {
    type: Boolean,
    default: undefined,
  },
  switchIndicator: {
    type: Boolean,
    default: undefined,
  },
  align: {
    type: String as PropType<
      "left" | "center" | "right" | "justify" | undefined
    >,
    default: "center",
  },
});

const tab = computed(() => {
  const page = getPagePathArray()?.[0] || "";
  return (
    links.value.find((it) => {
      return it.label.indexOf(page) >= 0;
    })?.label ||
    links.value[0]?.label ||
    ""
  );
});
</script>

<style scoped>
.light-mode {
  .q-tab--active {
    background-color: rgba(127, 127, 127, 0.25);
  }
}
.dark-mode {
  .q-tab--active {
    background-color: rgba(255, 255, 255, 0.25);
  }
}

:deep(.q-tab--active) {
  .q-tab__label {
    font-weight: 600 !important;
  }
}
</style>
