import { createBrowserRouter, RouterProvider } from "react-router";

import { OutletContainer } from "./components/OutletContainer";
import { Calendar } from "./containers/Calendar";

const router = createBrowserRouter([
  {
    path: "/",
    element: <OutletContainer />,
    children: [
      {
        path: "/",
        element: <Calendar />,
      }
    ],
  },
]);

export function App() {
  return <RouterProvider router={router} />;
}
