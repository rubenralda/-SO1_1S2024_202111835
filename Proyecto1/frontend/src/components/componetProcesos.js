import React, { useEffect, useRef, useState } from "react";
import { Network } from "vis-network";

function Procesos({ desplegable, value }) {
  const container = useRef(null);
  const [nodes, setNodes] = useState([]);
  const [edges, setEdges] = useState([]);
  const [arrayProcesos, setArrayProcesos] = useState([]);

  useEffect(() => {
    const nodes1 = [];
    const edges1 = [];
    const combobox = [];
    fetch("/api/procesos")
      .then((res) => res.json())
      .then((data) => {
        const resultado = JSON.parse(data.Procesos).processes;
        console.log(resultado);
        nodes1.push({
          id: 0,
          label: "Padre",
        });
        resultado.forEach((proceso) => {
          nodes1.push({
            id: proceso.pid,
            label: `${proceso.name} \nPID: ${proceso.pid.toString()}`,
          });
          combobox.push({ pid: proceso.pid, name: proceso.name });
          edges1.push({ from: proceso.parent, to: proceso.pid });
          // Si hay hijos, agregar aristas
          /* proceso.child.forEach((child) => {
            edges1.push({ from: proceso.pid, to: child.pid });
          }); */
        });
        setArrayProcesos(resultado);
        setNodes(nodes1);
        setEdges(edges1);
        desplegable(combobox);
      });
  }, [desplegable]);

  useEffect(() => {
    if (value === "") {
      console.log("No esta seleccionado algo", value);
      return;
    }
    console.log("Graficando");
    const procesosBuscar = [value];
    const nodes1 = [];
    const edges1 = [];
    while (procesosBuscar.length > 0) {
      arrayProcesos.forEach((proceso) => {
        if (procesosBuscar[0] !== proceso.pid.toString()) {
          return;
        }
        nodes1.push({
          id: proceso.pid,
          label: `${proceso.name} \nPID: ${proceso.pid.toString()}`,
        });
        edges1.push({ from: proceso.parent, to: proceso.pid });

        // Si hay hijos, agregar aristas
        proceso.child.forEach((child) => {
          procesosBuscar.push(child.pid.toString());
        });
      });
      procesosBuscar.shift()
    }
    setNodes(nodes1);
    setEdges(edges1);
  }, [value, arrayProcesos]);

  useEffect(() => {
    container.current &&
    new Network(
      container.current,
      { nodes, edges },
      {
        layout: {
          hierarchical: {
            direction: "UD",
            sortMethod: "directed",
          },
        },
      }
    );
  console.log("generando arbol");
  }, [nodes, edges, arrayProcesos])

  return (
    <div ref={container} style={{ height: "500px", width: "800px" }}></div>
  );
}

export default Procesos;
