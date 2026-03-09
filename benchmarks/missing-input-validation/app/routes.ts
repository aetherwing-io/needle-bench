import http from 'http';
import { Inventory } from './inventory';
import { parseBody } from './utils';

const inventory = new Inventory();

// Seed some initial products
inventory.addProduct({ sku: 'WIDGET-001', name: 'Standard Widget', quantity: 100, price: 9.99 });
inventory.addProduct({ sku: 'GADGET-002', name: 'Deluxe Gadget', quantity: 50, price: 24.99 });
inventory.addProduct({ sku: 'SENSOR-003', name: 'Precision Sensor', quantity: 200, price: 14.99 });

export async function handleInventory(
  req: http.IncomingMessage,
  res: http.ServerResponse,
  url: URL
): Promise<void> {
  // GET /inventory — list all products
  if (req.method === 'GET' && url.pathname === '/inventory') {
    const products = inventory.listProducts();
    res.end(JSON.stringify({ products, count: products.length }));
    return;
  }

  // GET /inventory/:sku — get single product
  const skuMatch = url.pathname.match(/^\/inventory\/([A-Z0-9-]+)$/);
  if (req.method === 'GET' && skuMatch) {
    const product = inventory.getProduct(skuMatch[1]);
    if (!product) {
      res.statusCode = 404;
      res.end(JSON.stringify({ error: 'product not found' }));
      return;
    }
    res.end(JSON.stringify(product));
    return;
  }

  // POST /inventory/:sku/adjust — adjust quantity
  const adjustMatch = url.pathname.match(/^\/inventory\/([A-Z0-9-]+)\/adjust$/);
  if (req.method === 'POST' && adjustMatch) {
    const body = await parseBody(req);
    const sku = adjustMatch[1];

    const product = inventory.getProduct(sku);
    if (!product) {
      res.statusCode = 404;
      res.end(JSON.stringify({ error: 'product not found' }));
      return;
    }

    const quantity = body.quantity as number;
    if (typeof quantity !== 'number' || isNaN(quantity)) {
      res.statusCode = 400;
      res.end(JSON.stringify({ error: 'quantity must be a number' }));
      return;
    }

    const result = inventory.adjustQuantity(sku, quantity);
    res.end(JSON.stringify(result));
    return;
  }

  res.statusCode = 405;
  res.end(JSON.stringify({ error: 'method not allowed' }));
}
