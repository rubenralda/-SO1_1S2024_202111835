const { Router } = require("express");
const router = Router();
const { MongoClient } = require("mongodb");

// Replace the uri string with your connection string.
const uri = "mongodb://admin:1234@34.66.138.214:27017/?authSource=admin";

function run() {
  return new Promise((resolve, reject) => {
    MongoClient.connect(uri)
      .then((db) => {
        const database = db.db("proyecto2");
        database
          .collection("logs")
          .find({})
          .sort({ fecha: -1 })
          .limit(20)
          .toArray()
          .then((result) => {
            resolve(result);
            database.close();
          })
          .catch((err) => reject(err));
      })
      .catch((err) => reject(err));
  });
}

router.get("/api/logs", async (req, res) => {
  console.log("Obteniendo datos...");
  run()
    .then((data) => res.send(data))
    .catch((err) => {
      console.log(err);
      res.status(404).send(err);
    });
});

module.exports = router;
