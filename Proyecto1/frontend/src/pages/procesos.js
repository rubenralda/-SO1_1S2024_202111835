import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import NavBar from "../components/navBar";
import ComponetProcesos from "../components/componetProcesos";
import Form from 'react-bootstrap/Form';
import { useState } from "react";

function Procesos() {
  const [Options, setOptions] = useState([]);
  const [value, setValor] = useState(0);

  const handleSelectChange = (event) => {
    setValor(event.target.value);
};
  return (
    <Container>
      <Row>
        <NavBar />
      </Row>
      <Row className="justify-content-md-center">
        <Col sm={12} md={2} lg={2} xl={2}>
          <Form.Select aria-label="Default select example" value={value} onChange={handleSelectChange}>
              <option value={""}>
                {"Escoge un proceso"}
              </option>
              {Options.map((proceso) => {
                return (
                  <option value={proceso.pid} key={proceso.pid}>
                    {proceso.name}
                  </option>
                );
              })}
          </Form.Select>
        </Col>
        <Col sm={12} md={10} lg={10} xl={10}>
          <ComponetProcesos desplegable={setOptions} value={value}/>
        </Col>
      </Row>
    </Container>
  );
}

export default Procesos;
