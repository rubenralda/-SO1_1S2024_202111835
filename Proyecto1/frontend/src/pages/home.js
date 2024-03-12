import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import NavBar from "../components/navBar";
import Ram from "../components/componetRam";
import Cpu from "../components/componetCPU";

function Home() {
  return (
    <Container>
      <Row>
        <NavBar />
      </Row>
      <Row className="justify-content-md-center">
        <Col sm={6} md={6} lg={6} xl={6}>
          <Ram />
        </Col>
        <Col sm={6} md={6} lg={6} xl={6}>
          <Cpu />
        </Col>
      </Row>
    </Container>
  );
}

export default Home;
