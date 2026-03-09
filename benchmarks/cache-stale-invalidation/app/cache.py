"""
In-memory cache with TTL support.
"""
import time


class Cache:
    def __init__(self, ttl_seconds=300):
        self.ttl_seconds = ttl_seconds
        self._store = {}  # key -> (value, expires_at)
        self._hits = 0
        self._misses = 0

    def get(self, key):
        """Get a value from cache. Returns None if not found or expired."""
        entry = self._store.get(key)
        if entry is None:
            self._misses += 1
            return None

        value, expires_at = entry
        if time.time() > expires_at:
            # Expired — remove and treat as miss
            del self._store[key]
            self._misses += 1
            return None

        self._hits += 1
        return value

    def set(self, key, value):
        """Set a value in cache with TTL."""
        expires_at = time.time() + self.ttl_seconds
        self._store[key] = (value, expires_at)

    def delete(self, key):
        """Delete a specific key from cache."""
        if key in self._store:
            del self._store[key]
            return True
        return False

    def invalidate_prefix(self, prefix):
        """Delete all keys matching a prefix."""
        keys_to_delete = [k for k in self._store if k.startswith(prefix)]
        for key in keys_to_delete:
            del self._store[key]
        return len(keys_to_delete)

    def clear(self):
        """Clear all cached entries."""
        self._store.clear()

    def stats(self):
        """Return cache statistics."""
        now = time.time()
        active = sum(1 for _, (_, exp) in self._store.items() if exp > now)
        return {
            'total_entries': len(self._store),
            'active_entries': active,
            'hits': self._hits,
            'misses': self._misses,
            'hit_rate': self._hits / max(self._hits + self._misses, 1),
            'ttl_seconds': self.ttl_seconds,
        }
