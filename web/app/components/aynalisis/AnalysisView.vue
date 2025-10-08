<template>
    <div class="overflow-y-scroll h-[calc(100vh-120px)]">
        <h6 class="ml-2">{{ $t('analysisView') }}</h6>

        <!-- 贡献者 -->
        <AuthorRanking :authors="authors" :since="since" sortBy="effectives" sortDirection="desc" />

        <!-- 提交频率图-->
        <PeriodCard :commits="commitsPeriod" :year="since" :title="$t('commitPeriod')" />

        <!-- Commit 活跃图 -->
        <HeatMapCard :commits="commitsAll" :year="since" :title="$t('commitHeatmap')" />

        <!-- Fix 活跃图 -->
        <HeatMapCard :commits="commitsFix" :year="since" :title="$t('fixHeatmap')" />

        <!-- Feat 活跃图 -->
        <HeatMapCard :commits="commitsFeat" :year="since" :title="$t('featHeatmap')" />

        <!-- Merge 活跃图 -->
        <HeatMapCard :commits="commitsMerge" :year="since" :title="$t('mergeHeatmap')" />
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import AuthorRanking from "@/components/aynalisis/AuthorRanking.vue";
import PeriodCard from "@/components/aynalisis/PeriodCard.vue";
import HeatMapCard from "~/components/aynalisis/HeatMapCard.vue";

const api = useApi()

const i18n = useI18n()

const authors = ref<any[]>([]);
const since = ref("");
const commitsAll = ref<any[]>([]);
const commitsFix = ref<any[]>([]);
const commitsFeat = ref<any[]>([]);
const commitsMerge = ref<any[]>([]);
const commitsPeriod = ref<any[]>([]);


const props = defineProps({
    filter: {
        type: Object,
        required: true
    }
})

watch(() => props.filter, () => {
    getData();
}, { deep: true })

const getData = async () => {
    api.getRanking(props.filter).then((resp: any) => {
        authors.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
    });
    api.getCommitHeatmap({
        ...props.filter,
        messageType: '',
    }).then((resp: any) => {
        commitsAll.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
    });
    api.getCommitHeatmap({
        ...props.filter,
        messageType: 'fix',
    }).then((resp: any) => {
        commitsFix.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
    });
    api.getCommitHeatmap({
        ...props.filter,
        messageType: 'feat',
    }).then((resp: any) => {
        commitsFeat.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
    });
    api.getCommitHeatmap({
        ...props.filter,
        messageType: '',
        isMerge: '1',
    }).then((resp: any) => {
        commitsMerge.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
    });
    api.getCommitPeriod({
        ...props.filter,
        period: 'day',
    }).then((resp: any) => {
        console.log(resp)
        commitsPeriod.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
    });
};
onMounted(() => {
    getData();
});

</script>