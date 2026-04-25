import dayjs from "dayjs";
import type { TimeSlot } from "../types";
import { Modal } from "./Modal";

export const TimeSlotConfirmationModal = ({
  isVisible,
  selectedTimeSlot,
  username,
  onUsernameChange,
  onClose,
  onConfirm,
}: {
  isVisible: boolean;
  selectedTimeSlot: TimeSlot | null;
  username: string;
  onUsernameChange: (username: string) => void;
  onClose: () => void;
  onConfirm: () => void;
}) => {
  return (
    <Modal
      confirmText="Confirm"
      onConfirm={onConfirm}
      isVisible={isVisible}
      onClose={onClose}
      title="Are you sure?">
      <p>
        Are you sure you want to confirm the time slot on{" "}
        <strong>{dayjs(selectedTimeSlot?.startDate).format("DD-MM-YYYY HH:mm")}</strong> to{" "}
        <strong>{dayjs(selectedTimeSlot?.endDate).format("DD-MM-YYYY HH:mm")}</strong>?
      </p>
      <p>
        Type your name or username below to confirm. It will be shared with others attending the event.
      </p>
      <input
        type="text"
        className="form-control"
        placeholder="Your name or username"
        value={username}
        onChange={(e) => { onUsernameChange(e.target.value); }}
      />
    </Modal>
  );
};
