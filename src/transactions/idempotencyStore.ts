import { EventEmitter } from 'events';

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

export interface StoredResponse {
  statusCode: number;
  body: unknown;
  headers?: Record<string, string>;
  createdAt: string; // ISO timestamp
}

export interface IdempotencyStore {
  /**
   * Store the response for a given idempotency key.
   * Should throw if the key already exists (to prevent silent overwrites).
   */
  save(key: string, response: StoredResponse, ttlSeconds: number): Promise<void>;

  /**
   * Retrieve a previously stored response by idempotency key.
   * Returns `null` if not found or expired.
   */
  get(key: string): Promise<StoredResponse | null>;

  /**
   * Delete an idempotency key manually (e.g. after successful cleanup).
   */
  delete(key: string): Promise<void>;

  /**
   * Check if a key exists and is still valid.
   */
  exists(key: string): Promise<boolean>;
}

// ---------------------------------------------------------------------------
// In-memory implementation (thread-safe within a single process)
// ---------------------------------------------------------------------------

export class InMemoryIdempotencyStore extends EventEmitter implements IdempotencyStore {
  private readonly store = new Map<string, { response: StoredResponse; expiresAt: number }>();
  private readonly cleanupInterval: NodeJS.Timeout;

  constructor(private readonly defaultTtlMs: number = 60_000) {
    super();
    // Periodic cleanup of expired entries (every 60s)
    this.cleanupInterval = setInterval(() => this.cleanup(), 60_000);
    this.cleanupInterval.unref(); // Don't keep process alive
  }

  async save(key: string, response: StoredResponse, ttlSeconds: number): Promise<void> {
    if (this.store.has(key)) {
      throw new Error(`Idempotency key already exists: ${key}`);
    }
    const expiresAt = Date.now() + ttlSeconds * 1000;
    this.store.set(key, { response, expiresAt });
    this.emit('saved', { key, ttlSeconds });
  }

  async get(key: string): Promise<StoredResponse | null> {
    const entry = this.store.get(key);
    if (!entry) return null;
    if (Date.now() > entry.expiresAt) {
      this.store.delete(key);
      return null;
    }
    return entry.response;
  }

  async delete(key: string): Promise<void> {
    this.store.delete(key);
    this.emit('deleted', { key });
  }

  async exists(key: string): Promise<boolean> {
    const entry = this.store.get(key);
    if (!entry) return false;
    if (Date.now() > entry.expiresAt) {
      this.store.delete(key);
      return false;
    }
    return true;
  }

  private cleanup(): void {
    const now = Date.now();
    for (const [key, entry] of this.store) {
      if (now > entry.expiresAt) {
        this.store.delete(key);
        this.emit('expired', { key });
      }
    }
  }

  destroy(): void {
    clearInterval(this.cleanupInterval);
    this.store.clear();
  }
}

// ---------------------------------------------------------------------------
// Redis implementation (requires ioredis)
// ---------------------------------------------------------------------------

import IORedis from 'ioredis';

export class RedisIdempotencyStore extends EventEmitter implements IdempotencyStore {
  private readonly client: IORedis.Redis;

  constructor(
    private readonly redisOptions?: IORedis.RedisOptions,
    private readonly keyPrefix: string = 'idempotency:',
  ) {
    super();
    this.client = new IORedis(redisOptions);

    this.client.on('error', (err) => {
      this.emit('redis-error', err);
    });
  }

  private buildKey(raw: string): string {
    return `${this.keyPrefix}${raw}`;
  }

  async save(key: string, response: StoredResponse, ttlSeconds: number): Promise<void> {
    const redisKey = this.buildKey(key);
    const serialized = JSON.stringify(response);

    // Use SET NX (not exists) to prevent overwriting existing keys
    const result = await this.client.set(redisKey, serialized, 'EX', ttlSeconds, 'NX');
    if (result === null) {
      throw new Error(`Idempotency key already exists (Redis): ${key}`);
    }
    this.emit('saved', { key, ttlSeconds });
  }

  async get(key: string): Promise<StoredResponse | null> {
    const redisKey = this.buildKey(key);
    const raw = await this.client.get(redisKey);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as StoredResponse;
    } catch {
      // If data is corrupted, delete and return null
      await this.client.del(redisKey);
      return null;
    }
  }

  async delete(key: string): Promise<void> {
    const redisKey = this.buildKey(key);
    await this.client.del(redisKey);
    this.emit('deleted', { key });
  }

  async exists(key: string): Promise<boolean> {
    const redisKey = this.buildKey(key);
    const exists = await this.client.exists(redisKey);
    return exists === 1;
  }

  async quit(): Promise<void> {
    await this.client.quit();
  }
}

// ---------------------------------------------------------------------------
// Factory
// ---------------------------------------------------------------------------

export interface StoreFactoryOptions {
  type?: 'redis' | 'memory';
  redisOptions?: IORedis.RedisOptions;
  defaultTtlSeconds?: number;
}

export function createIdempotencyStore(options: StoreFactoryOptions = {}): IdempotencyStore {
  const type = options.type || (process.env.IDEMPOTENCY_STORE === 'redis' ? 'redis' : 'memory');

  if (type === 'redis') {
    if (typeof IORedis === 'undefined') {
      throw new Error(
        'Redis idempotency store requires the "ioredis" package. Please install it or use in-memory store.'
      );
    }
    return new RedisIdempotencyStore(options.redisOptions);
  }

  // Default: in-memory
  return new InMemoryIdempotencyStore(
    (options.defaultTtlSeconds ?? 60) * 1000
  );
}

// ---------------------------------------------------------------------------
// Export a default instance for convenience (optional)
// ---------------------------------------------------------------------------

let defaultInstance: IdempotencyStore | null = null;

export function getDefaultStore(): IdempotencyStore {
  if (!defaultInstance) {
    defaultInstance = createIdempotencyStore();
  }
  return defaultInstance;
}