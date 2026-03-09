"""
Application configuration.
"""
import os

PORT = int(os.environ.get('PORT', 8080))
CACHE_TTL = int(os.environ.get('CACHE_TTL', 300))  # seconds
LOG_LEVEL = os.environ.get('LOG_LEVEL', 'INFO')
