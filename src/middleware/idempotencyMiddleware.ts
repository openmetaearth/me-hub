typescript
import { Request, Response, NextFunction } from 'express';
import crypto from 'crypto';

// ---------------------------------------------------------------------------
// Logging
// ---------------------------------------------------------------------------
const logger = {
  info: (msg: string, meta?: Record<string, unknown>) =>
    console.log(JSON.stringify({ level: 'info', msg, ...meta })),
  warn: (msg: string, meta?: Record<string, unknown>) =>
    console.warn(JSON.stringify({ level: 'warn', msg, ...meta })),
  error: (msg: string, meta?: Record<string, unknown>) =>
    console.error(JSON.stringify({ level: 'error', msg, ...meta })),
};

// ---------------------------------------------------------------------------
// Types & Interfaces
// ---------------------------------------------------------------------------

export interface IdempotencyOptions {
  /**
   * HTTP header name for the idempotency key.
   * @default 'Idempotency-Key'
   */
  headerName?: string;

  /**
   * Time‑to‑live in milliseconds for a completed response.
   * @default 86400000 (24 hours)
   */
  ttlCompleted?: number;

  /**
   * Time‑to‑live in milliseconds for an in‑flight (pending) response.
   * @default 30000 (30 seconds)
   */
  ttlInflight?: number;

  /**
   * HTTP methods that are considered idempotent.
   * @default ['POST', 'PUT', 'PATCH', 'DELETE']
   */
  idempotentMethods?: string[];

  /**
   * Custom storage backend. Defaults to an in‑memory Map with TTL sweep.
   */
  store?: IdempotencyStore;

  /**
   * Maximum request body size in bytes to cache.
   * Caching the body is required to later re‑play the request if needed.
   * Use -1 for unlimited (not recommended). Default: 65536 (64 KB).
   */
  maxCachedBodySize?: number;

  /**
   * If true, the middleware will reject requests without an idempotency key
   * with a 400 error (strict mode). Otherwise it passes through.
   * @default true
   */
  requireKey?: boolean;

  /**
   * Optional prefix for the idempotency key in storage (e.g. "user:123:").
   * Useful for namespacing.
   */
  keyPrefix?: string;

  /**
   * Optional function to compute a hash of the request body to detect
   * duplicate requests with different keys but same payload.
   * Return null to skip body hash check.
   */
  hashBody?: (body: unknown) => string | null;
}

export interface CachedResponse {
  statusCode: number;
  headers: Record<string, string | number | string[]>;
  body: any; // raw body as sent by res.send/res.json
  /** SHA‑256 of the original request body (optional, used for consistency check) */
  bodyHash?: string;
}

export interface InflightEntry {
  response: Promise<CachedResponse | null>;
  createdAt: number;
}

export interface CompletedEntry {
  response: CachedResponse;
  createdAt: number;
}

export interface IdempotencyStore {
  get(key: string): Promise<CompletedEntry | InflightEntry | undefined>;
  setInFlight(key: string, entry: InflightEntry, ttl: number): Promise<void>;
  setCompleted(key: string, entry: CompletedEntry, ttl: number): Promise<void>;
  delete(key: string): Promise<void>;
  /** Optional periodic cleanup */
  sweep?(): Promise<void>;
}

// ---------------------------------------------------------------------------
// Default In‑Memory Store with TTL Sweep
// ---------------------------------------------------------------------------

interface MapValue {
  type: 'completed' | 'inflight';
  entry: CompletedEntry | InflightEntry;
  expiresAt: number;
}

export class InMemoryIdempotencyStore implements IdempotencyStore {
  private readonly store = new Map<string, MapValue>();

  /**
   * @param sweepIntervalMs How often (ms) to run cleanup of expired entries.
   *                        Pass 0 to disable auto‑sweep.
   */
  constructor(sweepIntervalMs = 60_000) {
    if (sweepIntervalMs > 0) {
      const interval = setInterval(() => this.sweep().catch(logger.error), sweepIntervalMs);
      if (interval.unref) interval.unref();
    }
  }

  async get(key: string): Promise<CompletedEntry | InflightEntry | undefined> {
    const val = this.store.get(key);
    if (!val) return undefined;
    if (Date.now() > val.expiresAt) {
      this.store.delete(key);
      return undefined;
    }
    return val.entry;
  }

  async setInFlight(key: string, entry: InflightEntry, ttl: number): Promise<void> {
    this.store.set(key, {
      type: 'inflight',
      entry,
      expiresAt: Date.now() + ttl,
    });
  }

  async setCompleted(key: string, entry: CompletedEntry, ttl: number): Promise<void> {
    this.store.set(key, {
      type: 'completed',
      entry,
      expiresAt: Date.now() + ttl,
    });
  }

  async delete(key: string): Promise<void> {
    this.store.delete(key);
  }

  async sweep(): Promise<void> {
    const now = Date.now();
    for (const [key, val] of this.store.entries()) {
      if (now > val.expiresAt) {
        this.store.delete(key);
      }
    }
  }

  dispose(): void {
    // No timers to clear in this version (interval is unref'd)
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/**
 * Validates that an idempotency key is a non‑empty string with safe characters.
 * We accept UUIDs or any alphanumeric string with hyphens/underscores.
 */
function isValidKey(key: string): boolean {
  return typeof key === 'string' && /^[A-Za-z0-9\-_]{1,64}$/.test(key);
}

/**
 * Creates a SHA‑256 hex digest of the given body (JSON‑serialized).
 */
function defaultHashBody(body: unknown): string | null {
  if (body === undefined || body === null) return null;
  try {
    return crypto.createHash('sha256').update(JSON.stringify(body)).digest('hex');
  } catch {
    return null;
  }
}

// ---------------------------------------------------------------------------
// Middleware Factory
// ---------------------------------------------------------------------------

/**
 * Creates an Express middleware that enforces idempotency for mutating requests.
 * See {@link IdempotencyOptions} for configuration.
 */
export function createIdempotencyMiddleware(
  options: IdempotencyOptions = {}
) {
  const {
    headerName = 'Idempotency-Key',
    ttlCompleted = 86_400_000,       // 24h
    ttlInflight = 30_000,            // 30s
    idempotentMethods = ['POST', 'PUT', 'PATCH', 'DELETE'],
    store = new InMemoryIdempotencyStore(),
    maxCachedBodySize = 65_536,      // 64 KB
    requireKey = true,
    keyPrefix = '',
    hashBody = defaultHashBody,
  } = options;

  const headerNameLower = headerName.toLowerCase();

  return async function idempotencyMiddleware(
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> {
    // -----------------------------------------------------------------------
    // 1. Only apply to configured idempotent methods
    // -----------------------------------------------------------------------
    if (!idempotentMethods.includes(req.method.toUpperCase())) {
      return next();
    }

    // -----------------------------------------------------------------------
    // 2. Validate idempotency key
    // -----------------------------------------------------------------------
    const rawKey = req.headers[headerNameLower] as string | undefined;

    if (!rawKey) {
      if (requireKey) {
        logger.warn('Idempotency key missing', {
          method: req.method,
          path: req.path,
          ip: req.ip,
        });
        res.status(400).json({
          error: 'Bad Request',
          message: `Idempotency key is required via the '${headerName}' header.`,
        });
        return;
      }
      return next(); // Silent pass‑through
    }

    if (!isValidKey(rawKey)) {
      logger.warn('Invalid idempotency key', { key: rawKey });
      res.status(400).json({
        error: 'Bad Request',
        message: `Idempotency key must be a string of 1‑64 alphanumeric characters, hyphens, or underscores.`,
      });
      return;
    }

    const key = keyPrefix + rawKey;

    // -----------------------------------------------------------------------
    // 3. Check for existing response
    // -----------------------------------------------------------------------
    const existing = await store.get(key);
    if (existing) {
      if (existing.type === 'completed') {
        logger.info('Replaying completed response', { key });
        const { statusCode, headers, body } = existing.response;
        res.status(statusCode).set(headers).send(body);
        return;
      } else if (existing.type === 'inflight') {
        logger.info('Request already in flight', { key });
        res.status(409).json({
          error: 'Conflict',
          message: 'Request already in progress. Please retry later.',
        });
        return;
      }
    }

    // -----------------------------------------------------------------------
    // 4. Validate request body size
    // -----------------------------------------------------------------------
    const bodySize = req.headers['content-length'] ? parseInt(req.headers['content-length'], 10) : 0;
    if (bodySize > maxCachedBodySize) {
      logger.warn('Request body too large', { key, bodySize });
      res.status(413).json({
        error: 'Payload Too Large',
        message: `Request body exceeds maximum allowed size of ${maxCachedBodySize} bytes.`,
      });
      return;
    }

    // -----------------------------------------------------------------------
    // 5. Compute body hash
    // -----------------------------------------------------------------------
    const bodyHash = hashBody(req.body);

    // -----------------------------------------------------------------------
    // 6. Store in-flight entry
    // -----------------------------------------------------------------------
    const inflightEntry: InflightEntry = {
      response: new Promise((resolve) => {
        const originalSend = res.send.bind(res);
        const originalJson = res.json.bind(res);
        const originalEnd = res.end.bind(res);

        const response: CachedResponse = {
          statusCode: 200,
          headers: {},
          body: null,
        };

        res.send = (body?: any): Response => {
          response.body = body;
          return originalSend(body);
        };

        res.json = (body?: any): Response => {
          response.body = body;
          return originalJson(body);
        };

        res.end = (body?: any): void => {
          if (body) response.body = body;
          originalEnd(body);
        };

        res.on('finish', () => {
          response.statusCode = res.statusCode;
          response.headers = res.getHeaders();
          if (bodyHash) response.bodyHash = bodyHash;
          resolve(response);
        });
      }),
      createdAt: Date.now(),
    };

    await store.setInFlight(key, inflightEntry, ttlInflight);

    // -----------------------------------------------------------------------
    // 7. Proceed with request
    // -----------------------------------------------------------------------
    try {
      await next();
    } catch (error) {
      logger.error('Request failed', { key, error });
      await store.delete(key);
      throw error;
    }

    // -----------------------------------------------------------------------
    // 8. Store completed response
    // -----------------------------------------------------------------------
    const response = await inflightEntry.response;
    if (response) {
      await store.setCompleted(key, { response, createdAt: Date.now() }, ttlCompleted);
    }
  };
}