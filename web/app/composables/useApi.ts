import moment from "moment";

export const useApi = () => {
  const env: any = useRuntimeConfig().public.env;

  const profile = env.NODE_ENV;

  const baseUrl = (() => {
    if (
      profile === "production" ||
      profile === "prod" ||
      profile === "release"
    ) {
      return env?.BASE_URL || "";
    } else {
      return env?.BASE_URL || "http://localhost:8080";
    }
  })();

  const filterToParams = (filter?: any, offset?: number, limit?: number) => {
    const target: any = {
      offset: offset || filter?.offset,
      limit: limit || filter?.limit,
      since: filter?.since ? filter?.since + " 00:00:00" : "",
      until: filter?.until ? filter?.until + " 23:59:59" : "",
      repos: filter?.repos?.join?.(",") || "",
      branches: filter?.branches?.join?.(",") || "",
      authors: filter?.authors?.join?.(",") || "",
      messageType: filter?.messageType || "",
      isMerge: filter?.isMerge || "",
      period: filter?.period || "day",
      leEffective: filter?.leEffective || "",
      geEffective: filter?.geEffective || "",
    };
    return Object.keys(target)
      .filter((key: string) => target[key] !== "")
      .reduce((obj: any, key: string) => {
        obj[key] = target[key];
        return obj;
      }, {});
  };

  const getCommitLogs = async (filter: any, offset: number, limit: number) => {
    return await $fetch(baseUrl + "/v1/commits", {
      params: filterToParams(filter, offset, limit),
    });
  };

  const getContributors = async (filter?: any) => {
    return await $fetch(baseUrl + "/v1/contributors", {
      params: filterToParams(filter),
    });
  };

  const getRepoBranches = async (filter?: any) => {
    return await $fetch(baseUrl + "/v1/branches", {
      params: filterToParams(filter),
    });
  };

  const getRanking = async (filter?: any) => {
    return await $fetch(baseUrl + "/v1/ranking", {
      params: filterToParams(filter),
    });
  };

  const getCommitHeatmap = async (filter?: any) => {
    return await $fetch(baseUrl + "/v1/heatmap", {
      params: filterToParams(filter),
    });
  };

  const getCommitPeriod = async (filter?: any) => {
    return await $fetch(baseUrl + "/v1/period", {
      params: filterToParams(filter),
    });
  };

  return {
    getCommitLogs,
    getContributors,
    getRepoBranches,
    getRanking,
    getCommitHeatmap,
    getCommitPeriod,
  };
};
