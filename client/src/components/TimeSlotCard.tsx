import dayjs from "dayjs";
import type { TimeSlot } from "../types";

export const TimeSlotCard = ({ timeSlot, onClick }: {
  timeSlot: TimeSlot;
  onClick: (timeSlotId: string) => void;
}) => {
  return (
    <div
      className="btn btn-info mt-2 cursor-pointer"
      onClick={() => { onClick(timeSlot.id); }}>
      <div>{dayjs(timeSlot.startDate).format("HH:mm")} - {dayjs(timeSlot.endDate).format("HH:mm")}</div>
    </div>
  );
};
