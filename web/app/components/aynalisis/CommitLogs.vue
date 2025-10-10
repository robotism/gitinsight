<template>
  <h6 class="w-full flex nowrap content-center">
    📝 {{ $t('commitLogs') }}
    <q-toggle class="text-xs ml-16" dense v-model="autoRefresh" :label="$t('autoRefresh')" />
  </h6>
  <q-infinite-scroll class="w-full h-[calc(100vh-150px)] overflow-auto" :offset="50" @load="onLoad">
    <q-item class="w-full" v-for="(c, index) in commitLogs" :key="index">
      <div class="w-full flex flex-col cursor-pointer" @click="gotoCommitUrl(c)">
        <div class="w-full text-[10px] text-blue-900 font-bold nowrap">
          🏷️ {{ c.commitHash }}
        </div>
        <div class="w-full text-black-600 text-[10px] break-all">
          {{ c.message }}
        </div>
        <div class="w-full flex justify-between text-xs text-gray-500 ">
          <div class="text-[8px]">{{ c.repoUrl.split('/').pop().replace('.git', '') + '/' + c.branchName }}</div>
          <div class="text-green text-[10px]">+{{ c.additions }}</div>
          <div class="text-red text-[10px]">-{{ c.deletions }}</div>
          <div class="text-purple text-[10px]">+{{ c.effectives }}</div>
        </div>
        <div class="w-full flex justify-between">
          <q-badge class="text-xs" :style="{ backgroundColor: hashColor(c.nickname), color: '#fff' }">{{ c.nickname }}</q-badge>
          <q-badge class="text-[10px]" color="secondary">⏰ {{ fmt.localDate(c.date) }}</q-badge>
        </div>
      </div>
    </q-item>

    <template #loading>
      <div class="row justify-center q-my-md">
        <q-spinner-dots color="primary" size="40px" />
      </div>
    </template>
  </q-infinite-scroll>
</template>

<script setup lang="ts">
const api = useApi()
const fmt = useFormat()
const { open } = useWindow()
const { hashColor } = useColor()

const offset = ref(0)
const limit = ref(20)
const loading = ref(false)

const meta: Ref<any> = ref({})
const commitLogs: Ref<any[]> = ref([])

const props = defineProps({
  filter: Object
})

const gotoCommitUrl = (c: any) => {
  const configs = meta.value.config || []
  const commitUrlTmpl = configs.find((item: any) => c.repoUrl.indexOf(item.domain) !== -1)?.commitUrlTmpl
  if (commitUrlTmpl) {
    const commitUrl = commitUrlTmpl.replace('{{.RepoUrl}}', c.repoUrl.replace('.git', '')).replace('{{.BranchName}}', c.branchName).replace('{{.CommitHash}}', c.commitHash)
    open(commitUrl)
  }else{
    const commitUrl = c.repoUrl.replace('.git', '') + '/commit/' + c.commitHash+"?branch="+c.branchName
    open(commitUrl)
  }
}

const onLoad = async (index?: number, done?: (stop?: boolean) => void) => {
  if (loading.value) return
  loading.value = true

  if (index !== undefined && index >= 1) {
    offset.value = (index - 1) * limit.value
  }
  if (offset.value === 0) {
    commitLogs.value = []
  }

  const res: any = await api.getCommitLogs(props.filter, offset.value, limit.value)
  if (res.code === 200) {
    commitLogs.value.push(...res.data)
    meta.value = res.meta
    done?.(res.data.length < limit.value) // true = 停止
  } else {
    done?.(true)
  }

  loading.value = false
}


watch(() => props.filter, (val) => {
  offset.value = 0
  onLoad()
}, { deep: true })


const refreshData = () => {
  offset.value = 0
  onLoad()
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