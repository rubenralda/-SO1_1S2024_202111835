import { Doughnut } from "react-chartjs-2";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { useState, useEffect } from "react";

ChartJS.register(ArcElement, Tooltip, Legend);

function App() {
  const [data, setData] = useState({
    labels: ["En Uso", "Libre"],
    datasets: [
      {
        label: "# of Votes",
        data: [63.24, 36.76],
        backgroundColor: ["rgba(255, 99, 132, 0.2)", "rgba(54, 162, 235, 0.2)"],
        borderColor: ["rgba(255, 99, 132, 1)", "rgba(54, 162, 235, 1)"],
        borderWidth: 1,
      },
    ],
  });
  const options = {
    responsive: true,
    legend: {
      display: false,
    },
    title: {
      display: true,
      text: "RAM Usage",
    },
  };

  const [porcentaje, setPorcentaje] = useState(0);

  function verRam() {
    fetch("/api/ram/actual")
      .then((res) => {
        return res.json();
      })
      .then((informacion) => {
        setPorcentaje(informacion.Porcentaje);
        setData({
          labels: ["En Uso", "Libre"],
          datasets: [
            {
              label: "# of Votes",
              data: [informacion.Porcentaje, 100 - informacion.Porcentaje],
              backgroundColor: [
                "rgba(255, 99, 132, 0.2)",
                "rgba(54, 162, 235, 0.2)",
              ],
              borderColor: ["rgba(255, 99, 132, 1)", "rgba(54, 162, 235, 1)"],
              borderWidth: 1,
            },
          ],
        });
      })
      .catch((err) => console.log(err));
  }
  useEffect(() => {
    const interval = setInterval(verRam, 1000);

    return () => clearInterval(interval);
  }, []);
  return (
    <div>
      <label style={{ fontSize: "40px" }}>RAM</label>
      <br />
      <label style={{ fontSize: "25px" }}>{porcentaje}% en uso</label>
      <Doughnut data={data} options={options} />
    </div>
  );
}

export default App;
