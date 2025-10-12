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
  commitsAll: {
    type: Array,
    required: true,
    // 格式: [{ date: '2025-01-02', commits: 3 }, { date: '2025-01-03', commits: 12 }]
  },
  commits: {
    type: Array,
    required: true,
    // 格式: [{ date: '2025-01-02', commits: 3 }, { date: '2025-01-03', commits: 12 }]
  }
})


/**
 * 自动计算数据范围（智能扩展至 showDays，且居中 data 范围）
 * @param {Array} data - 数据数组，每项第一列是日期字符串
 * @param {number} showDays - 最少显示的天数（默认90）
 * @returns {[string, string]} [startDate, endDate]
 */
const calcRange = ( showDays = 90) => {
  const data = (props.commitsAll || props.commits || [])
    .filter(item => !!item.date)
    .map(item => [item.date.slice(0, 10), item.commits || 0])
    .sort((a, b) => a[0].localeCompare(b[0]))

  const today = dayjs()
  if (!data || data.length === 0) {
    // 没有数据，默认显示最近 showDays 天
    const end = today
    const start = end.subtract(showDays - 1, 'day')
    return [start.format('YYYY-MM-DD'), end.format('YYYY-MM-DD')]
  }

  // 从 data 中提取所有日期，并排序
  const dates = data.map(d => dayjs(d[0])).sort((a, b) => a - b)
  const first = dates[0]
  const last = dates[dates.length - 1]
  const diff = last.diff(first, 'day') + 1

  if (diff >= showDays) {
    // 如果数据跨度 >= showDays，直接返回数据范围
    return [first.format('YYYY-MM-DD'), last.format('YYYY-MM-DD')]
  }

  // 否则，居中扩展到 showDays
  const extraDays = showDays - diff
  const before = Math.floor(extraDays / 2)
  const after = extraDays - before

  let start = first.subtract(before, 'day')
  let end = last.add(after, 'day')

  // 限制不要超过今天（end 最大不能超过 today）
  if (end.isAfter(today)) {
    const offset = end.diff(today, 'day')
    end = today
    start = start.subtract(offset, 'day')
  }

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

  const range = calcRange()
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