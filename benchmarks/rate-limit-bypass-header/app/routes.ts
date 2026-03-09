import http from 'http';
import { parseBody, jsonResponse } from './utils';

const DATA_STORE: Record<string, any> = {
  '1': { id: '1', name: 'Alice', email: 'alice@example.com', balance: 1000 },
  '2': { id: '2', name: 'Bob', email: 'bob@example.com', balance: 2500 },
  '3': { id: '3', name: 'Carol', email: 'carol@example.com', balance: 750 },
};

export async function router(req: http.IncomingMessage, res: http.ServerResponse) {
  const url = new URL(req.url || '/', `http://localhost`);
  const path = url.pathname;

  if (req.method === 'GET' && path === '/health') {
    jsonResponse(res, 200, { status: 'ok' });
    return;
  }

  if (req.method === 'GET' && path === '/api/users') {
    const users = Object.values(DATA_STORE).map(u => ({
      id: u.id,
      name: u.name,
      email: u.email,
    }));
    jsonResponse(res, 200, { users });
    return;
  }

  if (req.method === 'GET' && path.startsWith('/api/users/')) {
    const id = path.split('/').pop();
    const user = id ? DATA_STORE[id] : null;
    if (!user) {
      jsonResponse(res, 404, { error: 'user not found' });
      return;
    }
    jsonResponse(res, 200, user);
    return;
  }

  if (req.method === 'POST' && path === '/api/transfer') {
    const body = await parseBody(req);
    const { from, to, amount } = body;

    const sender = DATA_STORE[from];
    const receiver = DATA_STORE[to];

    if (!sender || !receiver) {
      jsonResponse(res, 404, { error: 'user not found' });
      return;
    }

    if (sender.balance < amount) {
      jsonResponse(res, 400, { error: 'insufficient funds' });
      return;
    }

    sender.balance -= amount;
    receiver.balance += amount;

    jsonResponse(res, 200, {
      success: true,
      from: { id: sender.id, balance: sender.balance },
      to: { id: receiver.id, balance: receiver.balance },
    });
    return;
  }

  if (req.method === 'GET' && path === '/api/stats') {
    const { limiter } = require('./server');
    jsonResponse(res, 200, limiter.getStats());
    return;
  }

  jsonResponse(res, 404, { error: 'not found' });
}
