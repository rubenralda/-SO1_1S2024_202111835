import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import NavBar from "../components/navBar";
import RamHistorico from "../components/RamHistorico";
import CPUHistorico from "../components/cpuHistorico";

function Historico() {
  return (
    <Container>
      <Row>
        <NavBar />
      </Row>
      <Row>
        <h1>Monitoreo Historico</h1>
      </Row>
      <Row className="justify-content-md-center">
        <Col sm={6} md={6} lg={6} xl={6}>
          <RamHistorico />
        </Col>
        <Col sm={6} md={6} lg={6} xl={6}>
          <CPUHistorico />
        </Col>
      </Row>
    </Container>
  );
}

export default Historico;
