require("dotenv").config();
const amqp = require("amqplib");

const RABBITMQ_URL = process.env.RABBITMQ_URL || "amqp://localhost";

async function sendMessageToQueue(msg) {
  try {
    const connection = await amqp.connect(RABBITMQ_URL);
    const channel = await connection.createChannel();

    const queue = "login_events";
    const replyQueue = await channel.assertQueue("", { exclusive: true });

    const correlationId = generateUuid();

    return new Promise((resolve, reject) => {
      channel.consume(
        replyQueue.queue,
        (msg) => {
          if (msg.properties.correlationId === correlationId) {
            resolve(JSON.parse(msg.content.toString()));
            setTimeout(() => connection.close(), 500);
          }
        },
        { noAck: true }
      );

      channel.sendToQueue(queue, Buffer.from(JSON.stringify(msg)), {
        correlationId,
        replyTo: replyQueue.queue,
      });
    });
  } catch (error) {
    console.error("RabbitMQ error:", error);
    throw error;
  }
}

function generateUuid() {
  return Math.random().toString() + Math.random().toString();
}

module.exports = sendMessageToQueue;
