/**
 * Maximum number of authentication retry attempts
 * before falling back to alternative methods.
 * Based on UX feedback: up to 3 retries.
 */
export const AUTH_RETRY_MAX_ATTEMPTS = 3;

/**
 * Base delay (in milliseconds) for retry backoff.
 * Actual delay = base * (retry_count ^ 2) + jitter.
 */
export const AUTH_RETRY_BASE_DELAY_MS = 1_000;

/**
 * Default timeout (in milliseconds) for buyer order unavailability.
 * After this period, the order is automatically cancelled.
 */
export const ORDER_BUYER_TIMEOUT_MS = 5 * 60 * 1_000; // 5 minutes

/**
 * Time-to-live (in milliseconds) for idempotency keys.
 * Prevents duplicate processing of identical requests within this window.
 * Default: 24 hours.
 */
export const IDEMPOTENCY_TTL_MS = 24 * 60 * 60 * 1_000; // 24 hours

/**
 * Configuration objects for grouped usage.
 */
export const authRetryConfig = {
  maxAttempts: AUTH_RETRY_MAX_ATTEMPTS,
  baseDelayMs: AUTH_RETRY_BASE_DELAY_MS,
} as const;

export const orderConfig = {
  buyerTimeoutMs: ORDER_BUYER_TIMEOUT_MS,
} as const;

export const idempotencyConfig = {
  ttlMs: IDEMPOTENCY_TTL_MS,
} as const;