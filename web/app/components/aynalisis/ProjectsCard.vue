<template>
    <VChart :option="option" autoresize />
</template>

<script setup>

const i18n = useI18n()

const { hashColor } = useColor()

const props = defineProps({
    title: {
        type: String,
    },
    projects: {
        type: Array,
        required: true
    },
    metric: {
        type: String,
        default: 'effectives' // 可选: additions, deletions, effectives, commits
    }
})


// 解析 repo 名称
const branches = computed(() =>
    (props.projects || []).map((p) => ({
        repoBranch: ((p.repoUrl || '')?.split('/').pop().replace('.git', '') || p.repoUrl) + '/' + p.branchName,
        ...p
    }))
)

// 排序逻辑：按 metric 排序（默认按有效代码量）
const sortedBranches = computed(() => {
    const arr = [...branches.value]
    const key = props.metric
    return arr.sort((a, b) => (b[key] || 0) - (a[key] || 0))
})

// 图表配置
const option = computed(() => {
    const metricKey = props.metric
    const metricLabel = i18n.t(metricKey)

    return {
        title: {
            text: `${props.title}`,
            left: 'center',
            textStyle: { fontSize: 16, fontWeight: 'bold' }
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: { type: 'shadow' }
        },
        grid: { left: 100, right: 40, top: 60, bottom: 30 },
        xAxis: {
            type: 'value',
            axisLabel: { formatter: '{value}' }
        },
        yAxis: {
            type: 'category',
            inverse: true,
            data: sortedBranches.value.map((p) => p.repoBranch).map((p) => {
                try {
                    const cherrypick = '/cherry-pick-'
                    if (p.indexOf(cherrypick) > -1) {
                        const index = p.indexOf(cherrypick) + cherrypick.length
                        return p.substring(0, index) + p.substring(index, index + 6)
                    }
                } catch (e) {
                    console.error(e)
                }
                return p
            }),
            axisLabel: { fontSize: 8 }
        },
        series: [
            {
                name: metricLabel,
                type: 'bar',
                data: sortedBranches.value.map((p) => p[metricKey]),
                label: {
                    show: true,
                    position: 'right',
                    formatter: '{c}'
                }, itemStyle: {
                    color: (params) => {
                        return hashColor(params.name || params.data.repoUrl + '/' + params.data.branchName)
                    },
                },
                barMaxWidth: 25
            }
        ]
    }
})

</script>

<style scoped>
/* 可选：让容器更自适应 */
</style>