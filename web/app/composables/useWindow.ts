export const useWindow = () => {
  const router = useRouter();
  const win: Ref<Window | null> = ref(null);
  const visible = ref(true);
  const { getPagePathArray } = useLocalePage();
  const { t, te } = useI18n();

  const size = ref({});

  const open = (url: string, target: string | undefined = "_blank") => {
    if (!url) {
      return;
    }
    if (url.toLowerCase().startsWith("http")) {
      win.value?.open?.(url, target);
    } else {
      const route = router.resolve({
        name: url,
      });
      win.value?.open?.(route.href, target);
    }
  };

  const onResize =async (_?: Event) => {
    size.value = {
      width: win.value?.innerWidth,
      height: win.value?.innerHeight,
    };
  };

  const onVisibilityChange = (_: Event) => {
    if (document.visibilityState === "hidden") {
      visible.value = false;
      return;
    }
    if (document.visibilityState === "visible") {
      visible.value = true;
      return;
    }
  };

  const computedTitle = (customTitle: string | null = null) => {
    return computed(() => {
      const prefix = !visible.value ? "🤡 " : "👻 ";
      const suffix = "丨" + t("app.name");
      const root = getPagePathArray()[0];
      const title = te("nav." + root) ? t("nav." + root) : t("nav.home");
      return prefix + (customTitle || title) + suffix;
    });
  };

  onMounted(() => {
    win.value = window;
    window.addEventListener("resize", onResize);
    window.addEventListener("visibilitychange", onVisibilityChange);
    window.dispatchEvent(
      new Event("resize", { bubbles: false, cancelable: false })
    );
    window.dispatchEvent(
      new Event("visibilitychange", { bubbles: false, cancelable: false })
    );
     onResize()
  });

  onUnmounted(() => {
    win.value = null;
    window.removeEventListener("resize", onResize);
    window.removeEventListener("visibilitychange", onVisibilityChange);
  });

  return {
    window,
    visible,
    size,
    open,
    computedTitle,
  };
};
