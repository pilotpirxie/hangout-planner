import { useMonthGridData } from "../hooks/useMonthGridData";
import { useResponsive } from "../hooks/useResponsive";
import type { TimeSlot } from "../types";
import { MonthDayCell } from "./MonthDayCell";

interface MonthGridViewProps {
  timeSlots: TimeSlot[];
  currentMonth: Date;
  onDayClick: (slots: TimeSlot[], date: string) => void;
}

const WEEKDAY_LABELS = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];

export const MonthGridView = ({
  timeSlots,
  currentMonth,
  onDayClick,
}: MonthGridViewProps) => {
  const weeks = useMonthGridData(timeSlots, currentMonth);
  const { screenSize } = useResponsive();

  const getGridColumns = () => {
    if (screenSize === "mobile") return "1fr";
    if (screenSize === "tablet") return "repeat(2, 1fr)";
    return "repeat(7, 1fr)";
  };

  const showHeaders = screenSize === "desktop";

  const handleDayClick = (slots: TimeSlot[], date: string) => {
    if (slots.length > 0) {
      onDayClick(slots, date);
    }
  };

  return (
    <div className="month-grid-container">
      {showHeaders ? <div
        className="d-grid"
        style={{
          gridTemplateColumns: getGridColumns(),
          gap: "0.5rem",
        }}>
        {WEEKDAY_LABELS.map((label) => (
          <div
            key={label}
            className="bg-white text-center fw-bold p-2 border rounded">
            {label}
          </div>
        ))}
      </div> : null}

      <div
        className="d-grid mt-2"
        style={{
          gridTemplateColumns: getGridColumns(),
          gap: "0.5rem",
        }}>
        {weeks.map((week) =>
          week.days.map((dayData) => (
            <MonthDayCell
              key={dayData.dateString}
              dayData={dayData}
              onDayClick={(slots, date) => { handleDayClick(slots, date); }}
            />
          ))
        )}
      </div>
    </div>
  );
};
