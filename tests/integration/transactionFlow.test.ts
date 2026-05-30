// tests/integration/transactionFlow.test.ts
// Integration tests for ME-Hub Reliability & Idempotency Architecture

import { AuthService } from '../../src/services/auth.service';
import { TransactionService } from '../../src/services/transaction.service';
import { OrderTrackingService } from '../../src/services/order-tracking.service';
import { IdempotencyMiddleware } from '../../src/middleware/idempotency.middleware';
import { v4 as uuidv4 } from 'uuid';

// ---------------------------------------------------------------------
// Test utilities
// ---------------------------------------------------------------------

/**
 * Simulates an unreliable external authentication provider (Me Pass)
 * that fails on first attempt and succeeds on subsequent attempts.
 */
class UnreliableMePassProvider {
  private attemptCount = 0;

  async authenticate(credentials: { username: string; password: string }): Promise<{ token: string } | null> {
    this.attemptCount++;
    // First call fails, second and third succeed
    if (this.attemptCount === 1) {
      return null;
    }
    return { token: `session-${uuidv4()}` };
  }

  reset() {
    this.attemptCount = 0;
  }
}

/**
 * Simulates a buyer that times out after a configurable delay.
 */
class TimeoutBuyerSimulator {
  private timeout: number;

  constructor(timeoutMs = 200) {
    this.timeout = timeoutMs;
  }

  async confirmOrder(orderId: string, delayMs: number): Promise<boolean> {
    await new Promise((resolve) => setTimeout(resolve, delayMs));
    return delayMs <= this.timeout;
  }
}

// ---------------------------------------------------------------------
// Integration test suite
// ---------------------------------------------------------------------

describe('Transaction Flow - End-to-End Idempotency', () => {
  let authService: AuthService;
  let transactionService: TransactionService;
  let orderTrackingService: OrderTrackingService;
  let idempotencyMiddleware: IdempotencyMiddleware;
  let mePassProvider: UnreliableMePassProvider;

  const userCredentials = { username: 'test_user', password: 'test_pass' };
  const payoutRequest = {
    amount: 100,
    currency: 'USD',
    recipientAddress: '0x1234567890abcdef',
  };

  beforeEach(() => {
    // Reset all services to clean state
    mePassProvider = new UnreliableMePassProvider();
    authService = new AuthService(mePassProvider);
    transactionService = new TransactionService();
    orderTrackingService = new OrderTrackingService({ buyerTimeoutMs: 150 });
    idempotencyMiddleware = new IdempotencyMiddleware();
  });

  // -----------------------------------------------------------------
  // Test 1: Authentication with retry and fallback to QR
  // -----------------------------------------------------------------
  it('should successfully authenticate using retry logic (Me Pass -> QR fallback)', async () => {
    // Attempt 1: Me Pass fails
    const attempt1 = await authService.authenticateWithMePass(userCredentials);
    expect(attempt1.success).toBe(false);
    expect(attempt1.fallbackMethod).toBe('qr');

    // Attempt 2: QR fallback with retry
    const attempt2 = await authService.authenticateWithQr();
    expect(attempt2.success).toBe(true);
    expect(attempt2.sessionToken).toBeDefined();
    expect(attempt2.idempotencyKey).toBeDefined();
  });

  // -----------------------------------------------------------------
  // Test 2: Idempotency prevents duplicate transaction submissions
  // -----------------------------------------------------------------
  it('should return same response for duplicate transaction requests with same idempotency key', async () => {
    // First: authenticate to get session and idempotency key
    await authService.authenticateWithMePass(userCredentials);
    const authResult = await authService.authenticateWithQr();
    const idempotencyKey = authResult.idempotencyKey;

    // First transaction request
    const firstResponse = await idempotencyMiddleware.execute(
      { idempotencyKey, body: payoutRequest },
      async () => transactionService.processPayout(payoutRequest),
    );

    expect(firstResponse.status).toBe('success');
    expect(firstResponse.transactionId).toBeDefined();
    const firstTransactionId = firstResponse.transactionId;

    // Second request with same idempotency key
    const secondResponse = await idempotencyMiddleware.execute(
      { idempotencyKey, body: payoutRequest },
      async () => transactionService.processPayout(payoutRequest),
    );

    // Should return the exact same response (cached), not a new transaction
    expect(secondResponse.status).toBe('success');
    expect(secondResponse.transactionId).toBe(firstTransactionId);
  });

  // -----------------------------------------------------------------
  // Test 3: Idempotent transaction processing (end-to-end)
  // -----------------------------------------------------------------
  it('should complete a full idempotent transaction flow', async () => {
    // Authenticate
    const authSession = await authService.loginWithFallback(userCredentials);
    expect(authSession.token).toBeDefined();

    // Send payout request with idempotency key
    const idempotencyKey = authSession.idempotencyKey;
    const transactionResult = await transactionService.createTransaction({
      ...payoutRequest,
      idempotencyKey,
      sessionToken: authSession.token,
    });

    expect(transactionResult.status).toBe('completed');
    expect(transactionResult.amount).toBe(payoutRequest.amount);
    expect(transactionResult.recipient).toBe(payoutRequest.recipientAddress);
    expect(transactionResult.duplicate).toBe(false);

    // Verify idempotency – resubmitting returns same result
    const duplicateResult = await transactionService.createTransaction({
      ...payoutRequest,
      idempotencyKey,
      sessionToken: authSession.token,
    });
    expect(duplicateResult.transactionId).toBe(transactionResult.transactionId);
    expect(duplicateResult.duplicate).toBe(true); // marked as duplicate
  });

  // -----------------------------------------------------------------
  // Test 4: Automatic order cancellation when buyer is unavailable
  // -----------------------------------------------------------------
  it('should cancel order if buyer does not confirm within timeout', async () => {
    // Create an order
    const order = await orderTrackingService.createOrder({
      buyerId: 'buyer-1',
      amount: 50,
      timeoutMs: 150,
    });

    // Simulate a buyer that takes too long (220ms > 150ms timeout)
    const buyerSimulator = new TimeoutBuyerSimulator(150);
    const buyerConfirmed = await buyerSimulator.confirmOrder(order.id, 220);

    // Buyer didn't confirm in time
    expect(buyerConfirmed).toBe(false);

    // Wait for cancellation timer
    await orderTrackingService.checkForTimeout(order.id);

    const cancelledOrder = await orderTrackingService.getOrder(order.id);
    expect(cancelledOrder.status).toBe('cancelled');
    expect(cancelledOrder.cancellationReason).toContain('buyer unavailable');
  });

  // -----------------------------------------------------------------
  // Test 5: Full happy path – login, transaction, no duplicate
  // -----------------------------------------------------------------
  it('should complete a successful transaction without duplicates when allowed', async () => {
    // Login
    const session = await authService.loginWithFallback(userCredentials);
    const key1 = uuidv4();
    const key2 = uuidv4();

    // First payout
    const tx1 = await transactionService.createTransaction({
      idempotencyKey: key1,
      sessionToken: session.token,
      amount: 200,
      currency: 'EUR',
      recipientAddress: '0xfeedbeef',
    });
    expect(tx1.status).toBe('completed');

    // Second payout with different key (allowed)
    const tx2 = await transactionService.createTransaction({
      idempotencyKey: key2,
      sessionToken: session.token,
      amount: 300,
      currency: 'EUR',
      recipientAddress: '0xcafebabe',
    });
    expect(tx2.status).toBe('completed');
    expect(tx2.transactionId).not.toBe(tx1.transactionId); // different transactions
  });

  // -----------------------------------------------------------------
  // Test 6: Idempotency of order status transitions (cancel twice)
  // -----------------------------------------------------------------
  it('should maintain idempotency for order cancellation', async () => {
    const order = await orderTrackingService.createOrder({
      buyerId: 'buyer-2',
      amount: 10,
    });

    // Cancel once
    await orderTrackingService.cancelOrder(order.id, 'buyer timeout');

    // Cancel again – should be idempotent
    await expect(orderTrackingService.cancelOrder(order.id, 'buyer timeout')).resolves.not.toThrow();

    // Verify status
    const updatedOrder = await orderTrackingService.getOrder(order.id);
    expect(updatedOrder.status).toBe('cancelled');
    // Only one cancellation event should be emitted
    expect(updatedOrder.cancellationCount ?? 1).toBe(1);
  });
});