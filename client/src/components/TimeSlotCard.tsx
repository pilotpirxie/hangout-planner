import type { TimeSlot } from "../types";

interface TimeSlotCardProps {
  timeSlot: TimeSlot;
}

export const TimeSlotCard = ({ timeSlot }: TimeSlotCardProps) => {
  return (
    <div className="btn btn-info mt-2 cursor-pointer">
      {timeSlot.startTime} - {timeSlot.endTime}
    </div>
  );
};
