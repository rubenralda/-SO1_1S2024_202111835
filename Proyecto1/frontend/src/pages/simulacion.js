import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import NavBar from "../components/navBar";
import Form from "react-bootstrap/Form";
import { useEffect, useState } from "react";
import Botones from "../components/botonSimulacion";

function Simulacion() {
  const [Options, setOptions] = useState([]);
  const [pidSelect, setPidSelect] = useState(0);

  const handleSelectChange = (event) => {
    setPidSelect(event.target.value);
  };

  useEffect(() => {
    fetch("/api/procesos/simulados")
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
        setOptions(data.Procesos);
      });
  }, [pidSelect]);

  return (
    <Container>
      <Row>
        <NavBar />
      </Row>
      <Row className="justify-content-md-center">
        <Col sm={12} md={2} lg={2} xl={2}>
          <Form.Select value={pidSelect} onChange={handleSelectChange}>
            <option value={0}>{"Escoge un proceso"}</option>
            {Options.map((proceso) => {
              return (
                <option value={proceso.Pid} key={proceso.Id}>
                  {proceso.Pid.toString()}
                </option>
              );
            })}
          </Form.Select>
          <Botones setPid={setPidSelect} pidSelect={pidSelect} />
        </Col>
        <Col sm={12} md={10} lg={10} xl={10}></Col>
      </Row>
    </Container>
  );
}

export default Simulacion;
