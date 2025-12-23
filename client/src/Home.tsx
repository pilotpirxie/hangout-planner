import { Collapse } from "./components/Collapse";

export const Home = () => {
  return <div className="bg-success vh-100">
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

            <div className="mt-3">
              <button className="btn btn-primary w-100">
                Create a new hangout
              </button>
            </div>

            <div className="mt-3">
              <Collapse title="Advanced options">
                <div className="mt-2">
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
  </div>;
};