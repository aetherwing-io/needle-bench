import http from 'http';

export function jsonResponse(
  res: http.ServerResponse,
  status: number,
  data: any
): void {
  res.writeHead(status, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify(data));
}

export function parseBody(req: http.IncomingMessage): Promise<any> {
  return new Promise((resolve, reject) => {
    let body = '';
    req.on('data', (chunk: Buffer) => {
      body += chunk.toString();
    });
    req.on('end', () => {
      try {
        resolve(body ? JSON.parse(body) : {});
      } catch (err) {
        reject(new Error('Invalid JSON body'));
      }
    });
    req.on('error', reject);
  });
}

export function getTimestamp(): string {
  return new Date().toISOString();
}

export function randomId(): string {
  return Math.random().toString(36).substring(2, 10);
}
