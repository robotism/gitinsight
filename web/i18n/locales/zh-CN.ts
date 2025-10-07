import moment from "moment";
import "moment/dist/locale/zh-cn";

import quasarLangZh from "quasar/lang/zh-CN.js";

if (!moment.locales().includes("zh")) {
  moment.defineLocale("zh", {
    parentLocale: "zh-cn",
  });
}

export default defineI18nLocale(async (locale) => {
  return {
    nav: {
      "": "首页",
      home: "首页",
      analyzer: "分析器",
      contributors: "贡献者",
    },
    contributors: "贡献者",
    nickname: "昵称",
    name: "姓名",
    email: "邮箱",
    additions: "添加",
    deletions: "删除",
    effectives: "有效",
    commits: "提交",
    projects: "项目",
    since: "自",
    until: "至",
    now: "现在",
    ranking: "排行",
    timeRange: "时间范围",
    repos: "项目",
    branches: "分支",
    authors: "作者",

    commitLogs: "提交日志",
    conditionFilter: "条件过滤",
    analysisView: "统计结果",
    fixHeatmap: "修复热度",
    featHeatmap: "功能热度",
    mergeHeatmap: "合并热度",
    commitHeatmap: "提交热度",

    all: "全部",

    today: "今日",
    yesterday: "昨日",
    weekThis: "本周",
    weekLast: "上周",
    weekBeforeLast: "上上周",
    monthThis: "本月",
    monthLast: "上月",
    monthBeforeLast: "上上月",

    dashboard: {
      todayProjects: "今日贡献项目",
      todayRanking: "今日贡献排行",
      yesterdayProjects: "昨日贡献项目",
      yesterdayRanking: "昨日贡献排行",
      weekProjects: "本周贡献项目",
      weekRanking: "本周贡献排行",
      lastWeekProjects: "上周贡献项目",
      lastWeekRanking: "上周贡献排行",
      monthProjects: "本月贡献项目",
      monthRanking: "本月贡献排行",
      lastMonthProjects: "上月贡献项目",
      lastMonthRanking: "上月贡献排行",
    },
    quasar: {
      ...quasarLangZh,
    },
    ...useRuntimeConfig().public.locales[locale],
  };
});
