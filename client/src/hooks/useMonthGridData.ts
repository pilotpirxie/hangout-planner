import dayjs, { type Dayjs } from "dayjs";
import { useMemo } from "react";
import type { TimeSlot } from "../types";

export interface MonthDayData {
  date: Dayjs;
  dateString: string;
  slots: TimeSlot[];
  isCurrentMonth: boolean;
}

export interface WeekRow {
  weekNumber: number;
  days: MonthDayData[];
}

export const useMonthGridData = (
  timeSlots: TimeSlot[],
  currentMonth: Date
): WeekRow[] => {
  return useMemo(() => {
    const currentDayjs = dayjs(currentMonth);
    const year = currentDayjs.year();
    const month = currentDayjs.month();

    const firstDayOfMonth = dayjs(new Date(year, month, 1));
    const firstDayWeekday = firstDayOfMonth.day();
    const mondayOffset = firstDayWeekday === 0 ? 6 : firstDayWeekday - 1;

    const totalCells = 35;

    const cells: MonthDayData[] = [];

    for (let i = 0; i < totalCells; i++) {
      const dayOffset = i - mondayOffset;
      const cellDate = dayjs(new Date(year, month, dayOffset + 1));

      const dateString = cellDate.format("YYYY-MM-DD");
      const isCurrentMonth = cellDate.month() === month;

      const slotsForDate = timeSlots.filter(slot => slot.date === dateString);

      cells.push({
        date: cellDate,
        dateString,
        slots: slotsForDate,
        isCurrentMonth,
      });
    }

    const weeks: WeekRow[] = [];
    for (let i = 0; i < totalCells; i += 7) {
      weeks.push({
        weekNumber: Math.floor(i / 7),
        days: cells.slice(i, i + 7),
      });
    }

    return weeks;
  }, [timeSlots, currentMonth]);
};
