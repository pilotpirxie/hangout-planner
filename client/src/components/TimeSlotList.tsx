import type { TimeSlot } from "../types";

export const TimeSlotList = ({
  timeSlots,
  onDelete,
  onEdit
}: {
  timeSlots: TimeSlot[];
  onDelete: (id: string) => void;
  onEdit: (slot: TimeSlot) => void;
  }) => {
  if (timeSlots.length === 0) {
    return null;
  }

  const sortedSlots = [...timeSlots].sort((a, b) => {
    const dateTimeA = `${a.date}T${a.startTime}`;
    const dateTimeB = `${b.date}T${b.startTime}`;
    return dateTimeA.localeCompare(dateTimeB);
  });

  return (
    <div className="mt-2 max-h-300px overflow-y-auto bg-light card card-body">
      {sortedSlots.map((slot) => (
        <div
          key={slot.id}
          className="card card-body mt-2">
          <div className="d-flex justify-content-between align-items-center">
            <div>
              {slot.date} from {slot.startTime} to {slot.endTime}
            </div>
            <div className="d-flex gap-2">
              <button
                className="btn btn-sm btn-primary"
                onClick={() => { onEdit(slot); }}>
                <i className="ri-edit-line" />
              </button>
              <button
                className="btn btn-sm btn-danger"
                onClick={() => { onDelete(slot.id); }}>
                <i className="ri-delete-bin-6-line" />
              </button>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
