<template>
  <div class="w-full h-full flex flex-col justify-center items-center p-4">

    <div class="w-full h-full flex flex-col justify-center items-center p-4">
      <q-toggle class="text-xs ml-auto mr-32" dense v-model="autoRefresh" :label="$t('autoRefresh')" />
    </div>

    <AuthorRanking class="q-pa-md w-[80vw] max-w-[1080px]" :authors="rows" :since="since" :until="until" />

    <!-- 提交频率图-->
    <div class="q-pa-md w-[80vw] max-w-[1080px]" v-for="(item, index) in commitsPeriods" :key="index">
      <PeriodCard :commits="item" :year="since" :title="$t('commitPeriod') + '-' + item.nickname" />
    </div>

  </div>
</template>

<script lang="ts" setup>
import AuthorRanking from "@/components/aynalisis/AuthorRanking.vue";
import PeriodCard from "@/components/aynalisis/PeriodCard.vue";

const api = useApi();

const i18n = useI18n();

const rows = ref<any[]>([]);

const since = ref("");
const until = ref("");

const commitsPeriods = ref<any[]>([]);

const getContributors = async () => {
  const resp: any = await api.getContributors();
  rows.value = resp?.data || [];
  since.value = resp?.meta?.since || '';
  until.value = resp?.meta?.until || '';

  // getCommitsPeriods();
};

const getCommitsPeriods = async () => {
  for (const item of rows.value) {
    api.getCommitPeriod({
      period: 'day',
      authors: [item.nickname],
    }).then((resp: any) => {
      if (resp?.data) {
        resp.data.nickname = item.nickname || item.authorName || "";
        commitsPeriods.value.push(resp.data || []);
      }
    });
  }
};

const refreshData = () => {
  getContributors();
}

const interval = ref()
const autoRefresh = ref(true)

onMounted(() => {
  refreshData()
  interval.value = setInterval(() => {
    if (autoRefresh.value) {
      refreshData()
    }
  }, 30000)
})

onUnmounted(() => {
  clearInterval(interval.value)
})

useSeoMeta({
  title: i18n.t("contributors"),
  description: i18n.t("contributors"),
});
</script>
