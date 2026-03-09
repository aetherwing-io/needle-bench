import http from 'http';
import { RateLimiter } from './ratelimiter';
import { parseBody, jsonResponse } from './utils';
import { router } from './routes';

const PORT = parseInt(process.env.PORT || '8080', 10);

const limiter = new RateLimiter({
  windowMs: 5_000,    // 5 second window (short for testability)
  maxRequests: 10,     // 10 requests per window per IP
});

const server = http.createServer(async (req, res) => {
  // Apply rate limiting
  const clientIp = limiter.getClientIp(req);
  if (!limiter.allow(clientIp)) {
    jsonResponse(res, 429, {
      error: 'Too many requests',
      retryAfter: Math.ceil(limiter.getRetryAfter(clientIp) / 1000),
    });
    return;
  }

  // Route the request
  try {
    await router(req, res);
  } catch (err: any) {
    jsonResponse(res, 500, { error: err.message });
  }
});

server.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}`);
});

export { server, limiter };
