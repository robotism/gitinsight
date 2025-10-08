<template>
  <VChart class="w-full h-[400px] border border-gray-200 m-1" :option="option" autoresize />
</template>

<script setup lang="ts">
const i18n = useI18n()

// ✅ 接收父组件传入数据
const props = defineProps({
  title: {
    type: String,
    default: 'Commit Trend'
  },
  commits: {
    type: Array,
    required: true,
    // 格式: [{ period: '2025-01-02', nickname: 'robotism', effectives: 957 }]
  }
})

/**
 * 生成折线图配置
 */
const option = computed(() => {
  const data = props.commits || []

  // 1️⃣ 获取所有日期
  const datesSet = new Set<string>()
  data.forEach((d: any) => datesSet.add(d.period))
  const dates = Array.from(datesSet).sort()

  // 2️⃣ 获取所有昵称
  const nicknamesSet = new Set<string>()
  data.forEach((d: any) => nicknamesSet.add(d.nickname))
  const nicknames = Array.from(nicknamesSet)

  // 3️⃣ 构建每个人的 effectives 数组
  const series = nicknames.map((name) => {
    const seriesData = dates.map(date => {
      const item: any = data.find((d: any) => d.nickname === name && d.period === date)
      return item ? item.effectives : 0
    })
    return {
      name,
      type: 'line',
      smooth: true,
      data: seriesData,
      symbol: 'circle',
      symbolSize: 4
    }
  })

  // 4️⃣ 自动颜色池
  const colors = [
    '#409EFF', '#67C23A', '#F56C6C', '#E6A23C', '#909399', '#5c9ded', '#f0ad4e',
    '#d46b08', '#2f4554', '#61a0a8', '#c23531', '#2f4554', '#91c7ae', '#749f83'
  ]
  series.forEach((s: any, idx) => s.lineStyle = { color: colors[idx % colors.length] })

  return {
    title: {
      text: props.title || i18n.t('commitTrend'),
      left: 'center',
      top: 10,
      textStyle: { fontSize: 15, fontWeight: 'bold' }
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#fff',
      borderColor: '#ccc',
      borderWidth: 1,
      textStyle: { color: '#333' },
      formatter: (params: any[]) => {
        let html = `<b>${params[0].axisValue}</b><br/>`
        params.forEach(p => {
          html += `<span style="color:${p.color}">●</span> ${p.seriesName}: <b>${p.data}</b><br/>`
        })
        return html
      }
    },
    legend: {
      top: 35,
      left: 'center',
      type: 'scroll',
      data: nicknames
    },
    grid: {
      left: 50,
      right: 20,
      top: 70,
      bottom: 50
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: {
        fontSize: 10,
        rotate: 30,
        color: '#666'
      },
      boundaryGap: false,
      axisPointer: { type: 'shadow' }
    },
    yAxis: {
      type: 'value',
      axisLabel: { fontSize: 10, color: '#666' },
      splitLine: { lineStyle: { color: '#eee' } }
    },
    series,
    dataZoom: [
      {
        type: 'slider',
        xAxisIndex: 0, // 指定作用在 xAxis[0]
        start: 0,
        end: Math.min(30 / dates.length * 100, 100), // 默认显示最近 30 天或全部
        handleSize: 10
      },
      {
        type: 'inside',
        xAxisIndex: 0,
        start: 0,
        end: Math.min(30 / dates.length * 100, 100)
      }
    ]
  }
})

</script>

<style scoped>
div[ref="chartRef"] {
  min-width: 300px;
  min-height: 240px;
}
</style>
