typescript
// src/transactions/transactionService.ts

/**
 * Core business logic for processing payments with idempotency enforcement, retry logic,
 * full error handling, and event emission. Designed for ME-Hub reliability: prevents
 * duplicate withdrawals, double-spend, and ensures eventual consistency.
 */

// ──────────────────────────────────────
// Domain Types
// ──────────────────────────────────────

export interface PaymentRequest {
  /** Unique idempotency key provided by the client */
  idempotencyKey: string;

  /** Source wallet or account identifier */
  sourceAccountId: string;

  /** Destination wallet address (e.g., Binance) */
  destinationAddress: string;

  /** Amount in smallest unit (e.g., satoshi, wei) */
  amount: string;

  /** Asset/currency code (e.g., USDC, ETH) */
  assetCode: string;

  /** Optional metadata for audit */
  metadata?: Record<string, unknown>;
}

export interface TransactionRecord {
  id: string;
  idempotencyKey: string;
  sourceAccountId: string;
  destinationAddress: string;
  amount: string;
  assetCode: string;
  status: TransactionStatus;
  createdAt: Date;
  updatedAt: Date;
  failureReason?: string;
  metadata?: Record<string, unknown>;
}

export enum TransactionStatus {
  PENDING = 'pending',
  COMPLETED = 'completed',
  FAILED = 'failed',
  CANCELLED = 'cancelled',
}

export interface TransactionResult {
  success: boolean;
  transactionId?: string;
  status: TransactionStatus;
  failureReason?: string;
}

export interface TransactionEvent {
  type: 'transaction.created' | 'transaction.completed' | 'transaction.failed' | 'transaction.cancelled';
  transaction: TransactionRecord;
  timestamp: Date;
}

// ──────────────────────────────────────
// Custom Error Types
// ──────────────────────────────────────

export class PaymentValidationError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'PaymentValidationError';
  }
}

export class PaymentExecutionError extends Error {
  constructor(message: string, public readonly originalError?: unknown) {
    super(message);
    this.name = 'PaymentExecutionError';
  }
}

export class IdempotencyError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'IdempotencyError';
  }
}

// ──────────────────────────────────────
// Repository Interfaces (Dependency Inversion)
// ──────────────────────────────────────

export interface ITransactionRepository {
  create(record: Omit<TransactionRecord, 'id' | 'createdAt' | 'updatedAt'>): Promise<TransactionRecord>;
  findById(id: string): Promise<TransactionRecord | null>;
  findByKey(idempotencyKey: string): Promise<TransactionRecord | null>;
  updateStatus(id: string, status: TransactionStatus, reason?: string): Promise<void>;
}

export interface IIdempotencyCache {
  /** Retrieve cached result for a given idempotency key */
  get(key: string): Promise<TransactionResult | null>;
  /** Store result for idempotency key with TTL */
  set(key: string, result: TransactionResult, ttlSeconds?: number): Promise<void>;
  /** Remove a cached result (for cleanup) */
  delete(key: string): Promise<void>;
}

export interface IEventEmitter {
  emit(event: TransactionEvent): Promise<void>;
}

// ──────────────────────────────────────
// Logger Interface
// ──────────────────────────────────────

export interface ILogger {
  info(message: string, meta?: Record<string, unknown>): void;
  warn(message: string, meta?: Record<string, unknown>): void;
  error(message: string, meta?: Record<string, unknown>): void;
  debug(message: string, meta?: Record<string, unknown>): void;
}

// ──────────────────────────────────────
// External Payment Executor Interface
// ──────────────────────────────────────

export interface IPaymentExecutor {
  /**
   * Execute the actual payment transfer. Should throw on failure.
   * @returns true if payment was successful.
   */
  execute(request: PaymentRequest): Promise<boolean>;
}

// ──────────────────────────────────────
// Transaction Service
// ──────────────────────────────────────

export class TransactionService {
  private readonly transactionRepo: ITransactionRepository;
  private readonly idempotencyCache: IIdempotencyCache;
  private readonly eventEmitter: IEventEmitter;
  private readonly paymentExecutor: IPaymentExecutor;
  private readonly logger: ILogger;
  private readonly maxRetries: number;
  private readonly retryBaseDelayMs: number;

  constructor(
    transactionRepo: ITransactionRepository,
    idempotencyCache: IIdempotencyCache,
    eventEmitter: IEventEmitter,
    paymentExecutor: IPaymentExecutor,
    logger?: ILogger,
    maxRetries: number = 3,
    retryBaseDelayMs: number = 1000,
  ) {
    this.transactionRepo = transactionRepo;
    this.idempotencyCache = idempotencyCache;
    this.eventEmitter = eventEmitter;
    this.paymentExecutor = paymentExecutor;
    this.logger = logger ?? console as unknown as ILogger;
    this.maxRetries = maxRetries;
    this.retryBaseDelayMs = retryBaseDelayMs;
  }

  /**
   * Process a payment request with idempotency, validation, and retry logic.
   * 
   * @param request - The validated payment request.
   * @returns TransactionResult indicating final outcome.
   */
  public async processPayment(request: PaymentRequest): Promise<TransactionResult> {
    try {
      await this.validatePaymentRequest(request);

      const { idempotencyKey } = request;

      // 1. Check idempotency cache
      const cached = await this.idempotencyCache.get(idempotencyKey);
      if (cached) {
        this.logger.info(`Idempotency cache hit for key ${idempotencyKey}`, { result: cached });
        return cached;
      }

      // 2. Check repository for existing transaction
      const existing = await this.transactionRepo.findByKey(idempotencyKey);
      if (existing) {
        const existingResult = this.transactionRecordToResult(existing);
        await this.idempotencyCache.set(idempotencyKey, existingResult);
        return existingResult;
      }

      // 3. Validate business rules
      const validationError = await this.validateBusinessRules(request);
      if (validationError) {
        return await this.handleValidationFailure(request, validationError);
      }

      // 4. Create transaction record
      const newTransaction = await this.transactionRepo.create({
        idempotencyKey,
        sourceAccountId: request.sourceAccountId,
        destinationAddress: request.destinationAddress,
        amount: request.amount,
        assetCode: request.assetCode,
        status: TransactionStatus.PENDING,
        metadata: request.metadata,
      });

      this.logger.info(`Transaction created`, {
        transactionId: newTransaction.id,
        idempotencyKey,
      });

      // 5. Emit created event
      await this.emitTransactionEvent('transaction.created', newTransaction);

      // 6. Execute payment with retries
      try {
        const paymentSuccess = await this.executeWithRetries(request);

        if (paymentSuccess) {
          await this.transactionRepo.updateStatus(newTransaction.id, TransactionStatus.COMPLETED);
          const completedRecord = { ...newTransaction, status: TransactionStatus.COMPLETED };
          const result: TransactionResult = {
            success: true,
            transactionId: newTransaction.id,
            status: TransactionStatus.COMPLETED,
          };
          await this.idempotencyCache.set(idempotencyKey, result);
          await this.emitTransactionEvent('transaction.completed', completedRecord);
          return result;
        } else {
          throw new PaymentExecutionError('Payment execution failed after retries');
        }
      } catch (executionError) {
        await this.transactionRepo.updateStatus(
          newTransaction.id,
          TransactionStatus.FAILED,
          executionError instanceof Error ? executionError.message : 'Unknown error'
        );
        const failedRecord = { ...newTransaction, status: TransactionStatus.FAILED };
        const result: TransactionResult = {
          success: false,
          transactionId: newTransaction.id,
          status: TransactionStatus.FAILED,
          failureReason: executionError instanceof Error ? executionError.message : 'Unknown error',
        };
        await this.idempotencyCache.set(idempotencyKey, result);
        await this.emitTransactionEvent('transaction.failed', failedRecord);
        return result;
      }
    } catch (error) {
      this.logger.error(`Failed to process payment`, { error });
      throw new Error('Payment processing failed');
    }
  }

  private async validatePaymentRequest(request: PaymentRequest): Promise<void> {
    // Implement validation logic
  }

  private async validateBusinessRules(request: PaymentRequest): Promise<string | null> {
    // Implement business rule validation logic
    return null;
  }

  private async handleValidationFailure(request: PaymentRequest, error: string): Promise<TransactionResult> {
    this.logger.warn(`Validation failed for payment request`, { request, error });
    const result: TransactionResult = {
      success: false,
      status: TransactionStatus.FAILED,
      failureReason: error,
    };
    await this.idempotencyCache.set(request.idempotencyKey, result);
    return result;
  }

  private async executeWithRetries(request: PaymentRequest): Promise<boolean> {
    let attempts = 0;

    while (attempts < this.maxRetries) {
      try {
        const success = await this.paymentExecutor.execute(request);
        if (success) {
          return true;
        }
      } catch (error) {
        this.logger.warn(`Payment execution failed, retrying...`, { attempt: attempts + 1, error });
      }

      attempts++;
      await new Promise(resolve => setTimeout(resolve, this.retryBaseDelayMs * Math.pow(2, attempts - 1)));
    }

    return false;
  }

  private transactionRecordToResult(record: TransactionRecord): TransactionResult {
    return {
      success: record.status === TransactionStatus.COMPLETED,
      transactionId: record.id,
      status: record.status,
      failureReason: record.failureReason,
    };
  }

  private async emitTransactionEvent(type: TransactionEvent['type'], transaction: TransactionRecord): Promise<void> {
    const event: TransactionEvent = {
      type,
      transaction,
      timestamp: new Date(),
    };

    await this.eventEmitter.emit(event);
  }
}