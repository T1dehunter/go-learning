import WebSocket from 'ws';

const ws = new WebSocket('ws://localhost:3000/chat');

ws.on('error', console.error);

ws.on('open', function open() {
    ws.send('message from node client');
});

ws.on('message', function message(data) {
    console.log('received: %s', data);
});