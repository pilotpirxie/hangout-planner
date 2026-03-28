import dayjs from "dayjs";
import type { TimeSlot } from "../types";
import { Modal } from "./Modal";

export const TimeSlotConfirmationModal = ({
  isVisible,
  selectedTimeSlot,
  nickname,
  onNicknameChange,
  onClose,
  onConfirm,
}: {
  isVisible: boolean;
  selectedTimeSlot: TimeSlot | null;
  nickname: string;
  onNicknameChange: (nickname: string) => void;
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
        Type your name or nickname below to confirm. It will be shared with others attending the event.
      </p>
      <input
        type="text"
        className="form-control"
        placeholder="Your name or nickname"
        value={nickname}
        onChange={(e) => { onNicknameChange(e.target.value); }}
      />
    </Modal>
  );
};
