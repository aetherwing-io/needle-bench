import http from 'http';

interface RateLimitConfig {
  windowMs: number;
  maxRequests: number;
}

interface ClientRecord {
  count: number;
  windowStart: number;
}

export class RateLimiter {
  private config: RateLimitConfig;
  private clients: Map<string, ClientRecord>;

  constructor(config: RateLimitConfig) {
    this.config = config;
    this.clients = new Map();
  }

  /**
   * Get the client IP address from the request.
   *
   * BUG: Trusts X-Forwarded-For header without validation.
   * An attacker can set a different X-Forwarded-For value on each request
   * to get a fresh rate limit window every time, effectively bypassing
   * the rate limiter entirely.
   */
  getClientIp(req: http.IncomingMessage): string {
    // Trust proxy headers — common in load balancer setups
    const forwarded = req.headers['x-forwarded-for'];
    if (forwarded) {
      // Take the first IP in the chain
      const ip = Array.isArray(forwarded) ? forwarded[0] : forwarded.split(',')[0].trim();
      return ip;
    }

    // Fall back to socket address
    return req.socket.remoteAddress || '127.0.0.1';
  }

  /**
   * Check if a request from the given IP is allowed.
   */
  allow(ip: string): boolean {
    const now = Date.now();
    const record = this.clients.get(ip);

    if (!record || now - record.windowStart >= this.config.windowMs) {
      // New window
      this.clients.set(ip, { count: 1, windowStart: now });
      return true;
    }

    if (record.count >= this.config.maxRequests) {
      return false;
    }

    record.count++;
    return true;
  }

  /**
   * Get time in ms until the client can make another request.
   */
  getRetryAfter(ip: string): number {
    const record = this.clients.get(ip);
    if (!record) return 0;

    const elapsed = Date.now() - record.windowStart;
    return Math.max(0, this.config.windowMs - elapsed);
  }

  /**
   * Get current stats for monitoring.
   */
  getStats(): object {
    const activeClients = this.clients.size;
    let totalRequests = 0;
    for (const [, record] of this.clients) {
      totalRequests += record.count;
    }
    return {
      activeClients,
      totalRequests,
      config: this.config,
    };
  }

  /**
   * Clean up expired entries.
   */
  cleanup(): void {
    const now = Date.now();
    for (const [ip, record] of this.clients) {
      if (now - record.windowStart >= this.config.windowMs) {
        this.clients.delete(ip);
      }
    }
  }
}
