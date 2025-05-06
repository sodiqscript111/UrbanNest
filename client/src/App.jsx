import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './pages/Home.jsx';
import EditProperty from './pages/EditProperty.jsx';
import UploadFile from "./pages/UploadFile.jsx";
import Welcome from './pages/Welcome.jsx';
import ProjectDetails from "./pages/ProjectDetails.jsx";
import ErrorBoundary from "./component/ErrorBoundary.jsx";

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Welcome />} />
                <Route path="/edit/:id" element={<EditProperty />} />
                <Route path="/upload" element={<UploadFile />} />
                <Route path="/home" element={<Home />} />
                <Route
                    path="/listing/:id"
                    element={
                        <ErrorBoundary>
                            <ProjectDetails />
                        </ErrorBoundary>
                    }
                />
            </Routes>
        </BrowserRouter>
    );
}

export default App;