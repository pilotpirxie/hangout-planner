import { formatMonthYear, formatWeekRange } from "../utils/dateFormatters";

type ViewMode = "week" | "month";

const ViewModeToggle = ({
  viewMode,
  onViewModeChange
}: {
  viewMode: ViewMode;
  onViewModeChange: (mode: ViewMode) => void;
}) => (
  <div
    className="btn-group"
    role="group">
    <button
      type="button"
      className={`btn ${viewMode === "week" ? "btn-primary" : "btn-outline-primary"}`}
      onClick={() => { onViewModeChange("week"); }}>
      Week View
    </button>
    <button
      type="button"
      className={`btn ${viewMode === "month" ? "btn-primary" : "btn-outline-primary"}`}
      onClick={() => { onViewModeChange("month"); }}>
      Month View
    </button>
  </div>
);

const NavigationControls = ({
  label,
  onNavigate,
}: {
  label: string;
  onNavigate: (direction: "prev" | "next") => void;
}) => (
  <div className="d-flex align-items-center gap-2">
    <button
      className="btn btn-sm btn-outline-secondary"
      onClick={() => { onNavigate("prev"); }}>
      <i className="ri-arrow-left-s-line" />
    </button>
    <span
      className="fw-bold"
      style={{ minWidth: "150px", textAlign: "center" }}>
      {label}
    </span>
    <button
      className="btn btn-sm btn-outline-secondary"
      onClick={() => { onNavigate("next"); }}>
      <i className="ri-arrow-right-s-line" />
    </button>
  </div>
);

export const CalendarHeader = ({
  eventName,
  viewMode,
  onViewModeChange,
  currentMonth,
  currentWeek,
  onMonthChange,
  onWeekChange,
  onGoToToday,
}: {
  eventName: string;
  viewMode: ViewMode;
  onViewModeChange: (mode: ViewMode) => void;
  currentMonth?: Date;
  currentWeek?: Date;
  onMonthChange?: (direction: "prev" | "next") => void;
  onWeekChange?: (direction: "prev" | "next") => void;
  onGoToToday?: () => void;
}) => {
  const displayLabel =
    viewMode === "week" && currentWeek
      ? formatWeekRange(currentWeek)
      : viewMode === "month" && currentMonth
        ? formatMonthYear(currentMonth)
        : "";

  const handleNavigate = viewMode === "week" ? onWeekChange : onMonthChange;
  // eslint-disable-next-line @typescript-eslint/prefer-nullish-coalescing
  const showNavigation = (viewMode === "week" && currentWeek && onWeekChange) ||
    (viewMode === "month" && currentMonth && onMonthChange);

  return (
    <div className="rounded-top p-3 d-flex flex-column align-items-center justify-content-center">
      <h1>{eventName}</h1>
      <h3 className="mb-3 text-center">Choose the date that works best for you</h3>

      <div className="d-flex flex-column align-items-center">
        <div className="d-flex gap-3 align-items-center">
          <ViewModeToggle
            viewMode={viewMode}
            onViewModeChange={onViewModeChange}
          />

          {onGoToToday ? <button
            type="button"
            className="btn btn-outline-success"
            onClick={onGoToToday}>
            <i className="ri-calendar-check-line me-1" />
            Today
          </button> : null}
        </div>

        {showNavigation && handleNavigate ? <div className="mt-4">
          <NavigationControls
            label={displayLabel}
            onNavigate={handleNavigate}
          />
        </div> : null}
      </div>
    </div>
  );
};
