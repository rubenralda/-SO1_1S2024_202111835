import { Routes, Route } from "react-router-dom";
import Home from "./pages/home";
import Historico from "./pages/historico"
import Procesos from "./pages/procesos";

function Aplicacion() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="tiempo-real" element={<Home />} />
      <Route path="historico" element={<Historico />} />
      <Route path="procesos" element={<Procesos/>} />
    </Routes>
  );
}

export default Aplicacion;
