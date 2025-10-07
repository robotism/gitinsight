<template>
  <div class="flex flex-row items-center" v-if="locales.length > 1">
    <q-btn
      v-if="locales.length == 2"
      color="primary"
      dense
      flat
      padding="8px"
      class="flex row text-nowrap"
      icon="bi-translate"
      @click="toggleLocalePath"
    />
    <q-btn-dropdown
      v-else
      v-ripple.early
      dense
      flat
      padding="8px"
      color="primary"
      icon="bi-translate"
      :label="name"
    >
      <q-list>
        <div v-for="(lang, index) in langs" :key="index">
          <q-item clickable v-close-popup @click="lang.onSelect()">
            <q-item-section>
              <q-item-label>{{ lang.label }}</q-item-label>
            </q-item-section>
          </q-item>
        </div>
      </q-list>
    </q-btn-dropdown>
  </div>
</template>

<script lang="ts" setup>
const { locale, locales, messages } = useI18n();
const router = useRouter();
const switchLocalePath = useSwitchLocalePath();

import { Lang, type QuasarLanguage } from "quasar";

const name = computed(() => {
  return locales.value.find((i) => i.code === locale.value)?.name || "";
});

const langs = computed(() => {
  return locales.value.map((i) => {
    return {
      code: i.code,
      label: i.name || i.code,
      color: locale.value === i.code ? "primary" : "neutral",
      onSelect: () => {
        router.push(switchLocalePath(i.code));
      },
      click: () => {
        router.push(switchLocalePath(i.code));
      },
    };
  });
});

const toggleLocaleCode = computed(() => {
  const i = locales.value.find((i) => i.code !== locale.value);
  return i?.code || locale.value;
});

const toggleLocalePath = () => {
  const toggleCode = toggleLocaleCode.value;
  router.push(switchLocalePath(toggleCode));
};

const updateQuasarLang = () => {
  const message = messages.value[locale.value];
  Lang.set(message?.quasar as QuasarLanguage);
};

watch(
  () => locale.value,
  () => {
    updateQuasarLang();
  }
);

onMounted(() => {
  updateQuasarLang();
});
</script>
