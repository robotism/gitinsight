<template>
  <div class="w-full h-[calc(100vh-120px)]">
    <h6 class="ml-4 flex nowrap">📝 {{ $t('commitLogs') }}</h6>
    <q-infinite-scroll class="h-full overflow-auto" :offset="50" @load="onLoad">
      <q-item v-for="(c, index) in commitLogs" :key="index">
        <div class="flex flex-col">
          <div class="w-full text-[11px] text-blue-900 font-bold">
            🏷️ {{ c.commitHash }}
          </div>
          <div class="w-full text-black-600 text-[10px] break-all">
            {{ c.message }}
          </div>
          <div class="w-full flex justify-between text-xs text-gray-500 cursor-pointer" @click="gotoRepo(c.repoUrl)">
            <div class="text-[8px]">{{ c.repoUrl.split('/').pop().replace('.git', '') + '/' + c.branchName }}</div>
            <div class="text-green text-[10px]">+{{ c.additions }}</div>
            <div class="text-red text-[10px]">-{{ c.deletions }}</div>
            <div class="text-purple text-[10px]">+{{ c.effectives }}</div>
          </div>
          <div class="w-full flex justify-between">
            <q-badge class="text-xs" color="blue">{{ c.nickname }}</q-badge>
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
  </div>
</template>

<script setup lang="ts">
const api = useApi()
const fmt = useFormat()
const { open } = useWindow()

const offset = ref(0)
const limit = ref(20)
const loading = ref(false)
const commitLogs: Ref<any[]> = ref([])

const props = defineProps({
  filter: Object
})

const gotoRepo = (url: string) => open(url)

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

</script>