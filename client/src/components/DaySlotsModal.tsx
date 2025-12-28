import dayjs from "dayjs";
import { useFadeTransition } from "../hooks/useFadeTransition";
import type { TimeSlot } from "../types";

export const DaySlotsModal = ({
  isVisible,
  date,
  slots,
  onClose,
  onSelectSlot,
}: {
  isVisible: boolean;
  date: string;
  slots: TimeSlot[];
  onClose: () => void;
  onSelectSlot: (timeSlotId: string) => void;
}) => {
  const { shouldRender, isAnimating } = useFadeTransition(isVisible);

  if (!shouldRender) {
    return null;
  }

  const formattedDate = date
    ? dayjs(date).format("dddd, MMMM D, YYYY")
    : "";

  const sortedSlots = [...slots].sort((a, b) =>
    a.startTime.localeCompare(b.startTime)
  );

  return <>
    <div
      className={`modal fade ${isAnimating ? "show" : ""} d-block`}
      tabIndex={-1}
      role="dialog"
      style={{ zIndex: 1050 }}>
      <div
        className="modal-dialog"
        role="document">
        <div className="modal-content">
          <div className="modal-body">
            <div className="d-flex justify-content-between align-items-center mb-3">
              <h5 className="modal-title mb-0">Choose a Time Slot</h5>
              <button
                type="button"
                className="btn-close"
                onClick={onClose}
                aria-label="Close"
              />
            </div>
            <p className="text-muted mb-3">{formattedDate}</p>

            <div className="d-flex flex-column gap-2">
              {sortedSlots.map((slot) => (
                <button
                  key={slot.id}
                  className="btn btn-info text-start"
                  onClick={() => {
                    onSelectSlot(slot.id);
                    onClose();
                  }}>
                  <div className="fw-bold">
                    {slot.startTime} - {slot.endTime}
                  </div>
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
    <div
      className={`modal-backdrop fade ${isAnimating ? "show" : ""}`}
      style={{ zIndex: 1040 }}
    />
  </>;
};
