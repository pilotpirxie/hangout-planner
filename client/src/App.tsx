import { createBrowserRouter, RouterProvider } from "react-router";

import { StrictMode } from "react";
import { Provider } from "react-redux";
import { OutletContainer } from "./components/OutletContainer";
import { Calendar } from "./containers/Calendar";
import { Home } from "./containers/Home";
import { store } from "./data/store";

const router = createBrowserRouter([
  {
    path: "/",
    element: <OutletContainer />,
    children: [
      {
        path: "/",
        element: <Home />,
      },
      {
        path: "/calendar",
        element: <Calendar />,
      }
    ],
  },
]);

export function App() {
  return <StrictMode>
    <Provider store={store}>
      <RouterProvider router={router} />
    </Provider>
  </StrictMode>;
}
