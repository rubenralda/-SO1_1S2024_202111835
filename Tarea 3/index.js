const Redis = require("ioredis");
const redis = new Redis({
  host: "10.167.120.43",
  port: 6379,
  connectTimeout: 5000,
});

setInterval(() => {
  const message = { msg: "Hola a todos" };

  // Message can be either a string or a buffer
  redis
    .publish("test", JSON.stringify(message))
    .then(() => console.log("Mensaje publicado"))
    .catch((err) => console.log("error al publicar mensaje", err));
}, 5000);
