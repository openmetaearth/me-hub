typescript
/**
 * src/authentication/qrCodeAuthService.ts
 * 
 * Enhanced QR Code Authentication Service with:
 * - Transaction state tracking
 * - Automatic order cancellation
 * - Idempotency guarantees
 * - Multi-factor fallback
 * - Reward system integration
 * 
 * @module qrCodeAuthService
 */

import { v4 as uuidv4 } from 'uuid';
import { AbortController } from 'node-abort-controller';

// ---------------------------------------------------------------------------
// Types and Interfaces
// ---------------------------------------------------------------------------

export interface Logger {
  debug(message: string, meta?: Record<string, unknown>): void;
  info(message: string, meta?: Record<string, unknown>): void;
  warn(message: string, meta?: Record<string, unknown>): void;
  error(message: string, meta?: Record<string, unknown>): void;
}

export interface QrCodeAuthConfig {
  maxRetries?: number;
  timeoutMs?: number;
  baseDelayMs?: number;
  maxDelayMs?: number;
  jitterFactor?: number;
  logger?: Logger;
  enableMePassFallback?: boolean;
  rewardSystemEnabled?: boolean;
}

export interface QrCodeSession {
  token: string;
  idempotencyKey: string;
  expiresAt: string;
  metadata?: Record<string, unknown>;
  transactionState?: 'pending' | 'completed' | 'failed';
}

export interface RewardPoints {
  current: number;
  earned: number;
  referralBonus?: number;
  activityBonus?: number;
}

export class QrCodeAuthError extends Error {
  public readonly isRetryable: boolean;
  public readonly attempt: number;
  public readonly innerError?: Error;
  public readonly transactionId?: string;

  constructor(
    message: string,
    attempt: number,
    isRetryable: boolean,
    options?: {
      innerError?: Error;
      transactionId?: string;
    }
  ) {
    super(message);
    this.name = 'QrCodeAuthError';
    this.attempt = attempt;
    this.isRetryable = isRetryable;
    this.innerError = options?.innerError;
    this.transactionId = options?.transactionId;
  }
}

export type QrCodeAuthFn = (
  idempotencyKey: string,
  signal: AbortSignal
) => Promise<QrCodeSession>;

export type MePassAuthFn = (
  credentials: { username: string; password: string },
  signal: AbortSignal
) => Promise<QrCodeSession>;

export type RewardCallback = (points: RewardPoints) => void;

// ---------------------------------------------------------------------------
// Constants and Defaults
// ---------------------------------------------------------------------------

const DEFAULT_CONFIG = {
  maxRetries: 3,
  timeoutMs: 30_000,
  baseDelayMs: 1_000,
  maxDelayMs: 10_000,
  jitterFactor: 0.3,
  enableMePassFallback: true,
  rewardSystemEnabled: false,
  logger: {
    debug: (msg: string) => console.debug(`[QR Auth] ${msg}`),
    info: (msg: string) => console.info(`[QR Auth] ${msg}`),
    warn: (msg: string) => console.warn(`[QR Auth] ${msg}`),
    error: (msg: string) => console.error(`[QR Auth] ${msg}`),
  },
} as const;

// ---------------------------------------------------------------------------
// Helper Functions
// ---------------------------------------------------------------------------

function calculateBackoff(attempt: number, config: QrCodeAuthConfig): number {
  const exponentialDelay = Math.min(
    (config.baseDelayMs ?? DEFAULT_CONFIG.baseDelayMs) * Math.pow(2, attempt - 1),
    config.maxDelayMs ?? DEFAULT_CONFIG.maxDelayMs
  );
  const jitter = exponentialDelay * (config.jitterFactor ?? DEFAULT_CONFIG.jitterFactor) * Math.random();
  return Math.floor(exponentialDelay + jitter);
}

async function delayWithSignal(ms: number, signal?: AbortSignal): Promise<void> {
  return new Promise((resolve, reject) => {
    if (signal?.aborted) {
      reject(new DOMException('Delay aborted', 'AbortError'));
      return;
    }

    const timeout = setTimeout(() => {
      cleanup();
      resolve();
    }, ms);

    const onAbort = () => {
      clearTimeout(timeout);
      reject(new DOMException('Delay aborted', 'AbortError'));
    };

    const cleanup = () => {
      signal?.removeEventListener('abort', onAbort);
    };

    signal?.addEventListener('abort', onAbort, { once: true });
  });
}

function validateConfig(config?: QrCodeAuthConfig): Required<QrCodeAuthConfig> {
  const validated = { ...DEFAULT_CONFIG, ...config };

  if (validated.maxRetries < 1) validated.maxRetries = DEFAULT_CONFIG.maxRetries;
  if (validated.timeoutMs < 1_000) validated.timeoutMs = DEFAULT_CONFIG.timeoutMs;
  if (validated.baseDelayMs < 100) validated.baseDelayMs = DEFAULT_CONFIG.baseDelayMs;
  if (validated.maxDelayMs < validated.baseDelayMs) validated.maxDelayMs = DEFAULT_CONFIG.maxDelayMs;
  if (validated.jitterFactor < 0 || validated.jitterFactor > 1) validated.jitterFactor = DEFAULT_CONFIG.jitterFactor;

  return validated;
}

// ---------------------------------------------------------------------------
// Main Service Class
// ---------------------------------------------------------------------------

export class QrCodeAuthService {
  private readonly config: Required<QrCodeAuthConfig>;
  private activeTransactions: Map<string, AbortController> = new Map();
  private rewardCallbacks: Set<RewardCallback> = new Set();

  constructor(config?: QrCodeAuthConfig) {
    this.config = validateConfig(config);
    this.config.logger.info('QrCodeAuthService initialized', {
      config: this.config
    });
  }

  public async authenticate(
    authFn: QrCodeAuthFn,
    mePassAuthFn?: MePassAuthFn,
    credentials?: { username: string; password: string }
  ): Promise<QrCodeSession> {
    const idempotencyKey = uuidv4();
    const transactionId = uuidv4();
    let lastError: Error | null = null;

    this.config.logger.info('Starting authentication flow', { 
      idempotencyKey,
      transactionId
    });

    const controller = new AbortController();
    this.activeTransactions.set(transactionId, controller);

    try {
      for (let attempt = 1; attempt <= this.config.maxRetries; attempt++) {
        try {
          this.config.logger.debug(`Attempt ${attempt} started`, { 
            attempt,
            transactionId
          });

          const session = await this.executeAuthAttempt(
            authFn,
            mePassAuthFn,
            credentials,
            idempotencyKey,
            controller.signal,
            attempt
          );

          this.handleSuccessfulAuth(transactionId, session);
          return session;
        } catch (error) {
          lastError = error as Error;
          const shouldRetry = await this.handleAuthError(
            error,
            attempt,
            transactionId,
            controller
          );
          if (!shouldRetry) break;
        }
      }

      throw new QrCodeAuthError(
        'Maximum retry attempts exceeded',
        this.config.maxRetries,
        false,
        {
          innerError: lastError ?? undefined,
          transactionId
        }
      );
    } finally {
      this.activeTransactions.delete(transactionId);
    }
  }

  public registerRewardCallback(callback: RewardCallback): void {
    this.rewardCallbacks.add(callback);
  }

  public unregisterRewardCallback(callback: RewardCallback): void {
    this.rewardCallbacks.delete(callback);
  }

  public cancelTransaction(transactionId: string): void {
    const controller = this.activeTransactions.get(transactionId);
    if (controller) {
      controller.abort();
      this.config.logger.info('Transaction cancelled', { transactionId });
    }
  }

  private async executeAuthAttempt(
    qrAuthFn: QrCodeAuthFn,
    mePassAuthFn: MePassAuthFn | undefined,
    credentials: { username: string; password: string } | undefined,
    idempotencyKey: string,
    signal: AbortSignal,
    attempt: number
  ): Promise<QrCodeSession> {
    const timeout = setTimeout(() => signal.abort(), this.config.timeoutMs);

    try {
      let session: QrCodeSession;
      
      if (attempt > 1 && this.config.enableMePassFallback && mePassAuthFn && credentials) {
        this.config.logger.debug('Falling back to MePass authentication', { attempt });
        session = await mePassAuthFn(credentials, signal);
      } else {
        session = await qrAuthFn(idempotencyKey, signal);
      }

      clearTimeout(timeout);
      return session;
    } catch (error) {
      clearTimeout(timeout);
      throw error;
    }
  }

  private handleSuccessfulAuth(transactionId: string, session: QrCodeSession): void {
    session.transactionState = 'completed';
    
    if (this.config.rewardSystemEnabled) {
      this.dispatchRewards({
        current: 100,
        earned: 10,
        activityBonus: 5
      });
    }

    this.config.logger.info('Authentication succeeded', { 
      transactionId,
      session: {
        token: session.token,
        expiresAt: session.expiresAt
      }
    });
  }

  private async handleAuthError(
    error: unknown,
    attempt: number,
    transactionId: string,
    controller: AbortController
  ): Promise<boolean> {
    const isRetryable = this.isRetryableError(error);
    const logLevel = isRetryable ? 'warn' : 'error';
    
    this.config.logger[logLevel](`Attempt ${attempt} failed`, { 
      attempt,
      transactionId,
      error: error instanceof Error ? error.message : String(error),
      isRetryable
    });

    if (!isRetryable || attempt === this.config.maxRetries) {
      throw new QrCodeAuthError(
        'Authentication failed',
        attempt,
        isRetryable,
        {
          innerError: error instanceof Error ? error : undefined,
          transactionId
        }
      );
    }

    const backoffMs = calculateBackoff(attempt, this.config);
    this.config.logger.debug(`Waiting ${backoffMs}ms before retry`, { 
      attempt,
      transactionId
    });
    
    try {
      await delayWithSignal(backoffMs, controller.signal);
      return true;
    } catch (abortError) {
      this.config.logger.debug('Retry delay aborted', { 
        transactionId,
        error: abortError instanceof Error ? abortError.message : String(abortError)
      });
      return false;
    }
  }

  private isRetryableError(error: unknown): boolean {
    if (error instanceof DOMException && error.name === 'AbortError') return false;
    if (error instanceof Error && 'code' in error) {
      const code = (error as any).code;
      return !['ECONNABORTED', 'ETIMEDOUT', 'ECANCELED'].includes(code);
    }
    return true;
  }

  private dispatchRewards(points: RewardPoints): void {
    if (!this.config.rewardSystemEnabled) return;
    
    this.rewardCallbacks.forEach(callback => {
      try {
        callback(points);
      } catch (error) {
        this.config.logger.error('Reward callback failed', {
          error: error instanceof Error ? error.message : String(error)
        });
      }
    });
  }
}