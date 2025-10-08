<template>
  <VChart class="w-full h-[400px] border border-gray-200" :option="option" autoresize />
</template>

<script setup>
import dayjs from 'dayjs'

const i18n = useI18n()

// ✅ 接收父组件传入数据
const props = defineProps({
  title: {
    type: String,
  },
  commits: {
    type: Array,
    required: true,
    // 格式: [{ date: '2025-01-02', commits: 3 }, { date: '2025-01-03', commits: 12 }]
  }
})


/**
 * 自动计算数据范围（适配跨年情况，并固定显示至少 showDays 天）
 */
const calcRange = (data, showDays = 90) => {
  const today = new Date()
  const endDate = dayjs(today)
  const startDate = endDate.subtract(showDays - 1, 'day')

  const datesFromData = (data || []).map(d => d[0]).sort()
  const firstDataDate = datesFromData[0] ? dayjs(datesFromData[0]) : startDate

  // 开始日期取数据最早日期或 showDays 天前，两者取最早
  const start = firstDataDate.isBefore(startDate) ? firstDataDate : startDate
  const end = endDate

  return [start.format('YYYY-MM-DD'), end.format('YYYY-MM-DD')]
}

/**
 * 生成图表配置
 */
const option = computed(() => {
  const data = (props.commits || [])
    .filter(item => !!item.date)
    .map(item => [item.date.slice(0, 10), item.commits || 0])
    .sort((a, b) => a[0].localeCompare(b[0]))

  const range = calcRange(data)
  const maxVal = Math.max(...data.map(d => d[1]), 1)

  return {
    title: {
      text: props.title || i18n.t('commitHeatmap'),
      left: 'center',
      top: 10,
      textStyle: { fontSize: 15, fontWeight: 'bold' }
    },
    tooltip: {
      position: 'top',
      formatter: (p) => {
        const [date, value] = p.data
        return `${date}<br/><b>${value}</b> ${i18n.t('commits')}`
      }
    },
    visualMap: {
      type: 'piecewise',
      show: true,
      min: 0,
      max: maxVal,
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      top: 20,
      inRange: {
        color: ['#e6f7ff', '#91d5ff', '#1890ff', '#003a8c'] // 蓝色系渐变
      },
    },
    calendar: {
      top: 80,
      left: 20,
      right: 20,
      bottom: 10,
      cellSize: ['auto', 14],
      range,
      itemStyle: {
        borderWidth: 0.4,
        borderColor: '#ccc'
      },
      yearLabel: { show: false },
      monthLabel: { nameMap: 'cn', color: '#555', fontSize: 10 },
      dayLabel: { firstDay: 1, nameMap: 'cn', color: '#888', fontSize: 9 }
    },
    series: [
      {
        type: 'heatmap',
        coordinateSystem: 'calendar',
        data,
        emphasis: {
          focus: 'self',
          itemStyle: {
            borderColor: '#333',
            borderWidth: 1
          }
        }
      }
    ]
  }
})


</script>

<style scoped>
/* 让容器在小屏时自适应 */
div[ref="chartRef"] {
  min-width: 300px;
  min-height: 240px;
}
</style>