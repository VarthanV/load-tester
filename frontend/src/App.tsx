import { createBrowserRouter, RouterProvider } from "react-router-dom";

const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <div>
        <Navbar />
        <HomePage/>
      </div>
    ),
  },
  {
    path: "/test",
    element: (
      <>
        <Navbar />
        <CreateTestRequest />{" "}
      </>
    ),
  },
  {
    path: "/test/:id",
    element: (
      <>
        <Navbar />
        <MetricsPage />{" "}
      </>
    ),
  },
]);

import Navbar from "./components/Navbar";
import CreateTestRequest from "./pages/CreateTestRequest";
import { Toaster } from "react-hot-toast";
import MetricsPage from "./pages/MetricsPage";
import HomePage from "./pages/HomePage";

function App() {
  return (
    <>
      <Toaster />
      <RouterProvider router={router} />
    </>
  );
}

export default App;
