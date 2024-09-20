const amqp = require('amqplib/callback_api');

// Helper function to connect with retry logic
function connectRabbitMQ() {
  amqp.connect('amqp://rabbitmq', function (err, conn) {
    if (err) {
      console.log('Failed to connect to RabbitMQ. Retrying in 5 seconds...');
      setTimeout(connectRabbitMQ, 10000); // Retry after 5 seconds
      return;
    }

    // Connection successful, proceed with channel creation
    conn.createChannel(function (err, ch) {
      if (err) throw err;

      const requestQueue = 'scan_requests';
      const responseQueue = 'scan_responses';

      // Ensure the request queue exists
      ch.assertQueue(requestQueue, { durable: false });
      // Ensure the response queue exists
      ch.assertQueue(responseQueue, { durable: false });

      // Send a scan request every 5 seconds
      setInterval(() => {
        const scanRequest = JSON.stringify({ task: 'scan', file: 'file1.txt' });
        ch.sendToQueue(requestQueue, Buffer.from(scanRequest));
        console.log(" [x] Sent scan request: '%s'", scanRequest);
      }, 5000); // Send a message every 5 seconds

      // Consume response messages from scanner
      ch.consume(responseQueue, function (msg) {
        console.log(" [x] Received scan response: '%s'", msg.content.toString());
      }, { noAck: true });
    });
  });
}

// Start connection with retry mechanism
connectRabbitMQ();
