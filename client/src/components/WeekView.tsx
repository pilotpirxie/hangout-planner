import dayjs from "dayjs";
import { useResponsive } from "../hooks/useResponsive";
import { useWeekData } from "../hooks/useWeekData";
import type { TimeSlot } from "../types";
import { TimeSlotCard } from "./TimeSlotCard";

export const WeekView = ({
  timeSlots,
  currentWeek,
  onTimeSlotClick,
}: {
  timeSlots: TimeSlot[];
  currentWeek: Date;
  onTimeSlotClick: (timeSlotId: string) => void;
}) => {
  const weekDays = useWeekData(timeSlots, currentWeek);
  const { screenSize } = useResponsive();

  const getGridColumns = () => {
    if (screenSize === "mobile") return "1fr";
    if (screenSize === "tablet") return "repeat(2, 1fr)";
    return "repeat(7, 1fr)";
  };

  return (
    <div className="week-view-container">
      <div
        className="d-grid"
        style={{
          gridTemplateColumns: getGridColumns(),
          gap: "0.5rem",
        }}>
        {weekDays.map((dayData) => {
          const isToday = dayjs().isSame(dayData.date, "day");
          const sortedSlots = [...dayData.slots].sort((a, b) =>
            a.startTime.localeCompare(b.startTime)
          );

          const formattedDate = screenSize === "desktop"
            ? dayData.dayNumber.toString()
            : dayjs(dayData.date).format("MMM D");

          return (
            <div
              key={dayData.dateString}
              className={`card ${isToday ? "border-primary border-2" : ""}`}>
              <div
                className={`card-header text-center fw-bold ${
                  isToday ? "bg-white text-primary" : ""
                }`}>
                <div>{dayData.dayName}</div>
                <div className="fs-5">{formattedDate}</div>
              </div>
              <div className="card-body p-2">
                <div className="d-flex flex-column gap-2">
                  {sortedSlots.length > 0
                    ? sortedSlots.map((slot) => (
                      <TimeSlotCard
                        key={slot.id}
                        timeSlot={slot}
                        onClick={onTimeSlotClick}
                      />
                    ))
                    : (
                      <div className="text-muted text-center small py-3">
                        No slots
                      </div>
                    )}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
