<template>
  <div class="w-full h-full flex justify-center items-center p-4">
    <AuthorRanking :authors="rows" :since="since" />
  </div>
</template>

<script lang="ts" setup>
import AuthorRanking from "@/components/aynalisis/AuthorRanking.vue";

const api = useApi();

const i18n = useI18n();

const rows = ref<any[]>([]);

const since = ref("");

const getContributors = async () => {
  const resp: any = await api.getContributors();
  rows.value = resp?.data || [];
  since.value = resp?.meta?.since;
};
onMounted(() => {
  getContributors();
});
useSeoMeta({
  title: i18n.t("contributors"),
  description: i18n.t("contributors"),
});
</script>
