import type { TimeSlot } from "../types";

export const TimeSlotCard = ({ timeSlot, onClick }: {
  timeSlot: TimeSlot;
  onClick: (timeSlotId: string) => void;
}) => {
  return (
    <div
      className="btn btn-info mt-2 cursor-pointer"
      onClick={() => { onClick(timeSlot.id); }}>
      <div>{timeSlot.startTime} - {timeSlot.endTime}</div>
    </div>
  );
};
