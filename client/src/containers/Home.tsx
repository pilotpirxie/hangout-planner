import { Collapse } from "../components/Collapse";
import { Modal } from "../components/Modal";
import { QuickSlotGenerator } from "../components/QuickSlotGenerator";
import { TimeSlotForm } from "../components/TimeSlotForm";
import { TimeSlotList } from "../components/TimeSlotList";
import { useHangoutForm } from "../hooks/useHangoutForm";
import { useQuickSlotModal } from "../hooks/useQuickSlotModal";
import { useTimeSlotModal } from "../hooks/useTimeSlotModal";
import type { TimeSlot } from "../types";

export const Home = () => {
  const {
    timeSlots,
    addTimeSlot,
    updateTimeSlot,
    addTimeSlots,
    handleDeleteTimeSlot,
    handleClearAllTimeSlots,
  } = useHangoutForm();

  const {
    isModalOpen,
    modalData,
    setModalData,
    editingId,
    handleOpenModal,
    handleOpenEditModal,
    handleCloseModal,
    handleSaveTimeSlot,
    isFormValid,
  } = useTimeSlotModal();

  const {
    isQuickModalOpen,
    quickData,
    setQuickData,
    handleOpenQuickModal,
    handleCloseQuickModal,
    handleGenerateQuickSlots,
    isQuickFormValid,
  } = useQuickSlotModal();

  const handleEditTimeSlot = (slot: TimeSlot) => {
    handleOpenEditModal(slot);
  };

  return (
    <div className="bg-success vh-100 overflow-auto">
      <div className="container">
        <div className="row">
          <div className="col-md-6 offset-md-3 mt-5">
            <div className="card card-body">
              <h1>Plan a hangout</h1>

              <div className="mt-3">
                <label htmlFor="name">Title (optional)</label>
                <input
                  id="name"
                  type="text"
                  className="form-control"
                  placeholder="Title of the hangout or activity"
                />
              </div>

              <Modal
                title={editingId ? "Edit time slot" : "Add time slot"}
                isVisible={isModalOpen}
                onClose={handleCloseModal}
                onConfirm={() => { handleSaveTimeSlot(addTimeSlot, updateTimeSlot); }}
                confirmText={editingId ? "Save" : "Add"}
                isConfirmDisabled={!isFormValid}>
                <TimeSlotForm
                  date={modalData.date}
                  startTime={modalData.startTime}
                  endTime={modalData.endTime}
                  onDateChange={(date) => {
                    setModalData(prev => ({ ...prev, date }));
                  }}
                  onStartTimeChange={(startTime) => {
                    setModalData(prev => ({ ...prev, startTime }));
                  }}
                  onEndTimeChange={(endTime) => {
                    setModalData(prev => ({ ...prev, endTime }));
                  }}
                />
              </Modal>

              <Modal
                title="Create quick slots"
                isVisible={isQuickModalOpen}
                onClose={handleCloseQuickModal}
                onConfirm={() => { handleGenerateQuickSlots(addTimeSlots); }}
                confirmText="Generate"
                isConfirmDisabled={!isQuickFormValid}>
                <QuickSlotGenerator
                  startDate={quickData.startDate}
                  endDate={quickData.endDate}
                  dailyStartTime={quickData.dailyStartTime}
                  dailyEndTime={quickData.dailyEndTime}
                  duration={quickData.duration}
                  isOverlapping={quickData.isOverlapping}
                  isWholeDay={quickData.isWholeDay}
                  onStartDateChange={(startDate) => {
                    setQuickData(prev => ({ ...prev, startDate }));
                  }}
                  onEndDateChange={(endDate) => {
                    setQuickData(prev => ({ ...prev, endDate }));
                  }}
                  onDailyStartTimeChange={(dailyStartTime) => {
                    setQuickData(prev => ({ ...prev, dailyStartTime }));
                  }}
                  onDailyEndTimeChange={(dailyEndTime) => {
                    setQuickData(prev => ({ ...prev, dailyEndTime }));
                  }}
                  onDurationChange={(duration) => {
                    setQuickData(prev => ({ ...prev, duration }));
                  }}
                  onOverlappingChange={(isOverlapping) => {
                    setQuickData(prev => ({ ...prev, isOverlapping }));
                  }}
                  onWholeDayChange={(isWholeDay) => {
                    setQuickData(prev => ({ ...prev, isWholeDay }));
                  }}
                />
              </Modal>

              <div className="mt-3">
                <button className="btn btn-primary w-100">
                  Create a new hangout
                </button>
              </div>

              <div className="mt-3">
                <Collapse title="Advanced options">
                  <div className="mt-2">
                    <div className="d-flex align-items-center justify-content-between">
                      <h4>Time slots</h4>
                      <div className="d-flex gap-2">
                        <button
                          onClick={handleOpenModal}
                          className="btn btn-sm btn-info">
                          <i className="ri-add-line" /> Add
                        </button>
                        <button
                          onClick={handleOpenQuickModal}
                          className="btn btn-sm btn-success">
                          <i className="ri-flashlight-line" /> Quick slots
                        </button>
                        <button
                          onClick={handleClearAllTimeSlots}
                          className="btn btn-sm btn-danger">
                          <i className="ri-delete-bin-line" /> Clear all
                        </button>
                      </div>
                    </div>
                    <div className="mt-3">
                      {timeSlots.length === 0 ? <div className="alert alert-warning">
                        By default if no time slots are defined, participants can choose any time for the whole upcoming week with one hour interval.
                      </div> : null}
                      <TimeSlotList
                        timeSlots={timeSlots}
                        onDelete={handleDeleteTimeSlot}
                        onEdit={handleEditTimeSlot}
                      />
                    </div>
                  </div>

                  <div className="mt-3">
                    <label htmlFor="password">Password required to join</label>
                    <input
                      id="password"
                      type="text"
                      className="form-control"
                      placeholder="Password required to join"
                    />
                  </div>

                  <div className="mt-3">
                    <label htmlFor="description">Description</label>
                    <input
                      id="description"
                      type="text"
                      className="form-control"
                      placeholder="Description"
                    />
                  </div>

                  <div className="mt-3">
                    <label htmlFor="location">Location used for weather</label>
                    <input
                      id="location"
                      type="text"
                      className="form-control"
                      placeholder="Location"
                    />
                  </div>

                  <div className="mt-3">
                    <label htmlFor="end-date">Accept responses until</label>
                    <input
                      id="end-date"
                      type="date"
                      className="form-control"
                      placeholder="End date"
                    />
                  </div>
                </Collapse>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
