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
      return env?.BASE_URL || "localhost:8080";
    }
  })();

  const getCommitLogs = async (filter: any, offset: number, limit: number) => {
    return await $fetch(baseUrl + "/v1/commits", {
      params: {
        since: filter.since || "",
        until: filter.until || "",
        repos: filter.repos?.join?.(",") || "",
        branches: filter.branches?.join?.(",") || "",
        authors: filter.authors?.join?.(",") || "",
        offset: offset,
        limit: limit,
      },
    });
  };

  const getContributors = async (filter?: any) => {
    return await $fetch(baseUrl + "/v1/contributors", {
      params: {
        since: filter?.since || "",
        until: filter?.until || "",
        repos: filter?.repos?.join?.(",") || "",
        branches: filter?.branches?.join?.(",") || "",
        authors: filter?.authors?.join?.(",") || "",
      },
    });
  };

  const getRepoBranches = async (filter?: any) => {
    return await $fetch(baseUrl + "/v1/branches", {
      params: {
        since: filter?.since || "",
        until: filter?.until || "",
        repos: filter?.repos?.join?.(",") || "",
      },
    });
  };

  const getRanking = async (filter?: any) => {
    return await $fetch(baseUrl + "/v1/ranking", {
      params: {
        since: filter?.since || "",
        until: filter?.until || "",
        repos: filter?.repos?.join?.(",") || "",
        branches: filter?.branches?.join?.(",") || "",
        authors: filter?.authors?.join?.(",") || "",
      },
    });
  };

  const getCommitHeatmap = async (filter?: any) => {
    return await $fetch(baseUrl + "/v1/heatmap", {
      params: {
        since: filter?.since || "",
        until: filter?.until || "",
        repos: filter?.repos?.join?.(",") || "",
        branches: filter?.branches?.join?.(",") || "",
        authors: filter?.authors?.join?.(",") || "",
        messageType: filter?.messageType || "",
        isMerge: filter?.isMerge || "",
      },
    });
  };

  return {
    getCommitLogs,
    getContributors,
    getRepoBranches,
    getRanking,
    getCommitHeatmap,
  };
};
