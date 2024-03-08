const redis = require('redis');
const amqp = require('amqplib/callback_api');

const REDIS_HOST = '127.0.0.1';
const REDIS_PORT = 6379;
const QUEUE_NAME = 'counterQueue';

let redisClient = "";
let initialize=false;

// IIFE to run redis
(async function (){
    try{
        if(!initialize){

            redisClient=redis.createClient({ host: REDIS_HOST, port: REDIS_PORT });
            await redisClient.connect();
            console.log("connceted to redis")
        }
    }
    catch(e){
        console.log("error in redis connection",e);
        throw e;
    }
})(); 

const LOCK_KEY = 'counterLock';
const LOCK_EXPIRE = 1000; // Lock expires in 1000 ms

function acquireLock() {
    return new Promise(async (resolve) => {
        try {
            const result = await redisClient.set(LOCK_KEY, 'locked', {
                NX: true,
                PX: LOCK_EXPIRE,
            });
            resolve(result === 'OK');
        } catch (error) {
            console.error('Error acquiring lock:', error);
            resolve(false);
        }
    });
}

async function releaseLock() {
    await redisClient.del(LOCK_KEY); // Delete lock key
}
// creating two counter with one locked
const incrementCounter2 = async (channel) => {
    try {
        const lockAcquired = await acquireLock();
        console.log(lockAcquired);
        if (!lockAcquired) {
            console.log('Lock not acquired');
            return;
        }
        const counter = await redisClient.get('counter') || '0';
        const newCounter = parseInt(counter) + 1;
        await redisClient.set('counter', newCounter);
        const newCounterValue = newCounter.toString(); // Convert number to string
        const buffer = Buffer.from(newCounterValue); // Convert string to Buffer

        channel.publish("", QUEUE_NAME, buffer);
        console.log(`Counter updated to: ${newCounter}`);
        // await releaseLock();
    } catch (err) {
        console.error('Error processing counter:', err);
        await releaseLock();
    }
};
const incrementCounter = async (channel) => {
    try {
        const lockAcquired = await acquireLock();
        console.log(lockAcquired);
        if (!lockAcquired) {
            console.log('Lock not acquired');
            return;
        }
        const counter = await redisClient.get('counter') || '0';
        const newCounter = parseInt(counter) + 1;
        await redisClient.set('counter', newCounter);
        const newCounterValue = newCounter.toString(); // Convert number to string
        const buffer = Buffer.from(newCounterValue); // Convert string to Buffer

        channel.publish("", QUEUE_NAME, buffer);
        console.log(`Counter updated to: ${newCounter}`);
        await releaseLock();
    } catch (err) {
        console.error('Error processing counter:', err);
        await releaseLock();
    }
};
// Connect to RabbitMQ server

amqp.connect('amqp://localhost', async (error0, connection) => {
    if (error0) {
        console.error('Failed to connect to RabbitMQ:', error0);
        return;
    }

    const channel = connection.createChannel();
    channel.assertQueue(QUEUE_NAME, { durable: false });

    await incrementCounter2(channel);
    await incrementCounter(channel);
        channel.consume(QUEUE_NAME, async (msg) => {
            if (msg && msg.content) {
                console.log(`Received message from queue. ${msg.content.toString()}`);
                channel.ack(msg);
            }
        }, { noAck: false });
});

// Error handling
process.on('uncaughtException', (err) => {
    console.error('Unhandled Exception', err);
});
process.on('unhandledRejection', (err) => {
    console.error('Unhandled Rejection', err);
});
