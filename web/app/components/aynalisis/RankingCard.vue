<template>
    <div ref="chartRef" class="w-[360px] h-[270px]"></div>
</template>

<script setup lang="ts">
import * as echarts from "echarts";

const i18n = useI18n();
const { size } = useWindow();
const { rankColor } = useColor();

const props = defineProps({
    title: {
        type: String,
    },
    metric: {
        type: String,
        default: "effectives", // 可选: additions, deletions, effectives, commits
    },
    ranking: {
        type: Array<any>,
        default: [],
    },
});

const chartRef = ref<HTMLDivElement | null>(null);
const chart = ref<echarts.ECharts | null>(null);

const metricKey = computed(() => props.metric);
const metricLabel = computed(() => i18n.t(metricKey.value));

// 排序后的前十名（或全部）
const sortedRanking = computed(() => {
    const arr = [...props.ranking || []];
    const key = metricKey.value;
    return arr.sort((a, b) => (b[key] || 0) - (a[key] || 0)).slice(0, 10);
});

// ECharts 配置
const option = computed(() => ({
    title: {
        text: props.title,
        left: "center",
        textStyle: { fontSize: 18, fontWeight: "bold" },
    },
    tooltip: {
        trigger: "axis",
        axisPointer: { type: "shadow" },
        formatter: (params: any) => {
            const item = params[0];
            const c = sortedRanking.value[item.dataIndex] || {};
            return `
          <div>
            <strong>${c.nickname}</strong><br/>
            ${i18n.t("commits")}: ${c.commits}<br/>
            ${i18n.t("additions")}: <span style="color:green">${c.additions}</span><br/>
            ${i18n.t("deletions")}: <span style="color:red">${c.deletions}</span><br/>
            ${i18n.t("effectives")}: <span style="color:teal">${c.effectives}</span>
          </div>`;
        },
    },
    grid: { left: 100, right: 40, top: 80, bottom: 40 },
    xAxis: {
        type: "value",
        axisLabel: { fontSize: 13 },
    },
    yAxis: {
        type: "category",
        inverse: true,
        axisLabel: { fontSize: 10 },
        data: sortedRanking.value.map((c) => c.nickname || ""),
    },
    series: [
        {
            name: metricLabel.value,
            type: "bar",
            barMaxWidth: 25,
            data: sortedRanking.value.map((c) => c[metricKey.value] || 0),
            label: {
                show: true,
                position: "right",
                formatter: "{c}",
            },
            itemStyle: {
                color: (params: any) => {
                    return rankColor(params.dataIndex, sortedRanking.value.length);
                },
            },
        },
    ],
}));

// 自适应 resize
watch(() => size.value, () => {
    if (chart.value && !chart.value.isDisposed()) {
        chart.value.resize();
    }
}, { deep: true });

// 监听数据变化自动更新图表
watch(option, (val) => {
    if (chart.value) {
        chart.value.setOption(val, true); // 不合并旧配置，强制全量更新
    }
});

onMounted(async () => {
    chart.value = echarts.init(chartRef.value!);
    chart.value.setOption(option.value);
});

onUnmounted(() => {
    chart.value?.dispose();
});
</script>

<style scoped></style>