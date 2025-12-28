import type { ReactNode } from "react";
import { useFadeTransition } from "../hooks/useFadeTransition";

export const Modal = ({
  title,
  isVisible,
  children,
  confirmText,
  onConfirm,
  onClose,
  isConfirmDisabled = false,
}: {
  title: string;
  isVisible: boolean;
  children: ReactNode;
  confirmText: string;
  onConfirm: () => void;
  onClose: () => void;
  isConfirmDisabled?: boolean;
}) => {
  const { shouldRender, isAnimating } = useFadeTransition(isVisible);

  if (!shouldRender) {
    return null;
  }

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
            <div>
              <h5 className="modal-title">{title}</h5>
            </div>

            {children}

            <div className="mt-4 d-flex justify-content-end gap-2">
              <button
                type="button"
                className="btn btn-secondary"
                onClick={onClose}>
                Close
              </button>
              <button
                type="button"
                className="btn btn-primary"
                onClick={onConfirm}
                disabled={isConfirmDisabled}>
                {confirmText}
              </button>
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