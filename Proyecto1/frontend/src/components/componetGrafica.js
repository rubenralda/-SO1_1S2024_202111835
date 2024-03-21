import { useEffect, useState, useRef } from "react";
import { Network } from "vis-network";

function SimulacionGrafica({ estado }) {
  const container = useRef(null);
  const [nodes, setNodes] = useState([]);
  const [edges, setEdges] = useState([]);

  useEffect(() => {
    let nodesN = [
      { id: 0, label: "New" },
      { id: 1, label: "Stop" },
      { id: 2, label: "Running" },
    ];
    let edgesN = [];
    switch (estado) {
      case "new":
        edgesN = [
          { from: 0, to: 1 },
          { from: 1, to: 2 },
        ];
        break;
      case "stop":
        edgesN = [
          { from: 0, to: 1 },
          { from: 1, to: 2 },
          { from: 2, to: 1, arrows: {to: true}, color: "red" },
        ];
        break;
      case "resume":
        edgesN = [
          { from: 0, to: 1 },
          { from: 1, to: 2, arrows: {to: true}, color: "red" },
          { from: 2, to: 1 },
        ];
        break;
      default:
        nodesN = [];
        break;
    }
    setNodes(nodesN);
    setEdges(edgesN);
  }, [estado]);

  useEffect(() => {
    container.current &&
      new Network(
        container.current,
        { nodes, edges },
        {
          layout: {},
        }
      );
    console.log("generando arbol");
  }, [nodes, edges]);

  return (
    <div ref={container} style={{ height: "500px", width: "800px" }}></div>
  );
}

export default SimulacionGrafica;
