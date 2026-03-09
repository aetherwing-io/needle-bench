import http from 'http';

export function parseBody(req: http.IncomingMessage): Promise<Record<string, unknown>> {
  return new Promise((resolve, reject) => {
    const chunks: Buffer[] = [];
    req.on('data', (chunk: Buffer) => chunks.push(chunk));
    req.on('end', () => {
      try {
        const body = JSON.parse(Buffer.concat(chunks).toString());
        resolve(body);
      } catch {
        reject(new Error('Invalid JSON body'));
      }
    });
    req.on('error', reject);
  });
}

export function validateSku(sku: string): boolean {
  return /^[A-Z0-9-]+$/.test(sku);
}

export function formatCurrency(amount: number): string {
  return `$${amount.toFixed(2)}`;
}
