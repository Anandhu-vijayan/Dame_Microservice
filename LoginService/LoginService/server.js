const express = require("express");
const cors = require("cors");
const sendMessageToQueue = require("./rabbitmqService");
require("dotenv").config();

const app = express();

// Enable CORS
app.use(cors({
  origin: "http://localhost:3000",
  methods: "GET,POST,PUT,DELETE",
  allowedHeaders: "Content-Type,Authorization"
}));

app.use(express.json());

app.post("/api/login", async (req, res) => {
  const { username, password } = req.body;

  try {
    const loginEvent = { username, password };
    console.log( username, password);
    // Send message to RabbitMQ and wait for response
    const response = await sendMessageToQueue(loginEvent);

    // Send the response from the consumer to the client
    if (response.status === "success") {
      return res.json({ message: response.message, token: response.token });
    } else {
      return res.status(401).json({ message: response.message });
    }
  } catch (error) {
    console.error("Login error:", error);
    return res.status(500).json({ message: "Internal Server Error" });
  }
});

const PORT = process.env.PORT || 5000;
app.listen(PORT, () => console.log(`Login Service running on port ${PORT}`));
