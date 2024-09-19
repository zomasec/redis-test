const redis = require('redis');

const client = redis.createClient({
    host: 'redis', // Use the service name instead of localhost
    port: 6379
});

client.on('error', (err) => {
    console.error('Redis Client Error', err);
});

async function publishMessage() {
    await client.connect();

    const data = {
        message: "Hello from Node.js",
        timestamp: new Date().toISOString()
    };

    const jsonData = JSON.stringify(data);
    await client.publish('mychannel', jsonData);

    console.log('Published:', jsonData);
    await client.quit();
}

publishMessage().catch(console.error);
