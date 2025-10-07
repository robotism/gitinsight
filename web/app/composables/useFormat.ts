import moment from "moment";

export const useFormat = () => {
  const { locale } = useI18n();

  watch(
    () => locale.value,
    () => {
      moment.locale(locale.value);
    }
  );

  const localDate = (
    date: string | Date | undefined | null,
    fmt: string | undefined | null = null
  ) => {
    fmt = fmt || "YYYY-MM-DD HH:mm:ss";
    return moment(date || new Date()).format(fmt);
  };

  const fromNow = (date: string | Date | undefined | null) => {
    return moment(date || new Date()).fromNow();
  };

  return {
    localDate,
    fromNow,
  };
};
