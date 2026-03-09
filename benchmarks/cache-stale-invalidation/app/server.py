"""
Product catalog API with caching layer.
"""
import json
import sys
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs

from cache import Cache
from store import ProductStore


cache = Cache(ttl_seconds=300)  # 5 minute TTL
store = ProductStore()


class CatalogHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed = urlparse(self.path)
        path = parsed.path
        params = parse_qs(parsed.query)

        if path == '/products':
            self._list_products(params)
        elif path.startswith('/products/'):
            product_id = path.split('/')[-1]
            self._get_product(product_id)
        elif path == '/health':
            self._json_response(200, {'status': 'ok'})
        elif path == '/cache/stats':
            self._json_response(200, cache.stats())
        else:
            self._json_response(404, {'error': 'not found'})

    def do_PUT(self):
        parsed = urlparse(self.path)
        path = parsed.path

        if path.startswith('/products/'):
            product_id = path.split('/')[-1]
            body = self._read_body()
            self._update_product(product_id, body)
        else:
            self._json_response(404, {'error': 'not found'})

    def do_POST(self):
        parsed = urlparse(self.path)
        path = parsed.path

        if path == '/products':
            body = self._read_body()
            self._create_product(body)
        else:
            self._json_response(404, {'error': 'not found'})

    def _list_products(self, params):
        cache_key = 'products:list'
        cached = cache.get(cache_key)
        if cached is not None:
            self._json_response(200, cached)
            return

        products = store.list_products()
        cache.set(cache_key, products)
        self._json_response(200, products)

    def _get_product(self, product_id):
        cache_key = f'products:{product_id}'
        cached = cache.get(cache_key)
        if cached is not None:
            self._json_response(200, cached)
            return

        product = store.get_product(product_id)
        if product is None:
            self._json_response(404, {'error': f'product {product_id} not found'})
            return

        cache.set(cache_key, product)
        self._json_response(200, product)

    def _update_product(self, product_id, data):
        product = store.get_product(product_id)
        if product is None:
            self._json_response(404, {'error': f'product {product_id} not found'})
            return

        updated = store.update_product(product_id, data)

        # BUG: We update the store but never invalidate the cache.
        # The cache still has the old product data and will serve it
        # until the TTL expires (up to 5 minutes of stale reads).

        self._json_response(200, updated)

    def _create_product(self, data):
        product = store.create_product(data)

        # BUG: We create a new product but don't invalidate the list cache.
        # GET /products will return the stale list (missing the new product)
        # until the cache TTL expires.

        self._json_response(201, product)

    def _read_body(self):
        content_length = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(content_length)
        return json.loads(body)

    def _json_response(self, status, data):
        self.send_response(status)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(data).encode())

    def log_message(self, format, *args):
        # Suppress request logs during testing
        pass


def main():
    port = int(sys.argv[1]) if len(sys.argv) > 1 else 8080
    server = HTTPServer(('', port), CatalogHandler)
    print(f'Catalog server listening on port {port}')
    server.serve_forever()


if __name__ == '__main__':
    main()
