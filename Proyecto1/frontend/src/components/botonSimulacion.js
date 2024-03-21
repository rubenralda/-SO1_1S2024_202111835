import Button from "react-bootstrap/Button";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import SimulacionGrafica from "./componetGrafica";
import { useState, useEffect } from "react";

function Botones({ setPid, pidSelect }) {
  const [estado, setEstado] = useState("");

  useEffect(() => {
    if (pidSelect === 0) {
      return;
    }
    fetch("/api/procesos/" + pidSelect.toString())
      .then((res) => res.json())
      .then((data) => {
        setEstado(data.Estado);
      })
      .catch((error) => {
        console.log(error);
      });
  }, [pidSelect]);

  function nuevo() {
    fetch("/api/start")
      .then((res) => res.json())
      .then((data) => {
        console.log(data.Mensaje);
        setPid(data.Pid);
        setEstado("new");
      })
      .catch((error) => {
        console.log(error);
      });
  }
  function stop() {
    fetch("/api/stop/" + pidSelect.toString())
      .then((res) => res.json())
      .then((data) => {
        console.log(data.Mensaje);
        setEstado("stop");
      })
      .catch((error) => {
        console.log(error);
      });
  }
  function ready() {
    fetch("/api/resume/" + pidSelect.toString())
      .then((res) => res.json())
      .then((data) => {
        console.log(data.Mensaje);
        setEstado("resume");
      })
      .catch((error) => {
        console.log(error);
      });
  }
  function kill() {
    fetch("/api/kill/" + pidSelect.toString())
      .then((res) => res.json())
      .then((data) => {
        console.log(data.Mensaje);
        setPid(0);
        setEstado("");
      })
      .catch((error) => {
        console.log(error);
      });
  }
  return (
    <>
      <Row className="justify-content-md-center">
        <Col sm={12} md={4} lg={4} xl={4}>
          <Button variant="success" onClick={nuevo}>
            New
          </Button>
        </Col>
        <Col sm={12} md={4} lg={4} xl={4}>
          <Button variant="warning" onClick={stop}>
            Stop
          </Button>
        </Col>
        <Col sm={12} md={4} lg={4} xl={4}>
          <Button variant="secondary" onClick={ready}>
            Ready
          </Button>
        </Col>
        <Col sm={12} md={4} lg={4} xl={4}>
          <Button variant="danger" onClick={kill}>
            Kill
          </Button>
        </Col>
      </Row>
      <SimulacionGrafica estado={estado} />
    </>
  );
}

export default Botones;
