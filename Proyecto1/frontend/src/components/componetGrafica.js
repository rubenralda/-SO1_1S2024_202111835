import { useEffect, useState } from "react";

function SimulacionGrafica({ estado }) {
  const [estadoActual, setEstadoActual] = useState("");
  useEffect(() => {
    setEstadoActual(estado);
  }, [estado]);
  switch (estadoActual) {
    case "new":
      return <h1>Nuevo</h1>;
    case "stop":
      return <h1>stop</h1>;
    case "resume":
      return <h1>resume</h1>;
    default:
      return <h1>Escoge un proceso</h1>;
  }
}

export default SimulacionGrafica;
