const amqp = require("amqplib");
const { Pool } = require("pg");
require("dotenv").config();

// PostgreSQL Database Connection
const pool = new Pool({
  user: process.env.DB_USER,
  host: process.env.DB_HOST,
  database: process.env.DB_NAME,
  password: process.env.DB_PASS,
  port: process.env.DB_PORT,
});

async function startConsumer() {
  try {
    const connection = await amqp.connect(process.env.RABBITMQ_URL || "amqp://localhost");
    console.log("✅ Connected to RabbitMQ!");
    const channel = await connection.createChannel();
    const queue = "login_events";

    await channel.assertQueue(queue, { durable: false });
    console.log("Waiting for login requests...");
    console.log(`🎧 Waiting for messages in queue: ${queue}`);

    channel.consume(queue, async (msg) => {
      const loginData = JSON.parse(msg.content.toString());
      console.log("Processing login for:", loginData);
        // ✅ Print the received queue data
        console.log("\n🔹 Received Login Data from Queue:");
        console.log(loginData); // Print entire object
        console.log(`🔸 Username: ${loginData.username}`);
        console.log(`🔸 Password: ${loginData.password}`);

      try {
        const client = await pool.connect();
        const result = await client.query(
          "SELECT * FROM user_login WHERE email = $1",
          [loginData.username]
        );

        if (result.rowCount === 0) {
          sendResponse(channel, msg, { status: "error", message: "User not found" });
        } else {
          const user = result.rows[0];

          // Check Password
          const bcrypt = require("bcrypt");
          const match = await bcrypt.compare(loginData.password, user.password);

          if (match) {
            sendResponse(channel, msg, {
              status: "success",
              message: "Login successful",
              token: "mock_jwt_token"
            });
          } else {
            sendResponse(channel, msg, { status: "error", message: "Invalid credentials" });
          }
        }

        client.release();
      } catch (error) {
        console.error("Database error:", error);
        sendResponse(channel, msg, { status: "error", message: "Database error" });
      }

      channel.ack(msg);
    });
  } catch (error) {
    console.error("RabbitMQ Consumer error:", error);
  }
}

function sendResponse(channel, msg, response) {
  channel.sendToQueue(
    msg.properties.replyTo,
    Buffer.from(JSON.stringify(response)),
    { correlationId: msg.properties.correlationId }
  );
}

startConsumer();
