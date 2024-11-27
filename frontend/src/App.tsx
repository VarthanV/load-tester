
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

const router = createBrowserRouter([
  {
    path: "/",
    element: <div>
      <Navbar/>
    </div>,
  },
  {
    path:"/test",
    element:<CreateTestRequest/>,
  }
]);



import Navbar from './components/Navbar';
import CreateTestRequest from './pages/CreateTestRequest';

function App() {
  return (
    <>
    <RouterProvider  router={router} />
    </>
  )
}

export default App
