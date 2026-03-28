import dayjs from "dayjs";
import isSameOrBefore from "dayjs/plugin/isSameOrBefore";
import type { TimeSlot } from "../types";

dayjs.extend(isSameOrBefore);

export const generateTimeSlots = ({
  startDate,
  endDate,
  dailyStartTime,
  dailyEndTime,
  duration,
  isOverlapping,
  isWholeDay,
}: {
  startDate: string;
  endDate: string;
  dailyStartTime: string;
  dailyEndTime: string;
  duration: string;
  isOverlapping: boolean;
  isWholeDay: boolean;
}): TimeSlot[] => {
  const generated: TimeSlot[] = [];
  const start = dayjs(startDate);
  const end = dayjs(endDate);

  if (isWholeDay) {
    for (let date = start; date.isSameOrBefore(end, "day"); date = date.add(1, "day")) {
      const dateStr = date.format("YYYY-MM-DD");
      generated.push({
        id: crypto.randomUUID(),
        startDate: dayjs(`${dateStr}T00:00:00`).toDate(),
        endDate: dayjs(`${dateStr}T23:59:59`).toDate(),
      });
    }
    return generated;
  }

  const durationHours = parseFloat(duration);
  const intervalHours = isOverlapping
    ? durationHours / 2
    : durationHours;

  const [dailyStartHour, dailyStartMin] = dailyStartTime.split(":").map(Number);
  const [dailyEndHour, dailyEndMin] = dailyEndTime.split(":").map(Number);
  const dailyStartMinutes = (dailyStartHour * 60) + dailyStartMin;
  const dailyEndMinutes = (dailyEndHour * 60) + dailyEndMin;

  for (let date = start; date.isSameOrBefore(end, "day"); date = date.add(1, "day")) {
    const dateStr = date.format("YYYY-MM-DD");

    let currentMinutes = dailyStartMinutes;
    while (currentMinutes + (durationHours * 60) <= dailyEndMinutes) {
      const startHour = Math.floor(currentMinutes / 60);
      const startMin = currentMinutes % 60;
      const endMinutes = currentMinutes + (durationHours * 60);
      const endHour = Math.floor(endMinutes / 60);
      const endMin = endMinutes % 60;

      const startTime = `${String(startHour).padStart(2, "0")}:${String(startMin).padStart(2, "0")}`;
      const endTime = `${String(endHour).padStart(2, "0")}:${String(endMin).padStart(2, "0")}`;

      const slotStart = `${dateStr}T${startTime}:00`;
      const slotEnd = `${dateStr}T${endTime}:00`;

      generated.push({
        id: crypto.randomUUID(),
        startDate: dayjs(slotStart).toDate(),
        endDate: dayjs(slotEnd).toDate(),
      });

      currentMinutes += intervalHours * 60;
    }
  }

  return generated;
};
