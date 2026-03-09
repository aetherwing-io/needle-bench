const http = require('http');
const EventEmitter = require('events');
const { Monitor } = require('./monitor');
const { DataProcessor } = require('./processor');
const { TaskQueue } = require('./queue');

const PORT = process.env.PORT || 8080;

// Global event bus for inter-component communication
const eventBus = new EventEmitter();
eventBus.setMaxListeners(0); // Disable warnings — we "know" what we're doing

const monitor = new Monitor(eventBus);
const processor = new DataProcessor(eventBus);
const taskQueue = new TaskQueue(eventBus);

const server = http.createServer((req, res) => {
  if (req.method === 'POST' && req.url === '/process') {
    let body = '';
    req.on('data', chunk => { body += chunk; });
    req.on('end', () => {
      try {
        const data = JSON.parse(body);
        const result = processor.process(data);
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(result));
      } catch (err) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: err.message }));
      }
    });
  } else if (req.method === 'GET' && req.url === '/stats') {
    const stats = monitor.getStats();
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify(stats));
  } else if (req.method === 'POST' && req.url === '/task') {
    let body = '';
    req.on('data', chunk => { body += chunk; });
    req.on('end', () => {
      try {
        const task = JSON.parse(body);
        taskQueue.enqueue(task);
        res.writeHead(202, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ status: 'queued' }));
      } catch (err) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: err.message }));
      }
    });
  } else if (req.method === 'GET' && req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ status: 'ok', uptime: process.uptime() }));
  } else {
    res.writeHead(404);
    res.end('Not Found');
  }
});

server.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}`);
});

module.exports = { server, eventBus };
