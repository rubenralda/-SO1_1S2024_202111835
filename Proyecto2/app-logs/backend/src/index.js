const express = require("express");
const morgan = require("morgan");
const cors = require("cors");


// creacion de la app o mi API
const app = express();

//configuraciones
app.set("port", 3000);

// usando morgan para middlewares
app.use(morgan("dev")); // para poder visualizar los estados de nuestro servidor
app.use(express.json()); // para poder manejar los json
app.use(cors());
app.use(require("./router"))

// inicializando mi servidor
app.listen(app.get("port"), () => {
  console.log("Servidor iniciado en el puerto: " + app.get("port"));
});
