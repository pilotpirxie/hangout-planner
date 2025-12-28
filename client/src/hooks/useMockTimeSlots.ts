import { useState } from "react";
import type { TimeSlot } from "../types";

export const useMockTimeSlots = () => {
  const [timeSlots] = useState<TimeSlot[]>(() =>
    Array.from({ length: 50 }, (_, i) => {
      const day = Math.floor(Math.random() * 28) + 1;
      const startHour = Math.floor(Math.random() * 8) + 9;
      const endHour = startHour + 1;
      return {
        id: `slot-${(i + 1).toString()}`,
        date: `2025-12-${day.toString().padStart(2, "0")}`,
        startTime: `${startHour.toString()}:00`,
        endTime: `${endHour.toString()}:00`,
      };
    })
  );

  return timeSlots;
};
