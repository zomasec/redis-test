const amqp = require('amqplib/callback_api');

// Connect to RabbitMQ
amqp.connect('amqp://rabbitmq', function (err, conn) {
  if (err) throw err;

  // Create a channel
  conn.createChannel(function (err, ch) {
    if (err) throw err;

    const queue = 'scan_requests';
    const msg = JSON.stringify({ task: 'scan', file: 'file1.txt' });

    // Ensure the queue exists
    ch.assertQueue(queue, { durable: false });

    // Send a message to the queue
    ch.sendToQueue(queue, Buffer.from(msg));
    console.log(" [x] Sent '%s'", msg);
  });

  // Close the connection after a short delay
  setTimeout(function () {
    conn.close();
    process.exit(0);
  }, 500);
});
