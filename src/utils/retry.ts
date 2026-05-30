/**
 * Flatiron retry utility with exponential backoff, jitter, timeout, and abort support.
 *
 * @module retry
 */

/** Configuration options for the retry operation. */
export interface RetryOptions {
  /** Maximum number of retry attempts (excluding the initial call). Default: 3 */
  maxRetries?: number;
  /** Base delay in milliseconds before the first retry. Default: 1000 */
  baseDelay?: number;
  /** Exponential factor applied to delay after each attempt. Default: 2 */
  factor?: number;
  /** Maximum delay in milliseconds between retries. Default: 30000 */
  maxDelay?: number;
  /** Whether to apply full jitter (random delay between 0 and current delay). Default: true */
  jitter?: boolean;
  /** Optional function to decide whether a retry should be attempted based on the error. Default: retry all errors */
  retryIf?: (error: unknown) => boolean;
  /** Optional callback invoked before each retry (receives attempt number and delay). */
  onRetry?: (attempt: number, delayMs: number) => void;
  /** Optional AbortSignal to cancel the retry loop. */
  signal?: AbortSignal;
  /** Optional timeout in milliseconds for each individual attempt. Enforced via AbortController. */
  attemptTimeout?: number;
}

/** Default retry options. */
const DEFAULT_OPTIONS: Required<RetryOptions> = {
  maxRetries: 3,
  baseDelay: 1000,
  factor: 2,
  maxDelay: 30000,
  jitter: true,
  retryIf: () => true,
  onRetry: () => {},
  signal: undefined as AbortSignal | undefined,
  attemptTimeout: undefined as number | undefined,
};

/**
 * Sleep for the given number of milliseconds, abortable via signal.
 *
 * @param ms - Milliseconds to sleep.
 * @param signal - Optional abort signal.
 * @returns A promise that resolves after the sleep duration, rejects if aborted.
 */
function sleep(ms: number, signal?: AbortSignal): Promise<void> {
  return new Promise((resolve, reject) => {
    if (signal?.aborted) {
      return reject(new DOMException('Aborted', 'AbortError'));
    }

    const onAbort = () => {
      clearTimeout(timer);
      reject(new DOMException('Aborted', 'AbortError'));
    };

    const timer = setTimeout(() => {
      signal?.removeEventListener('abort', onAbort);
      resolve();
    }, ms);

    if (signal) {
      signal.addEventListener('abort', onAbort, { once: true });
    }
  });
}

/**
 * Executes an async operation with exponential backoff retry logic.
 *
 * The function will call `operation` immediately. If it throws and retry options allow,
 * it will wait for an exponentially increasing delay (with optional jitter) and retry.
 * All retry attempts can be cancelled via an `AbortSignal`.
 * An optional per-attempt timeout can be specified via `attemptTimeout`.
 *
 * @typeParam T - The return type of the operation.
 * @param operation - Async function to retry.
 * @param options - Configuration for retry behavior.
 * @returns The result of the successful operation.
 * @throws The last error encountered if all retries fail or if the operation is aborted.
 *
 * @example
 *