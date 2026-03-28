import dayjs from "dayjs";
import { useState } from "react";
import { useNavigate } from "react-router";
import { Collapse } from "../components/Collapse";
import { Modal } from "../components/Modal";
import { QuickSlotGenerator } from "../components/QuickSlotGenerator";
import { TimeSlotForm } from "../components/TimeSlotForm";
import { TimeSlotList } from "../components/TimeSlotList";
import { useCreateCalendarMutation, useCreateCalendarTimeSlotsMutation } from "../data/calendarsApi";
import { useHangoutForm } from "../hooks/useHangoutForm";
import { useQuickSlotModal } from "../hooks/useQuickSlotModal";
import { useTimeSlotModal } from "../hooks/useTimeSlotModal";

export const Home = () => {
  const [title, setTitle] = useState("");
  const [password, setPassword] = useState("");
  const [description, setDescription] = useState("");
  const [location, setLocation] = useState("");
  const [acceptResponsesUntil, setAcceptResponsesUntil] = useState("");

  const navigate = useNavigate();

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
    isQuickSlotModalOpen,
    quickSlotModalData,
    setQuickSlotModalData,
    handleOpenQuickSlotModal,
    handleCloseQuickSlotModal,
    handleGenerateQuickSlots,
    isQuickSlotFormValid,
  } = useQuickSlotModal();

  const [createCalendar, { isLoading: isCreatingCalendar }] = useCreateCalendarMutation();
  const [createCalendarTimeSlots, { isLoading: isCreatingTimeSlots }] = useCreateCalendarTimeSlotsMutation();
  const handleCreateCalendar = async () => {
    try {
      if (password.length > 0 && password.length < 3) {
        alert("Password must be optional or at least 3 characters long.");
        return;
      }

      if (title.length > 0 && (title.length < 3 || title.length > 256)) {
        alert("Title must be between 3 and 256 characters long.");
        return;
      }

      if (description.length > 1024) {
        alert("Description cannot exceed 1024 characters.");
        return;
      }

      if (location.length > 512) {
        alert("Location cannot exceed 512 characters.");
        return;
      }

      if (acceptResponsesUntil && isNaN(Date.parse(acceptResponsesUntil))) {
        alert("Accept responses until must be a valid date or blank.");
        return;
      }

      const calendar = await createCalendar({
        title: title,
        description,
        location,
        accept_responses_until: acceptResponsesUntil,
        password
      }).unwrap();

      await createCalendarTimeSlots({
        calendar_id: calendar.id,
        admin_token: calendar.admin_token,
        time_slots: timeSlots.map(slot => ({
          start_date: dayjs(slot.startDate).toISOString(),
          end_date: dayjs(slot.endDate).toISOString(),
        })),
      }).unwrap();

      void navigate(`/calendar/${calendar.id}`);
    } catch (err) {
      console.error("Failed to create calendar:", err);
    }
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
                  value={title}
                  minLength={3}
                  maxLength={256}
                  onChange={(e) => { setTitle(e.target.value); }}
                />
              </div>

              <Modal
                title={editingId ? "Edit time slot" : "Add time slot"}
                isVisible={isModalOpen}
                onClose={handleCloseModal}
                onConfirm={() => { handleSaveTimeSlot(editingId ? updateTimeSlot : addTimeSlot); }}
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
                isVisible={isQuickSlotModalOpen}
                onClose={handleCloseQuickSlotModal}
                onConfirm={() => { handleGenerateQuickSlots(addTimeSlots); }}
                confirmText="Generate"
                isConfirmDisabled={!isQuickSlotFormValid}>
                <QuickSlotGenerator
                  startDate={quickSlotModalData.startDate}
                  endDate={quickSlotModalData.endDate}
                  dailyStartTime={quickSlotModalData.dailyStartTime}
                  dailyEndTime={quickSlotModalData.dailyEndTime}
                  duration={quickSlotModalData.duration}
                  isOverlapping={quickSlotModalData.isOverlapping}
                  isWholeDay={quickSlotModalData.isWholeDay}
                  onStartDateChange={(startDate) => {
                    setQuickSlotModalData(prev => ({ ...prev, startDate }));
                  }}
                  onEndDateChange={(endDate) => {
                    setQuickSlotModalData(prev => ({ ...prev, endDate }));
                  }}
                  onDailyStartTimeChange={(dailyStartTime) => {
                    setQuickSlotModalData(prev => ({ ...prev, dailyStartTime }));
                  }}
                  onDailyEndTimeChange={(dailyEndTime) => {
                    setQuickSlotModalData(prev => ({ ...prev, dailyEndTime }));
                  }}
                  onDurationChange={(duration) => {
                    setQuickSlotModalData(prev => ({ ...prev, duration }));
                  }}
                  onOverlappingChange={(isOverlapping) => {
                    setQuickSlotModalData(prev => ({ ...prev, isOverlapping }));
                  }}
                  onWholeDayChange={(isWholeDay) => {
                    setQuickSlotModalData(prev => ({ ...prev, isWholeDay }));
                  }}
                />
              </Modal>

              <div className="mt-3">
                <button
                  disabled={isCreatingCalendar || isCreatingTimeSlots}
                  className="btn btn-primary w-100"
                  onClick={() => { void handleCreateCalendar(); }}>
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
                          onClick={handleOpenQuickSlotModal}
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
                        onEdit={handleOpenEditModal}
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
                      value={password}
                      minLength={3}
                      maxLength={128}
                      onChange={(e) => { setPassword(e.target.value); }}
                    />
                  </div>

                  <div className="mt-3">
                    <label htmlFor="description">Description</label>
                    <input
                      id="description"
                      type="text"
                      className="form-control"
                      placeholder="Description"
                      value={description}
                      maxLength={1024}
                      onChange={(e) => { setDescription(e.target.value); }}
                    />
                  </div>

                  <div className="mt-3">
                    <label htmlFor="location">Location used for weather</label>
                    <input
                      id="location"
                      type="text"
                      className="form-control"
                      placeholder="Location"
                      value={location}
                      maxLength={512}
                      onChange={(e) => { setLocation(e.target.value); }}
                    />
                  </div>

                  <div className="mt-3">
                    <label htmlFor="end-date">Accept responses until</label>
                    <input
                      id="end-date"
                      type="date"
                      className="form-control"
                      placeholder="End date"
                      value={acceptResponsesUntil}
                      onChange={(e) => { setAcceptResponsesUntil(e.target.value); }}
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
