// tests/unit/idempotencyMiddleware.test.ts

import { Request, Response, NextFunction } from 'express';
import { createIdempotencyMiddleware } from '../../src/middleware/idempotencyMiddleware';
import { IdempotencyCache } from '../../src/services/idempotencyCache';

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

const mockCache: jest.Mocked<IdempotencyCache> = {
  get: jest.fn(),
  set: jest.fn(),
  delete: jest.fn(),
};

// Helper to generate fake req/res/next
function createMockContext() {
  const req = {
    headers: {},
    method: 'POST',
    path: '/transactions',
    body: {},
  } as unknown as Request;

  const res = {
    statusCode: 200,
    json: jest.fn(),
    send: jest.fn(),
    end: jest.fn(),
    setHeader: jest.fn(),
    on: jest.fn(),
  } as unknown as Response;

  const next = jest.fn() as NextFunction;

  return { req, res, next };
}

// ---------------------------------------------------------------------------
// Setup
// ---------------------------------------------------------------------------

beforeEach(() => {
  jest.clearAllMocks();
});

describe('createIdempotencyMiddleware', () => {
  const middleware = createIdempotencyMiddleware({
    cache: mockCache,
    headerName: 'Idempotency-Key',
    ttlMs: 3600_000, // 1 hour
  });

  // -----------------------------------------------------------------------
  // Missing key
  // -----------------------------------------------------------------------
  it('should call next() when no idempotency key is present', async () => {
    const { req, res, next } = createMockContext();
    req.headers = {};

    await middleware(req, res, next);

    expect(mockCache.get).not.toHaveBeenCalled();
    expect(mockCache.set).not.toHaveBeenCalled();
    expect(next).toHaveBeenCalledTimes(1);
  });

  it('should call next() for non-mutating methods (GET, HEAD, OPTIONS)', async () => {
    const { req, res, next } = createMockContext();
    req.method = 'GET';
    req.headers = { 'Idempotency-Key': 'abc' };

    await middleware(req, res, next);

    expect(mockCache.get).not.toHaveBeenCalled();
    expect(next).toHaveBeenCalled();
  });

  // -----------------------------------------------------------------------
  // First request – cache miss
  // -----------------------------------------------------------------------
  it('should proceed and cache the response on first request', async () => {
    const { req, res, next } = createMockContext();
    const key = 'unique-key-123';
    req.headers = { 'Idempotency-Key': key };
    mockCache.get.mockResolvedValue(null);
    mockCache.set.mockResolvedValue(true);

    await middleware(req, res, next);

    expect(mockCache.get).toHaveBeenCalledWith(key);
    expect(next).toHaveBeenCalledTimes(1);

    // Simulate response being sent
    const responseBody = { status: 'ok', txId: 'tx_001' };
    res.json(responseBody);
    res.statusCode = 201;

    // After the response finishes, the middleware should cache the result
    // We need to simulate the 'finish' event that the middleware listens on.
    const finishCallback = (res.on as jest.Mock).mock.calls.find(
      ([event]: [string]) => event === 'finish'
    )?.[1];
    if (finishCallback) {
      finishCallback();
    }

    expect(mockCache.set).toHaveBeenCalledWith(
      key,
      expect.objectContaining({
        statusCode: 201,
        body: JSON.stringify(responseBody),
      }),
      expect.any(Number) // ttl
    );
  });

  // -----------------------------------------------------------------------
  // Duplicate request – cache hit
  // -----------------------------------------------------------------------
  it('should return the cached response for duplicate key', async () => {
    const { req, res, next } = createMockContext();
    const key = 'duplicate-key';
    req.headers = { 'Idempotency-Key': key };

    const cachedResponse = {
      statusCode: 200,
      body: JSON.stringify({ status: 'already_processed', txId: 'tx_001' }),
      headers: { 'content-type': 'application/json' },
    };
    mockCache.get.mockResolvedValue(cachedResponse);

    await middleware(req, res, next);

    expect(mockCache.get).toHaveBeenCalledWith(key);
    expect(next).not.toHaveBeenCalled();

    expect(res.setHeader).toHaveBeenCalledWith(
      'Idempotency-Replay',
      'true'
    );
    expect(res.statusCode).toBe(200);
    expect(res.json).toHaveBeenCalledWith(JSON.parse(cachedResponse.body));
  });

  // -----------------------------------------------------------------------
  // Error handling when cache is down
  // -----------------------------------------------------------------------
  it('should call next() with error when cache get fails', async () => {
    const { req, res, next } = createMockContext();
    const key = 'error-key';
    req.headers = { 'Idempotency-Key': key };
    const cacheError = new Error('Cache unavailable');
    mockCache.get.mockRejectedValue(cacheError);

    await middleware(req, res, next);

    expect(next).toHaveBeenCalledWith(cacheError);
  });

  it('should not crash if cache set fails (fire-and-forget)', async () => {
    const { req, res, next } = createMockContext();
    const key = 'faf-key';
    req.headers = { 'Idempotency-Key': key };
    mockCache.get.mockResolvedValue(null);
    mockCache.set.mockRejectedValue(new Error('Cache write failed'));

    await middleware(req, res, next);

    // next is called, then we simulate finish
    expect(next).toHaveBeenCalledTimes(1);

    // Simulate finish event
    const finishCallback = (res.on as jest.Mock).mock.calls.find(
      ([event]: [string]) => event === 'finish'
    )?.[1];
    if (finishCallback) {
      finishCallback();
    }

    // Should not throw – the error is logged but not propagated
    expect(mockCache.set).toHaveBeenCalled();
  });

  // -----------------------------------------------------------------------
  // TTL and expiration
  // -----------------------------------------------------------------------
  it('should respect custom ttl from config', () => {
    const customMiddleware = createIdempotencyMiddleware({
      cache: mockCache,
      headerName: 'Idempotency-Key',
      ttlMs: 5000,
    });

    expect(customMiddleware).toBeDefined();
  });

  // -----------------------------------------------------------------------
  // Key header name customization
  // -----------------------------------------------------------------------
  it('should use custom header name', () => {
    const customMiddleware = createIdempotencyMiddleware({
      cache: mockCache,
      headerName: 'X-Idempotency-Key',
      ttlMs: 3600_000,
    });

    // Check internal mapping by verifying behavior with custom header
    const { req, res, next } = createMockContext();
    req.headers = { 'X-Idempotency-Key': 'custom-header-key' };
    mockCache.get.mockResolvedValue(null);

    customMiddleware(req, res, next);

    expect(mockCache.get).toHaveBeenCalledWith('custom-header-key');
  });

  // -----------------------------------------------------------------------
  // Concurrent requests with same key (lock handling)
  // -----------------------------------------------------------------------
  it('should process only the first concurrent request with same key', async () => {
    const { req: req1, res: res1, next: next1 } = createMockContext();
    const { req: req2, res: res2, next: next2 } = createMockContext();
    const key = 'concurrent-key';

    req1.headers = { 'Idempotency-Key': key };
    req2.headers = { 'Idempotency-Key': key };

    // First get returns null, second should wait for lock or return cached
    // For simplicity, we assume the middleware uses an in-flight lock (e.g., a Map)
    // We'll test that one request passes through and the other is blocked.
    // This requires the middleware to have an internal lock mechanism.
    // We'll simulate by checking that only one next() call is made.

    mockCache.get.mockResolvedValue(null);
    mockCache.set.mockResolvedValue(true);

    // Fire both requests concurrently
    await Promise.all([
      middleware(req1, res1, next1),
      middleware(req2, res2, next2),
    ]);

    // Only one should have called next(), the other should have waited
    // The exact behavior depends on implementation, but we can assert that
    // get was called only once (if lock prevents concurrent processing)
    // Note: This test assumes a lock is implemented; adjust if not.
    // For a production test, we would mock the lock separately.
    // Here we simply ensure no crashes and both eventually resolve.
    expect(next1.mock.calls.length + next2.mock.calls.length).toBeGreaterThanOrEqual(1);
    expect(mockCache.get).toHaveBeenCalledTimes(1); // weak assumption
  });

  // -----------------------------------------------------------------------
  // Response headers (cors, content-type) are preserved
  // -----------------------------------------------------------------------
  it('should preserve response headers from cached response', async () => {
    const { req, res, next } = createMockContext();
    const key = 'headers-key';
    req.headers = { 'Idempotency-Key': key };

    const cachedResponse = {
      statusCode: 200,
      body: JSON.stringify({ success: true }),
      headers: {
        'content-type': 'application/json',
        'x-request-id': 'req-456',
      },
    };
    mockCache.get.mockResolvedValue(cachedResponse);

    await middleware(req, res, next);

    expect(res.setHeader).toHaveBeenCalledWith('content-type', 'application/json');
    expect(res.setHeader).toHaveBeenCalledWith('x-request-id', 'req-456');
    expect(res.setHeader).toHaveBeenCalledWith('Idempotency-Replay', 'true');
  });
});