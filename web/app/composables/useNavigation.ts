export const useNavigation = () => {
  const router = useRouter();
  const localePath = useLocalePath();

  // https://icons8.com/line-awesome
  const links = ref([
    {
      label: "nav.home",
      icon: "la la-compass",
      root: true,
      click: () => {
        router.push(localePath("/"));
      },
    },
    {
      label: "nav.analyzer",
      icon: "la la-heartbeat",
      click: () => {
        router.push(localePath("/analyzer"));
      },
    },
    {
      label: "nav.contributors",
      icon: "la la-users",
      click: () => {
        router.push(localePath("/contributors"));
      },
    },
  ]);

  return {
    links,
  };
};
