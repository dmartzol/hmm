import "./App.css";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from "./components/Login/Login";
import PageLayout from "./components/PageLayout/PageLayout";

function App() {
  return (
    <BrowserRouter>
      <PageLayout>
        <Routes>
          <Route path="/" element={<Login />} />
        </Routes>
      </PageLayout>
    </BrowserRouter>
  );
}

export default App;
