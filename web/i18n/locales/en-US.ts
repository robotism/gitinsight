import moment from "moment";
import "moment/dist/locale/en-gb";

import quasarLangEn from "quasar/lang/en-US.js";

if (!moment.locales().includes("en")) {
  moment.defineLocale("en", {
    parentLocale: "en-gb",
  });
}

export default defineI18nLocale(async (locale) => {
  return {
    nav: {
      "": "Home",
      home: "Home",
      analyzer: "Analyzer",
      contributors: "Contributors",
    },
    contributors: "Contributors",
    nickname: "Nickname",
    name: "Name",
    email: "Email",
    additions: "Additions",
    deletions: "Deletions",
    effectives: "Effectives",
    commits: "Commits",
    projects: "Projects",
    since: "Since",
    until: "Until",
    now: "Now",
    ranking: "Ranking",
    timeRange: "Time Range",
    repos: "Repos",
    branches: "Branches",
    authors: "Authors",

    commitLogs: "Commit Logs",
    conditionFilter: "Condition Filter",
    analysisView: "Statistics Result",
    commitPeriod: "Commit Statistics",
    commitHeatmap: "Commit Heatmap",
    fixHeatmap: "Fix Heatmap",
    featHeatmap: "Feat Heatmap",
    mergeHeatmap: "Merge Heatmap",

    all: "All",
    default: "Default",

    today: "Today",
    yesterday: "Yesterday",
    weekThis: "This Week",
    weekLast: "Last Week",
    weekBeforeLast: "Before Last Week",
    monthThis: "This Month",
    monthLast: "Last Month",
    monthBeforeLast: "Before Last Month",

    autoRefresh: "Auto Refresh",

    dashboard: {
      todayProjects: "Today Projects",
      todayRanking: "Today Ranking",
      yesterdayProjects: "Yesterday Projects",
      yesterdayRanking: "Yesterday Ranking",
      weekProjects: "Week Projects",
      weekRanking: "Week Ranking",
      lastWeekProjects: "Last Week Projects",
      lastWeekRanking: "Last Week Ranking",
      monthProjects: "Month Projects",
      monthRanking: "Month Ranking",
      lastMonthProjects: "Last Month Projects",
      lastMonthRanking: "Last Month Ranking",
    },
    quasar: {
      ...quasarLangEn,
    },
    ...useRuntimeConfig().public.locales[locale],
  };
});
