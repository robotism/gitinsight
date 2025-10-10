<template>
  <div class="w-full content-center">
    <div class="w-full max-w-[1600px]">
      <div class="flex nowrap w-full h-full" v-if="$q.screen.gt.sm">
        <div class="w-[320px] h-[calc(100vh-130px)]">
          <ConditionFilter v-model:filter="filter" />
        </div>
        <div class="w-[calc(100%-640px)] h-[calc(100vh-130px)]">
          <AnalysisView v-model:filter="filter" />
        </div>
        <div class="w-[320px] h-[calc(100vh-130px)]">
          <CommitLogs v-model:filter="filter" />
        </div>
      </div>

      <div class="flex nowrap w-full h-full" v-else>
        <div class="w-full content-center">
          <ConditionFilter v-model:filter="filter" />
        </div>
        <div class="w-full min-h-[120px]">
          <AnalysisView v-model:filter="filter" />
        </div>
        <div class="w-full min-h-[120px]">
          <CommitLogs v-model:filter="filter" />
        </div>
      </div>
    </div>
  </div>

</template>

<script lang="ts" setup>
import ConditionFilter from "../../components/aynalisis/ConditionFilter.vue"
import AnalysisView from "../../components/aynalisis/AnalysisView.vue"
import CommitLogs from "../../components/aynalisis/CommitLogs.vue"

import moment from "moment"

const filter = ref({
  since: moment().subtract(2, "month").startOf("month").format("YYYY-MM-DD"),
  until: moment().format("YYYY-MM-DD"),
  branches: [],
  contributors: [],
})

</script>

<style scoped>
:deep(.q-drawer) {
  position: relative !important;
}
</style>
