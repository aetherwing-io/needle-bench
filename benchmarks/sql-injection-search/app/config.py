"""
Application configuration.
"""
import os

PORT = int(os.environ.get('PORT', 8080))
DB_PATH = os.environ.get('DB_PATH', 'products.db')
LOG_LEVEL = os.environ.get('LOG_LEVEL', 'INFO')
MAX_SEARCH_RESULTS = int(os.environ.get('MAX_RESULTS', 100))
