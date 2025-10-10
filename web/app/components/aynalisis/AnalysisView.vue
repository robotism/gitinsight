<template>

    <div class="w-full pl-4 pr-4 flex-col ">
        <h6 class="w-full content-center">
            {{ $t('analysisView') }}
            <q-toggle class="text-xs ml-16" dense v-model="autoRefresh" :label="$t('autoRefresh')" />
        </h6>
        <div class="w-full h-[calc(100vh-130px)] overflow-y-scroll flex flex-col items-center">
            <div class="w-full flex flex-col" :class="$q.screen.gt.sm ? 'max-w-[700px]' : 'w-full'">
                <!-- 贡献者 -->
                <AuthorRanking :class="$q.screen.gt.sm ? 'max-w-[700px]' : 'w-full'" :authors="authors" :since="since"
                    :until="until" sortBy="effectives" sortDirection="desc" :hideColumns="['name', 'email']" />

                <!-- 提交频率图-->
                <PeriodCard :commits="commitsPeriodDay" :year="since" :title="$t('commitPeriod')" />

                <!-- 提交频率图-->
                <PeriodCard :commits="commitsPeriodWeek" :year="since" :title="$t('commitPeriod')" />

                <!-- 提交频率图-->
                <PeriodCard :commits="commitsPeriodMonth" :year="since" :title="$t('commitPeriod')" />

                <!-- Commit 活跃图 -->
                <HeatMapCard :commits="commitsAll" :year="since" :title="$t('commitHeatmap')" />

                <!-- Fix 活跃图 -->
                <HeatMapCard :commits="commitsFix" :year="since" :title="$t('fixHeatmap')" />

                <!-- Feat 活跃图 -->
                <HeatMapCard :commits="commitsFeat" :year="since" :title="$t('featHeatmap')" />

                <!-- Merge 活跃图 -->
                <HeatMapCard :commits="commitsMerge" :year="since" :title="$t('mergeHeatmap')" />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import AuthorRanking from "@/components/aynalisis/AuthorRanking.vue";
import PeriodCard from "@/components/aynalisis/PeriodCard.vue";
import HeatMapCard from "~/components/aynalisis/HeatMapCard.vue";

const api = useApi()

const authors = ref<any[]>([]);
const since = ref("");
const until = ref("");
const commitsAll = ref<any[]>([]);
const commitsFix = ref<any[]>([]);
const commitsFeat = ref<any[]>([]);
const commitsMerge = ref<any[]>([]);
const commitsPeriodDay = ref<any[]>([]);
const commitsPeriodWeek = ref<any[]>([]);
const commitsPeriodMonth = ref<any[]>([]);

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
        until.value = resp?.meta?.until || until.value;
    });
    api.getCommitHeatmap({
        ...props.filter,
        messageType: '',
    }).then((resp: any) => {
        commitsAll.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
        until.value = resp?.meta?.until || until.value;
    });
    api.getCommitHeatmap({
        ...props.filter,
        messageType: 'fix',
    }).then((resp: any) => {
        commitsFix.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
        until.value = resp?.meta?.until || until.value;
    });
    api.getCommitHeatmap({
        ...props.filter,
        messageType: 'feat',
    }).then((resp: any) => {
        commitsFeat.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
        until.value = resp?.meta?.until || until.value;
    });
    api.getCommitHeatmap({
        ...props.filter,
        messageType: '',
        isMerge: '1',
    }).then((resp: any) => {
        commitsMerge.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
        until.value = resp?.meta?.until || until.value;
    });
    api.getCommitPeriod({
        ...props.filter,
        period: 'day',
    }).then((resp: any) => {
        commitsPeriodDay.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
        until.value = resp?.meta?.until || until.value;
    });
    api.getCommitPeriod({
        ...props.filter,
        period: 'week',
    }).then((resp: any) => {
        commitsPeriodWeek.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
        until.value = resp?.meta?.until || until.value;
    });
    api.getCommitPeriod({
        ...props.filter,
        period: 'month',
    }).then((resp: any) => {
        commitsPeriodMonth.value = resp?.data || [];
        since.value = resp?.meta?.since || since.value;
        until.value = resp?.meta?.until || until.value;
    });
};

const debounce = ref()
const refreshData = () => {
    clearTimeout(debounce.value)
    debounce.value = setTimeout(() => {
        getData()
    }, 300)
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


</script>