import type { Dayjs } from "dayjs";
import dayjs from "dayjs";
import type { MonthDayData } from "../hooks/useMonthGridData";
import { useResponsive } from "../hooks/useResponsive";
import type { TimeSlot } from "../types";

interface MonthDayCellProps {
  dayData: MonthDayData;
  onDayClick: (slots: TimeSlot[], date: string) => void;
}

export const MonthDayCell = ({ dayData, onDayClick }: MonthDayCellProps) => {
  const { date, slots, isCurrentMonth, dateString } = dayData;
  const dateObj: Dayjs = date;
  const dayNumber = dateObj.date();
  const hasSlots = slots.length > 0;
  const { screenSize } = useResponsive();

  const isToday = dayjs().isSame(dateObj, "day");

  const dayName = dateObj.format("ddd");
  const showDayName = screenSize !== "desktop";

  return (
    <div
      className={`month-day-cell border rounded p-2 
        ${!isCurrentMonth ? "text-muted bg-light" : "bg-white"} 
        ${isToday ? "border-primary border-2" : ""} ${hasSlots ? "cursor-pointer" : ""}`}
      onClick={() => {
        if (hasSlots) {
          onDayClick(slots, dateString);
        }
      }}
      style={{
        minHeight: "100px",
        cursor: hasSlots ? "pointer" : "default",
      }}>
      <div className="d-flex justify-content-between align-items-start">
        <div>
          {showDayName ? <div className="small text-muted">{dayName}</div> : null}
          <span className={`fw-bold ${isToday ? "text-primary fs-5" : ""}`}>
            {showDayName
              ? dateObj.format("MMM D")
              : dayNumber}
          </span>
        </div>
        {hasSlots
          ? (
            <span className="badge bg-primary rounded-pill">
              {slots.length}
            </span>
          )
          : null}
      </div>
      {hasSlots
        ? (
          <div className="mt-2 small text-truncate">
            {slots.slice(0, 2).map((slot) => (
              <div
                key={slot.id}
                className="text-truncate">
                {slot.startTime}
              </div>
            ))}
            {slots.length > 2
              ? (
                <div className="text-muted fst-italic">+{slots.length - 2} more</div>
              )
              : null}
          </div>
        )
        : null}
    </div>
  );
};
