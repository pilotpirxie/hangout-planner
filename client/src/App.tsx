import { createBrowserRouter, RouterProvider } from "react-router";

import { OutletContainer } from "./components/OutletContainer";
import { Home } from "./Home";

const router = createBrowserRouter([
  {
    path: "/",
    element: <OutletContainer />,
    children: [
      {
        path: "/",
        element: <Home />,
      }
    ],
  },
]);

export function App() {
  return <RouterProvider router={router} />;
}
