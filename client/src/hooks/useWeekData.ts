import dayjs, { type Dayjs } from "dayjs";
import { useMemo } from "react";
import type { TimeSlot } from "../types";

export interface WeekDayData {
  date: Dayjs;
  dateString: string;
  dayName: string;
  dayNumber: number;
  slots: TimeSlot[];
}

export const useWeekData = (
  timeSlots: TimeSlot[],
  currentWeek: Date
): WeekDayData[] => {
  return useMemo(() => {
    let startOfWeek = dayjs(currentWeek);
    const day = startOfWeek.day();
    const daysToMonday = day === 0 ? 6 : day - 1;
    startOfWeek = startOfWeek.subtract(daysToMonday, "days").startOf("day");

    const weekDays: WeekDayData[] = [];

    for (let i = 0; i < 7; i++) {
      const date = startOfWeek.add(i, "days");

      const dateString = date.format("YYYY-MM-DD");

      const dayNames = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];
      const dayIndex = date.day();
      const dayName = dayNames[dayIndex === 0 ? 6 : dayIndex - 1];
      const dayNumber = date.date();

      const slotsForDate = timeSlots.filter(slot => slot.date === dateString);

      weekDays.push({
        date,
        dateString,
        dayName,
        dayNumber,
        slots: slotsForDate,
      });
    }

    return weekDays;
  }, [timeSlots, currentWeek]);
};
