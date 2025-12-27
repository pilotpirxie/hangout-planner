import { useMemo } from "react";
import type { TimeSlot } from "../types";

interface DaySlots {
  date: string;
  slots: TimeSlot[];
}

export const useGroupedTimeSlots = (timeSlots: TimeSlot[]): DaySlots[] => {
  return useMemo(() => {
    const groupedByDay = new Map<string, TimeSlot[]>();

    timeSlots.forEach(slot => {
      if (!groupedByDay.has(slot.date)) {
        groupedByDay.set(slot.date, []);
      }
      groupedByDay.get(slot.date)?.push(slot);
    });

    return Array.from(groupedByDay, ([date, slots]) => ({ date, slots }))
      .sort((a, b) => a.date.localeCompare(b.date));
  }, [timeSlots]);
};
