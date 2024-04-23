import { Routes, Route } from "react-router-dom";
import Home from "./pages/home";
import Historico from "./pages/historico"
import Procesos from "./pages/procesos";
import Simulacion from "./pages/simulacion";

function Aplicacion() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="tiempo-real" element={<Home />} />
      <Route path="historico" element={<Historico />} />
      <Route path="procesos" element={<Procesos/>} />
      <Route path="simulacion" element={<Simulacion/>} />
    </Routes>
  );
}

export default Aplicacion;
