import { Router, Request, Response, NextFunction } from 'express';
import crypto from 'crypto';

// ──────────────────────────────────────────────
// Types & Interfaces
// ──────────────────────────────────────────────

export type OrderStatus = 'pending' | 'active' | 'completed' | 'cancelled' | 'failed';

export interface Order {
  id: string;
  buyerId: string;
  sellerId: string;
  amount: number;
  asset: string;
  status: OrderStatus;
  idempotencyKey: string;
  createdAt: Date;
  updatedAt: Date;
  cancelledAt?: Date;
  cancelReason?: string;
}

export interface CreateOrderRequest {
  buyerId: string;
  sellerId: string;
  amount: number;
  asset: string;
  idempotencyKey: string;
}

export interface OrderResponse {
  success: boolean;
  data?: Order;
  error?: string;
}

// ──────────────────────────────────────────────
// In‑Memory Stores (replace with persistent DB)
// ──────────────────────────────────────────────

const orders: Map<string, Order> = new Map();
const idempotencyCache: Map<string, { status: number; body: OrderResponse }> = new Map();
const timeoutTimers: Map<string, NodeJS.Timeout> = new Map();

// ──────────────────────────────────────────────
// Idempotency Middleware
// ──────────────────────────────────────────────

function idempotencyMiddleware(req: Request, res: Response, next: NextFunction): void {
  const idempotencyKey = req.headers['idempotency-key'] as string;
  if (!idempotencyKey) {
    return next();
  }

  const cached = idempotencyCache.get(idempotencyKey);
  if (cached) {
    return res.status(cached.status).json(cached.body);
  }

  res.on('finish', () => {
    if (res.statusCode >= 200 && res.statusCode < 500) {
      idempotencyCache.set(idempotencyKey, {
        status: res.statusCode,
        body: res.locals.responseBody,
      });
      // Optional: expire cached responses after 1 hour
      setTimeout(() => idempotencyCache.delete(idempotencyKey), 60 * 60 * 1000);
    }
  });

  next();
}

// ──────────────────────────────────────────────
// Order State Machine
// ──────────────────────────────────────────────

function isValidTransition(from: OrderStatus, to: OrderStatus): boolean {
  const transitions: Record<OrderStatus, OrderStatus[]> = {
    pending:   ['active', 'cancelled'],
    active:    ['completed', 'cancelled', 'failed'],
    completed: [],
    cancelled: [],
    failed:    [],
  };
  return transitions[from]?.includes(to) ?? false;
}

// ──────────────────────────────────────────────
// Helper: create order object
// ──────────────────────────────────────────────

function createOrder(req: CreateOrderRequest): Order {
  return {
    id: crypto.randomUUID(),
    buyerId: req.buyerId,
    sellerId: req.sellerId,
    amount: req.amount,
    asset: req.asset,
    status: 'pending',
    idempotencyKey: req.idempotencyKey,
    createdAt: new Date(),
    updatedAt: new Date(),
  };
}

// ──────────────────────────────────────────────
// Timeout: auto‑cancel if buyer not active
// ──────────────────────────────────────────────

const BUYER_TIMEOUT_MS = 30_000; // 30 seconds – configurable

function startBuyerTimeout(orderId: string): void {
  const timer = setTimeout(() => {
    const order = orders.get(orderId);
    if (order && order.status === 'pending') {
      order.status = 'cancelled';
      order.cancelledAt = new Date();
      order.cancelReason = 'Buyer did not become active within timeout';
      order.updatedAt = new Date();
      console.log(`[Order ${orderId}] Auto‑cancelled due to buyer inactivity`);
      // Notify parties (webhook / event emitter) – placeholder
    }
    timeoutTimers.delete(orderId);
  }, BUYER_TIMEOUT_MS);

  timeoutTimers.set(orderId, timer);
}

function cancelBuyerTimeout(orderId: string): void {
  const timer = timeoutTimers.get(orderId);
  if (timer) {
    clearTimeout(timer);
    timeoutTimers.delete(orderId);
  }
}

// ──────────────────────────────────────────────
// Router Definition
// ──────────────────────────────────────────────

const router = Router();

// Apply idempotency middleware globally to mutations
router.use('/orders', (req: Request, _res: Response, next: NextFunction) => {
  if (['POST', 'PUT', 'PATCH', 'DELETE'].includes(req.method)) {
    return idempotencyMiddleware(req, _res, next);
  }
  next();
});

// POST /orders – create a new order
router.post('/orders', async (req: Request, res: Response, next: NextFunction) => {
  try {
    const { buyerId, sellerId, amount, asset, idempotencyKey } = req.body;

    // Validation
    if (!buyerId || !sellerId || !amount || !asset || !idempotencyKey) {
      const response: OrderResponse = { success: false, error: 'Missing required fields' };
      res.locals.responseBody = response;
      return res.status(400).json(response);
    }

    // Check if order with same idempotency key already exists (safety)
    const existing = Array.from(orders.values()).find(o => o.idempotencyKey === idempotencyKey);
    if (existing) {
      const response: OrderResponse = { success: true, data: existing };
      res.locals.responseBody = response;
      return res.status(200).json(response);
    }

    const order = createOrder({ buyerId, sellerId, amount, asset, idempotencyKey });
    orders.set(order.id, order);

    // Start timeout for buyer activity
    startBuyerTimeout(order.id);

    const response: OrderResponse = { success: true, data: order };
    res.locals.responseBody = response;
    return res.status(201).json(response);
  } catch (err) {
    next(err);
  }
});

// GET /orders – list orders (with optional status filter)
router.get('/orders', async (req: Request, res: Response, next: NextFunction) => {
  try {
    const statusFilter = req.query.status as OrderStatus | undefined;
    let result: Order[] = Array.from(orders.values());
    if (statusFilter) {
      result = result.filter(o => o.status === statusFilter);
    }
    res.locals.responseBody = { success: true, data: result };
    return res.status(200).json({ success: true, data: result });
  } catch (err) {
    next(err);
  }
});

// GET /orders/:id – get order by ID
router.get('/orders/:id', async (req: Request, res: Response, next: NextFunction) => {
  try {
    const order = orders.get(req.params.id);
    if (!order) {
      return res.status(404).json({ success: false, error: 'Order not found' });
    }
    res.locals.responseBody = { success: true, data: order };
    return res.status(200).json({ success: true, data: order });
  } catch (err) {
    next(err);
  }
});

// POST /orders/:id/cancel – cancel an order manually
router.post('/orders/:id/cancel', async (req: Request, res: Response, next: NextFunction) => {
  try {
    const order = orders.get(req.params.id);
    if (!order) {
      return res.status(404).json({ success: false, error: 'Order not found' });
    }

    if (order.status === 'completed' || order.status === 'cancelled' || order.status === 'failed') {
      return res.status(409).json({ success: false, error: `Order already in terminal state: ${order.status}` });
    }

    if (!isValidTransition(order.status, 'cancelled')) {
      return res.status(409).json({ success: false, error: `Cannot cancel order in status ${order.status}` });
    }

    // Cancel timeout
    cancelBuyerTimeout(order.id);

    order.status = 'cancelled';
    order.cancelledAt = new Date();
    order.cancelReason = req.body.reason || 'Manual cancellation';
    order.updatedAt = new Date();

    // Notify parties – placeholder

    const response: OrderResponse = { success: true, data: order };
    res.locals.responseBody = response;
    return res.status(200).json(response);
  } catch (err) {
    next(err);
  }
});

// POST /orders/:id/activate – simulate buyer becoming active (for testing)
router.post('/orders/:id/activate', async (req: Request, res: Response, next: NextFunction) => {
  try {
    const order = orders.get(req.params.id);
    if (!order) {
      return res.status(404).json({ success: false, error: 'Order not found' });
    }

    if (order.status !== 'pending') {
      return res.status(409).json({ success: false, error: `Order not in pending state: ${order.status}` });
    }

    cancelBuyerTimeout(order.id);
    order.status = 'active';
    order.updatedAt = new Date();

    const response: OrderResponse = { success: true, data: order };
    res.locals.responseBody = response;
    return res.status(200).json(response);
  } catch (err) {
    next(err);
  }
});

// ──────────────────────────────────────────────
// Global Error Handler
// ──────────────────────────────────────────────

router.use((err: Error, _req: Request, res: Response, _next: NextFunction) => {
  console.error('[OrderRouter] Unhandled error:', err);
  res.status(500).json({ success: false, error: 'Internal server error' });
});

export default router;