import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import NavBar from "../components/navBar";
import ComponetProcesos from "../components/componetProcesos";

function Procesos() {
  return (
    <Container>
      <Row>
        <NavBar />
      </Row>
      <Row className="justify-content-md-center">
        <Col sm={12} md={2} lg={2} xl={2}>
          <label>PRUEBA</label>
        </Col>
        <Col sm={12} md={10} lg={10} xl={10}>
          <ComponetProcesos />
        </Col>
      </Row>
    </Container>
  );
}

export default Procesos;
