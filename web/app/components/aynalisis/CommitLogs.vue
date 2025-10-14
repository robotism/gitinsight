<template>
  <h6 class="w-full flex nowrap content-center">
    📝 {{ $t('commitLogs') }}
    <q-toggle class="text-xs ml-16" dense v-model="autoRefresh" :label="$t('autoRefresh')" />
  </h6>
  
  <!-- 添加key属性，filter变更时强制重新创建组件 -->
  <q-infinite-scroll 
    :key="filterKey"
    ref="infiniteScroll" 
    class="w-full h-[calc(100vh-150px)] overflow-auto" 
    :offset="100" 
    @load="onLoad"
  >
    <q-item class="w-full" v-for="(c, index) in commitLogs" :key="c.commitHash">
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
          <q-badge class="text-xs" :style="{ backgroundColor: hashColor(c.nickname), color: '#fff' }">{{ c.nickname
          }}</q-badge>
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
import { ref, watch, onMounted, onUnmounted, computed } from 'vue';
import { useApi } from '@/composables/useApi';
import { useFormat } from '@/composables/useFormat';
import { useWindow } from '@/composables/useWindow';
import { useColor } from '@/composables/useColor';

const api = useApi()
const fmt = useFormat()
const { open } = useWindow()
const { hashColor } = useColor()

const offset = ref(0)
const limit = ref(20)
const loading = ref(false)
const hasMore = ref(true)  // 新增：跟踪是否还有更多数据
const filterKey = ref(0)   // 新增：用于强制刷新组件的key

const infiniteScroll = ref()

const meta: Ref<any> = ref({})
const commitLogs: Ref<any[]> = ref([])

const props = defineProps({
  filter: Object
})

// 计算属性：生成filter的字符串表示，用于检测实际变化
const filterString = computed(() => {
  return JSON.stringify(props.filter);
});

const gotoCommitUrl = (c: any) => {
  const configs = meta.value.config || []
  const commitUrlTmpl = configs.find((item: any) => c.repoUrl.indexOf(item.domain) !== -1)?.commitUrlTmpl
  if (commitUrlTmpl) {
    const commitUrl = commitUrlTmpl
      .replace('{{.RepoUrl}}', c.repoUrl.replace('.git', ''))
      .replace('{{.BranchName}}', c.branchName)
      .replace('{{.CommitHash}}', c.commitHash);
    open(commitUrl)
  } else {
    const commitUrl = `${c.repoUrl.replace('.git', '')}/commit/${c.commitHash}?branch=${c.branchName}`;
    open(commitUrl)
  }
}

const onLoad = async (index?: number, done?: (stop?: boolean) => void) => {
  // 如果正在加载或没有更多数据，直接返回
  if (loading.value || !hasMore.value) {
    done?.(!hasMore.value);
    return;
  }
  
  loading.value = true;
  
  try {
    // 计算偏移量
    if (index !== undefined && index >= 1) {
      offset.value = (index - 1) * limit.value;
    }
    
    // 调用API获取数据
    const res: any = await api.getCommitLogs(props.filter, offset.value, limit.value);
    
    if (res.code === 200) {
      // 如果是第一页，清空现有数据
      if (offset.value === 0) {
        commitLogs.value = [];
      }
      
      // 添加新数据
      commitLogs.value.push(...res.data);
      meta.value = res.meta;
      
      // 判断是否还有更多数据
      hasMore.value = res.data.length >= limit.value;
      done?.(!hasMore.value);
    } else {
      hasMore.value = false;
      done?.(true);
    }
  } catch (e) {
    console.error('加载提交日志失败:', e);
    hasMore.value = false;
    done?.(true);
  } finally {
    loading.value = false;
  }
}

// 监听filter的实际变化
watch(filterString, () => {
  refreshData();
}, { deep: true })

const refreshData = () => {
  // 重置所有状态
  offset.value = 0;
  hasMore.value = true;
  loading.value = false;
  commitLogs.value = [];
  
  // 变更key以强制重新创建infinite-scroll组件
  filterKey.value += 1;
  
  // 重置组件并加载数据
  if (infiniteScroll.value) {
    infiniteScroll.value.reset();
    // 在下一个事件循环中触发加载，确保组件已重置
    setTimeout(() => {
      onLoad();
    }, 0);
  }
}

const interval = ref()
const autoRefresh = ref(true)

onMounted(() => {
  refreshData();
  interval.value = setInterval(() => {
    if (autoRefresh.value) {
      refreshData();
    }
  }, 30000);
})

onUnmounted(() => {
  if (interval.value) {
    clearInterval(interval.value);
  }
})
</script>
