const express = require("express"); //llamamos a Express
const morgan = require("morgan");
const cors = require("cors");
const mongoose = require('mongoose');
//const Image = require('../models/Image');

const app = express();

app.set("port", 5000);

app.use(morgan("dev"));
app.use(express.json());
app.use(cors());

mongoose.connect("mongodb://localhost:27017/imageDB", {
  useNewUrlParser: true,
  useUnifiedTopology: true,
});
// Define Image schema
const imageSchema = new mongoose.Schema({
  data: Buffer,
  contentType: String,
});
const ImageModel = mongoose.model("Image", imageSchema);

app.get("/", async (req, res) => {
  try {
    const imageData = req.body.image;
    const base64Data = imageData.replace(/^data:image\/png;base64,/, "");
    const bufferData = Buffer.from(base64Data, "base64");
    const image = new ImageModel({
      data: bufferData,
      contentType: "image/png",
    });
    await image.save();
    res.status(200).send("Image uploaded successfully");
  } catch (error) {
    console.error("Error uploading image:", error);
    res.status(500).send("Internal Server Error");
  }
});

// iniciamos nuestro servidor
app.listen(app.get("port"), () => {
  console.log("Servidor iniciado en el puerto: " + app.get("port"));
});
