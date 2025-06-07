// main.tsx
import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import DecryptPage from "./components/DecryptPage";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/decrypt" element={<DecryptPage />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);
