export const useLocalePage = () => {
  const { defaultLocale } = useI18n();
  const router = useRouter();
  const localeRoute = useLocaleRoute();

  const getLocalePath = (): string => {
    const loc = localeRoute("/")?.path || "";
    return loc.length <= 1 ? "/" + defaultLocale : loc;
  };

  const getPagePath = (withLocale: boolean | undefined = false): string => {
    const loc = getLocalePath();
    const cur = router.currentRoute.value.path;
    if (withLocale) {
      return cur.startsWith(loc) ? cur : loc + cur;
    } else {
      return cur.replace(loc, "");
    }
  };

  const getPagePathArray = (
    withLocale: boolean | undefined = false
  ): string[] => {
    return getPagePath(withLocale).replace(/^\//, "").split("/");
  };

  return {
    getLocalePath,
    getPagePath,
    getPagePathArray,
  };
};
