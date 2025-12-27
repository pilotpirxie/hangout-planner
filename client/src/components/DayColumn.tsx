import type { TimeSlot } from "../types";
import { TimeSlotCard } from "./TimeSlotCard";

interface DayColumnProps {
  date: string;
  slots: TimeSlot[];
}

export const DayColumn = ({ date, slots }: DayColumnProps) => {
  const sortedSlots = [...slots].sort((a, b) => {
    return a.startTime.localeCompare(b.startTime);
  });

  return (
    <div className="day-column">
      <div className="text-center fw-bold mb-2 p-2">
        {date}
      </div>
      <div className="d-flex flex-column gap-2">
        {sortedSlots.map(timeSlot => (
          <TimeSlotCard
            key={timeSlot.id}
            timeSlot={timeSlot}
          />
        ))}
      </div>
    </div>
  );
};
