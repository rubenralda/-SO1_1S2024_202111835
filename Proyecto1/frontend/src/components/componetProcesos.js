import React, { useEffect, useRef, useState } from "react";
import { Network } from "vis-network";

function Procesos({ desplegable, value }) {
  const container = useRef(null);

  const [nodes, setNodes] = useState([]);
  const [edges, setEdges] = useState([]);
  const [arrayProcesos, setArrayProcesos] = useState([]);

  const options = {
    layout: {
      hierarchical: {
        direction: "UD",
        sortMethod: "directed",
      },
    },
  };
  useEffect(() => {
    fetch("/api/procesos")
      .then((res) => res.json())
      .then((data) => {
        const nodes1 = [];
        const edges1 = [];
        const combobox = [];
        const resultado = JSON.parse(data.Procesos).processes;
        console.log(resultado)
        setArrayProcesos(resultado);
        console.log(arrayProcesos, "FETCH");
        nodes1.push({
          id: 0,
          label: "Padre",
        });
        /* nodes1.push({
          id: 3748,
          label: "Padre",
        }); */
        arrayProcesos.forEach((proceso) => {
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
        setNodes(nodes1);
        setEdges(edges1);
        desplegable(combobox);
      });
  }, []);

  /* useEffect(() => {
    if (value === 0) {
      console.log("No es", value);
      return;
    }
    if (value === 1) {
      console.log("Padre");
      fecthProcesos();
      return;
    }
    console.log("Cambio", value);
    const procesosBuscar = [value];
    const nodes1 = [];
    const edges1 = [];
    while (procesosBuscar.length > 0) {
      console.log(procesosBuscar);
      arrayProcesos.forEach((proceso) => {
        if (procesosBuscar[0] !== proceso.pid) {
          console.log("Entro", typeof procesosBuscar[0], typeof proceso.pid)
          return;
        }
        nodes1.push({
          id: proceso.pid,
          label: `${proceso.name} \nPID: ${proceso.pid.toString()}`,
        });
        edges1.push({ from: 1, to: proceso.pid });

        // Si hay hijos, agregar aristas
        proceso.child.forEach((child) => {
          console.log(child, "hijos")
          procesosBuscar.push(proceso.pid);
          edges1.push({ from: proceso.pid, to: child.pid });
        });
      });
      console.log(procesosBuscar.shift());
    }
    setNodes(nodes1);
    setEdges(edges1);
  }, [value]); */

  useEffect(() => {
    const network =
      container.current &&
      new Network(container.current, { nodes, edges }, options);
    console.log("generando arbol");
  }, [container, nodes, edges]);
  //console.log(nodes)
  /* const network =
      container.current &&
      new Network(container.current, { nodes, edges }, options); */

  return (
    <div ref={container} style={{ height: "500px", width: "800px" }}></div>
  );
}

export default Procesos;
