import { Line } from "react-chartjs-2";
import { useState, useEffect } from "react";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

function RamHistorico() {
  const [data, setData] = useState({
    labels: [],
    datasets: [
      {
        label: "",
        data: [],
        borderColor: "rgb(255, 99, 132)",
        backgroundColor: "rgba(255, 99, 132, 0.5)",
      },
    ],
  });
  const options = {
    responsive: true,
    plugins: {
      legend: {
        display: false,
      },
      title: {
        display: true,
        text: "Rendimiento a lo largo del tiempo",
      },
    },
  };

  function ramHistorico() {
    fetch("/api/ram")
      .then((res) => {
        return res.json();
      })
      .then((data) => {
        //console.log(data)
        setData({
          labels: data.Labels,
          datasets: [
            {
              label: "Porcentaje",
              data: data.Data,
              borderColor: "rgb(255, 99, 132)",
              backgroundColor: "rgba(255, 99, 132, 0.5)",
            },
          ],
        });
      })
      .catch((err) => console.log(err));
  }
  useEffect(() => {
    const interval = setInterval(ramHistorico, 1000);
    return () => clearInterval(interval);
  }, []);

  return (
    <>
      <label style={{ fontSize: "40px" }}>Memoria RAM</label>
      <Line data={data} options={options} />
    </>
  );
}

export default RamHistorico;
