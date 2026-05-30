// src/transactions/transactionController.ts
import { Router, Request, Response } from 'express';
import { v4 as uuidv4 } from 'uuid';
import { authenticate } from '../middleware/authenticate';
import { idempotencyMiddleware } from '../middleware/idempotency';
import { validateRequest } from '../middleware/validateRequest';
import {
  createTransaction,
  getTransactionById,
  cancelTransaction,
} from './transactionService';
import { CreateTransactionDTO, TransactionResponse } from './types';
import { logger } from '../utils/logger';

const router = Router();

/**
 * POST /transactions
 *
 * Creates a new transaction (withdrawal/payout) with idempotency guarantee.
 * The request must include a valid session token and an `idempotencyKey` header.
 * If a previous request with the same key was already processed, the cached
 * response is returned without executing the transaction again.
 *
 * @body {CreateTransactionDTO} - Transaction details (amount, destination, currency)
 * @header {string} idempotency-key - Unique key to prevent duplicate submissions
 * @returns {TransactionResponse} - Created transaction object
 * @throws 400 - Validation error
 * @throws 401 - Unauthenticated
 * @throws 409 - Conflict (duplicate request with different payload)
 */
router.post(
  '/',
  authenticate,                    // Ensures valid session and provides req.user
  idempotencyMiddleware,           // Checks/creates idempotency entry; stores original request
  validateRequest(CreateTransactionDTO), // Validates body against schema
  async (req: Request, res: Response): Promise<void> => {
    try {
      const dto: CreateTransactionDTO = req.body;
      const userId = req.user!.id;
      const idempotencyKey = req.headers['idempotency-key'] as string;

      // The idempotency middleware already deduplicates. Here we process the request.
      const result: TransactionResponse = await createTransaction(userId, dto, idempotencyKey);

      logger.info('Transaction created', {
        transactionId: result.id,
        userId,
        idempotencyKey,
        amount: dto.amount,
        currency: dto.currency,
      });

      res.status(201).json(result);
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Unknown error';
      logger.error('Transaction creation failed', { error: message });
      res.status(500).json({ error: 'Failed to create transaction' });
    }
  },
);

/**
 * GET /transactions/:id
 *
 * Retrieves the status and details of a specific transaction.
 * The user must own the transaction or have admin privileges.
 *
 * @param {string} id - Transaction ID
 * @returns {TransactionResponse} - Transaction object with current state
 * @throws 401 - Unauthenticated
 * @throws 404 - Transaction not found
 */
router.get(
  '/:id',
  authenticate,
  async (req: Request, res: Response): Promise<void> => {
    try {
      const { id } = req.params;
      const userId = req.user!.id;

      const transaction = await getTransactionById(id, userId);
      if (!transaction) {
        res.status(404).json({ error: 'Transaction not found' });
        return;
      }

      res.json(transaction);
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Unknown error';
      logger.error('Get transaction failed', { error: message });
      res.status(500).json({ error: 'Failed to retrieve transaction' });
    }
  },
);

/**
 * POST /transactions/:id/cancel
 *
 * Cancels a pending or processing transaction (order).
 * This is typically used when the buyer is unavailable or the transaction
 * has exceeded the grace period. The cancellation is recorded and the
 * underlying order cancellation logic handles state transitions and notifications.
 *
 * @param {string} id - Transaction ID to cancel
 * @returns {TransactionResponse} - Updated transaction with 'cancelled' status
 * @throws 401 - Unauthenticated
 * @throws 404 - Transaction not found
 * @throws 409 - Transaction cannot be cancelled (e.g., already completed)
 */
router.post(
  '/:id/cancel',
  authenticate,
  async (req: Request, res: Response): Promise<void> => {
    try {
      const { id } = req.params;
      const userId = req.user!.id;
      const reason = req.body.reason || 'Buyer unavailable';

      const updated = await cancelTransaction(id, userId, reason);
      if (!updated) {
        res.status(404).json({ error: 'Transaction not found or not cancellable' });
        return;
      }

      logger.info('Transaction cancelled', {
        transactionId: id,
        userId,
        reason,
      });

      res.json(updated);
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Unknown error';
      logger.error('Cancel transaction failed', { error: message });
      res.status(500).json({ error: 'Failed to cancel transaction' });
    }
  },
);

export default router;