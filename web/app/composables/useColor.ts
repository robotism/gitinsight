export const useColor = () => {
  const hashColor = (str: string) => {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    const hue = Math.abs(hash) % 360;
    return `hsl(${hue}, 70%, 55%)`; // 饱和度70%，亮度55%
  };
  const rankColor = (index: number, total: number) => {
    const goldShades = ["#FFB300", "#FFC107", "#FFD700"]; // 前三名金色
    const blueColor = "#60A5FA"; // 中间默认蓝色
    const dangerShades = ["#FF7043", "#F44336", "#D32F2F"]; // 最后三名红色/橙色系

    if (index < 3) return goldShades[index]; // 前三名
    if (total > 0 && index >= total - 3)
      return dangerShades[index - (total - 3)]; // 最后三名
    return blueColor; // 其他默认蓝色
  };

  return {
    hashColor,
    rankColor,
  };
};
