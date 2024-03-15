import React, { useEffect, useRef, useState } from "react";
import { Network } from "vis-network";

function Procesos() {
    const container = useRef(null);

    const [nodes, setNodes] = useState([]);
    const [edges, setEdges] = useState([]);
  
    const options = {
        layout: {
            hierarchical: {
                direction: "UD",
                sortMethod: "directed"
            }
        }
    };

    useEffect(() => {
        fetch("/api/procesos")
          .then((res) => res.json())
          .then((data) => {
              const nodes1 = [];
              const edges1 = [];
              const arrayProcesos = JSON.parse(data.Procesos).processes;
              arrayProcesos.forEach((proceso) => {
                 nodes1.push({ id: proceso.pid, label: `${proceso.name} \nPID: ${proceso.pid.toString()}`});
                 edges1.push({ from: 1, to: proceso.pid });
                  proceso.child.forEach((child) => {
                    if (proceso.pid != 3816) {
                          return
                    }
                    edges1.push({ from: proceso.pid, to: child.pid });  
                  });
              });
              console.log(edges1);
              setNodes(nodes1);
              setEdges(edges1);
          });
    }, []);
    
    useEffect(() => {
      const network =
        container.current &&
        new Network(container.current, { nodes, edges }, options);
        console.log("generando arbol");
    }, [container, nodes, edges]);
  
    return <div ref={container} style={{ height: '500px', width: '800px' }} />;
}

export default Procesos;
