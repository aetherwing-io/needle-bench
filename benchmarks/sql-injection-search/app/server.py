"""
Product search API backed by SQLite.
"""
import json
import sys
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs

from database import Database
from search import SearchEngine


db = Database('products.db')
search = SearchEngine(db)


class SearchHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed = urlparse(self.path)
        path = parsed.path
        params = parse_qs(parsed.query)

        if path == '/search':
            query = params.get('q', [''])[0]
            category = params.get('category', [None])[0]
            self._search(query, category)
        elif path == '/products':
            self._list_products()
        elif path.startswith('/products/'):
            product_id = path.split('/')[-1]
            self._get_product(product_id)
        elif path == '/health':
            self._json_response(200, {'status': 'ok'})
        else:
            self._json_response(404, {'error': 'not found'})

    def _search(self, query, category):
        if not query:
            self._json_response(400, {'error': 'query parameter q is required'})
            return

        try:
            results = search.search(query, category)
            self._json_response(200, {
                'query': query,
                'count': len(results),
                'results': results,
            })
        except Exception as e:
            self._json_response(500, {'error': str(e)})

    def _list_products(self):
        products = db.execute('SELECT * FROM products')
        self._json_response(200, {'products': products})

    def _get_product(self, product_id):
        rows = db.execute('SELECT * FROM products WHERE id = ?', (product_id,))
        if not rows:
            self._json_response(404, {'error': f'product {product_id} not found'})
            return
        self._json_response(200, rows[0])

    def _json_response(self, status, data):
        self.send_response(status)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(data).encode())

    def log_message(self, format, *args):
        pass


def main():
    port = int(sys.argv[1]) if len(sys.argv) > 1 else 8080

    # Initialize database with sample data
    db.initialize()

    server = HTTPServer(('', port), SearchHandler)
    print(f'Search server listening on port {port}')
    server.serve_forever()


if __name__ == '__main__':
    main()
