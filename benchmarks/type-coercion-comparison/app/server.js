const http = require('http');
const { filterProducts } = require('./filter');
const { products } = require('./data');
const { parseQuery } = require('./utils');

const PORT = 3000;

const server = http.createServer((req, res) => {
  res.setHeader('Content-Type', 'application/json');

  const url = new URL(req.url, `http://localhost:${PORT}`);

  if (url.pathname === '/health') {
    res.end(JSON.stringify({ status: 'ok' }));
    return;
  }

  if (url.pathname === '/products') {
    const query = parseQuery(url.searchParams);
    const results = filterProducts(products, query);

    res.end(JSON.stringify({
      query: query,
      count: results.length,
      products: results,
    }));
    return;
  }

  if (url.pathname === '/products/count') {
    const query = parseQuery(url.searchParams);
    const results = filterProducts(products, query);

    res.end(JSON.stringify({ count: results.length }));
    return;
  }

  res.statusCode = 404;
  res.end(JSON.stringify({ error: 'not found' }));
});

server.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});

module.exports = server;
