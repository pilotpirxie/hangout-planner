import dayjs from "dayjs";

export const formatMonthYear = (date: Date): string => {
  return dayjs(date).format("MMMM YYYY");
};

export const formatWeekRange = (date: Date): string => {
  const day = dayjs(date);
  const start = day.day() === 0 ? day.subtract(6, "days") : day.subtract(day.day() - 1, "days");
  const end = start.add(6, "days");

  return `${start.format("MMM D")} - ${end.format("MMM D")}`;
};
