<template>
  <div class="w-full h-full flex flex-col justify-center items-center p-4">
    <q-toggle class="text-xs ml-auto mr-32" dense v-model="autoRefresh" :label="$t('autoRefresh')" />
  </div>
  <div class="w-full content-center">
    <div class="max-w-[1600px] flex flex-row content-center">

      <ProjectsCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.todayProjects')"
        :projects="todayProjects" />
      <RankingCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.todayRanking')"
        :ranking="todayRanking" />

      <ProjectsCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.yesterdayProjects')"
        :projects="yesterdayProjects" />
      <RankingCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.yesterdayRanking')"
        :ranking="yesterdayRanking" />

      <ProjectsCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.weekProjects')"
        :projects="weekProjects" />
      <RankingCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.weekRanking')"
        :ranking="weekRanking" />

      <ProjectsCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.lastWeekProjects')"
        :projects="lastWeekProjects" />
      <RankingCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.lastWeekRanking')"
        :ranking="lastWeekRanking" />

      <ProjectsCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.monthProjects')"
        :projects="monthProjects" />
      <RankingCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.monthRanking')"
        :ranking="monthRanking" />

      <ProjectsCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.lastMonthProjects')"
        :projects="lastMonthProjects" />
      <RankingCard :class="$q.screen.gt.md ? 'card' : 'card-mobile'" :title="$t('dashboard.lastMonthRanking')"
        :ranking="lastMonthRanking" />

    </div>
  </div>

</template>

<script setup>
import ProjectsCard from "../components/aynalisis/ProjectsCard.vue"
import RankingCard from "../components/aynalisis/RankingCard.vue"

import moment from "moment"

const api = useApi()

const todayProjects = ref([])
const todayRanking = ref([])
const yesterdayProjects = ref([])
const yesterdayRanking = ref([])

const weekProjects = ref([])
const weekRanking = ref([])
const lastWeekProjects = ref([])
const lastWeekRanking = ref([])

const monthProjects = ref([])
const monthRanking = ref([])
const lastMonthProjects = ref([])
const lastMonthRanking = ref([])


function getDateRange(type) {
  const FORMAT = "YYYY-MM-DD HH:mm:ss"
  const m = moment()

  switch (type) {
    case "today":
      return [m.startOf("day").format(FORMAT), m.endOf("day").format(FORMAT)]
    case "yesterday":
      return [
        m.clone().subtract(1, "day").startOf("day").format(FORMAT),
        m.clone().subtract(1, "day").endOf("day").format(FORMAT)
      ]
    case "week":
      return [m.startOf("week").format(FORMAT), m.endOf("week").format(FORMAT)]
    case "lastWeek":
      return [
        m.clone().subtract(1, "week").startOf("week").format(FORMAT),
        m.clone().subtract(1, "week").endOf("week").format(FORMAT)
      ]
    case "month":
      return [m.startOf("month").format(FORMAT), m.endOf("month").format(FORMAT)]
    case "lastMonth":
      return [
        m.clone().subtract(1, "month").startOf("month").format(FORMAT),
        m.clone().subtract(1, "month").endOf("month").format(FORMAT)
      ]
    default:
      return []
  }
}


const refreshData = () => {
  const [todaySince, todayUntil] = getDateRange("today")
  const [yesterdaySince, yesterdayUntil] = getDateRange("yesterday")
  const [weekSince, weekUntil] = getDateRange("week")
  const [lastWeekSince, lastWeekUntil] = getDateRange("lastWeek")
  const [monthSince, monthUntil] = getDateRange("month")
  const [lastMonthSince, lastMonthUntil] = getDateRange("lastMonth")

  api.getRepoBranches({ since: todaySince, until: todayUntil }).then(res => {
    if (res.code === 200) {
      todayProjects.value = res.data
    }
  })

  api.getRanking({ since: todaySince, until: todayUntil }).then(res => {
    if (res.code === 200) {
      todayRanking.value = res.data
    }
  })

  api.getRepoBranches({ since: yesterdaySince, until: yesterdayUntil }).then(res => {
    if (res.code === 200) {
      yesterdayProjects.value = res.data
    }
  })

  api.getRanking({ since: yesterdaySince, until: yesterdayUntil }).then(res => {
    if (res.code === 200) {
      yesterdayRanking.value = res.data
    }
  })

  api.getRepoBranches({ since: weekSince, until: weekUntil }).then(res => {
    if (res.code === 200) {
      weekProjects.value = res.data
    }
  })

  api.getRanking({ since: weekSince, until: weekUntil }).then(res => {
    if (res.code === 200) {
      weekRanking.value = res.data
    }
  })

  api.getRepoBranches({ since: lastWeekSince, until: lastWeekUntil }).then(res => {
    if (res.code === 200) {
      lastWeekProjects.value = res.data
    }
  })

  api.getRanking({ since: lastWeekSince, until: lastWeekUntil }).then(res => {
    if (res.code === 200) {
      lastWeekRanking.value = res.data
    }
  })

  api.getRepoBranches({ since: monthSince, until: monthUntil }).then(res => {
    if (res.code === 200) {
      monthProjects.value = res.data
    }
  })

  api.getRanking({ since: monthSince, until: monthUntil }).then(res => {
    if (res.code === 200) {
      monthRanking.value = res.data
    }
  })

  api.getRepoBranches({ since: lastMonthSince, until: lastMonthUntil }).then(res => {
    if (res.code === 200) {
      lastMonthProjects.value = res.data
    }
  })

  api.getRanking({ since: lastMonthSince, until: lastMonthUntil }).then(res => {
    if (res.code === 200) {
      lastMonthRanking.value = res.data
    }
  })
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


<style scoped>
.card {
  width: 360px;
  height: 270px;
}

.card-mobile {
  width: 96vw;
  height: 270px;
}
</style>