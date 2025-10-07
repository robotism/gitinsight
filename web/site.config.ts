import fs from "fs";

const env = (() => {
  const env: any = {};
  Object.keys(process.env).forEach((key) => {
    env[key] = process.env[key];
  });
  Object.keys(import.meta.env).forEach((key) => {
    env[key] = import.meta.env[key];
  });
  return Object.keys(env)
    .filter((key: string) => key.startsWith("VITE_") || key.startsWith("NODE_"))
    .reduce((obj: any, key) => {
      obj[key] = env[key];
      return obj;
    }, {});
})();

const CNAME = fs.readFileSync("CNAME", "utf8") || env?.CNAME || "localhost";

export default {
  site: {
    host: CNAME,
    url: "http://" + CNAME,
    // github: "https://github.com/robotism/gitinsight",
  },
  style: {
    header: {
      height: "64px",
    },
    footer: {
      height: "48px",
    },
  },
  locales: {
    en: {
      app: {
        name: "Git Insight",
      },
    },
    zh: {
      app: {
        name: "Git 洞察",
      },
    },
  },
  env: env,
};
