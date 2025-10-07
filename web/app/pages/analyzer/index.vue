<template>
  <div class="flex nowrap w-full h-full">
    <div class="w-1/4 h-[calc(100vh-120px)]">
      <ConditionFilter v-model:filter="filter" />
    </div>
    <div class="w-2/4 h-[calc(100vh-120px)]">
      <AnalysisView v-model:filter="filter" />
    </div>
    <div class="w-1/4 h-[calc(100vh-120px)]">
      <CommitLogs v-model:filter="filter" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import ConditionFilter from "../../components/aynalisis/ConditionFilter.vue"
import AnalysisView from "../../components/aynalisis/AnalysisView.vue"
import CommitLogs from "../../components/aynalisis/CommitLogs.vue"

const filter = ref({
  since: "",
  until: "",
  branches: [],
  contributors: [],
})

</script>

<style scoped>
:deep(.q-drawer) {
  position: relative !important;
}
</style>
