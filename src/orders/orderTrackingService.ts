// src/orders/orderTrackingService.ts

/**
 * Order status lifecycle:
 *   PENDING → ACTIVE (buyer available) → COMPLETED
 *   PENDING → CANCELLED (timeout or explicit cancellation)
 *   ACTIVE  → CANCELLED (if buyer becomes unavailable or manual cancel)
 */
export enum OrderStatus {
  PENDING = 'PENDING',
  ACTIVE = 'ACTIVE',
  COMPLETED = 'COMPLETED',
  CANCELLED = 'CANCELLED',
}

export interface Order {
  id: string;
  status: OrderStatus;
  buyerId: string;
  createdAt: Date;
  updatedAt: Date;
  /** If not null, the order should be cancelled after this timestamp */
  timeoutAt: Date | null;
  cancellationReason: string | null;
}

export interface OrderTrackingConfig {
  /** Grace period in milliseconds after which an unresponsive order is cancelled (default 5 minutes) */
  buyerTimeoutMs: number;
  /** Interval in ms for checking timeouts (default 10s) */
  timeoutCheckIntervalMs: number;
}

const DEFAULT_CONFIG: OrderTrackingConfig = {
  buyerTimeoutMs: 5 * 60 * 1000, // 5 minutes
  timeoutCheckIntervalMs: 10_000, // 10 seconds
};

/**
 * OrderTrackingService – manages order lifecycle, buyer availability detection,
 * and automatic cancellation of stale orders.
 */
export class OrderTrackingService {
  private orders: Map<string, Order> = new Map();
  private config: OrderTrackingConfig;
  private timeoutTimer: NodeJS.Timeout | null = null;

  constructor(config?: Partial<OrderTrackingConfig>) {
    this.config = { ...DEFAULT_CONFIG, ...config };
  }

  // ─────────────────────────────────────────────────────────────────
  // Public API
  // ─────────────────────────────────────────────────────────────────

  /**
   * Creates a new order in PENDING status and starts the buyer availability timer.
   * @param orderId   Unique identifier for the order (e.g. from transaction service)
   * @param buyerId   Identifier of the buyer
   * @returns The newly created Order
   */
  createOrder(orderId: string, buyerId: string): Order {
    if (this.orders.has(orderId)) {
      throw new Error(`Order with id '${orderId}' already exists`);
    }

    const now = new Date();
    const timeoutAt = new Date(now.getTime() + this.config.buyerTimeoutMs);

    const order: Order = {
      id: orderId,
      status: OrderStatus.PENDING,
      buyerId,
      createdAt: now,
      updatedAt: now,
      timeoutAt,
      cancellationReason: null,
    };

    this.orders.set(orderId, order);
    this.ensureTimeoutCheckerRunning();
    console.log(`[OrderTracking] Created order ${orderId} for buyer ${buyerId}, timeout at ${timeoutAt.toISOString()}`);
    return order;
  }

  /**
   * Marks a buyer as available, moving the order to ACTIVE status and clearing the timeout.
   * @param orderId ID of the order
   * @returns The updated Order
   */
  markBuyerAvailable(orderId: string): Order {
    const order = this.getOrderOrThrow(orderId);
    if (order.status !== OrderStatus.PENDING) {
      throw new Error(`Cannot mark buyer available for order ${orderId} in status ${order.status}`);
    }

    order.status = OrderStatus.ACTIVE;
    order.timeoutAt = null;
    order.updatedAt = new Date();
    console.log(`[OrderTracking] Buyer available for order ${orderId} → ACTIVE`);

    return order;
  }

  /**
   * Completes an active order.
   * @param orderId ID of the order
   * @returns The updated Order
   */
  completeOrder(orderId: string): Order {
    const order = this.getOrderOrThrow(orderId);
    if (order.status !== OrderStatus.ACTIVE) {
      throw new Error(`Cannot complete order ${orderId} in status ${order.status}`);
    }

    order.status = OrderStatus.COMPLETED;
    order.updatedAt = new Date();
    console.log(`[OrderTracking] Order ${orderId} completed`);
    return order;
  }

  /**
   * Cancels an order with a reason.
   * @param orderId ID of the order
   * @param reason  Human-readable cancellation reason
   * @returns The updated Order
   */
  cancelOrder(orderId: string, reason: string): Order {
    const order = this.getOrderOrThrow(orderId);
    if (order.status === OrderStatus.COMPLETED) {
      throw new Error(`Cannot cancel completed order ${orderId}`);
    }

    order.status = OrderStatus.CANCELLED;
    order.timeoutAt = null;
    order.cancellationReason = reason;
    order.updatedAt = new Date();
    console.log(`[OrderTracking] Order ${orderId} cancelled: ${reason}`);
    this.notifyCancellation(order);
    return order;
  }

  /**
   * Retrieves an order by ID.
   */
  getOrder(orderId: string): Order | undefined {
    return this.orders.get(orderId);
  }

  /**
   * Returns all orders currently in a given status.
   */
  getOrdersByStatus(status: OrderStatus): Order[] {
    const result: Order[] = [];
    for (const order of this.orders.values()) {
      if (order.status === status) result.push(order);
    }
    return result;
  }

  /**
   * Returns all pending orders (those waiting for buyer availability).
   */
  getPendingOrders(): Order[] {
    return this.getOrdersByStatus(OrderStatus.PENDING);
  }

  /**
   * Stops the background timeout checker if running.
   */
  shutdown(): void {
    if (this.timeoutTimer !== null) {
      clearInterval(this.timeoutTimer);
      this.timeoutTimer = null;
      console.log('[OrderTracking] Timeout checker stopped');
    }
  }

  // ─────────────────────────────────────────────────────────────────
  // Private helpers
  // ─────────────────────────────────────────────────────────────────

  private getOrderOrThrow(orderId: string): Order {
    const order = this.orders.get(orderId);
    if (!order) {
      throw new Error(`Order not found: '${orderId}'`);
    }
    return order;
  }

  private ensureTimeoutCheckerRunning(): void {
    if (this.timeoutTimer !== null) return;

    this.timeoutTimer = setInterval(() => {
      this.checkTimeouts();
    }, this.config.timeoutCheckIntervalMs);

    // Prevent the timer from keeping the process alive if nothing else is running
    if (this.timeoutTimer && typeof this.timeoutTimer.unref === 'function') {
      this.timeoutTimer.unref();
    }

    console.log('[OrderTracking] Timeout checker started');
  }

  /**
   * Finds all pending orders whose timeout has passed and cancels them.
   */
  private checkTimeouts(): void {
    const now = new Date();
    for (const order of this.orders.values()) {
      if (order.status === OrderStatus.PENDING && order.timeoutAt && now >= order.timeoutAt) {
        this.cancelOrder(order.id, 'Buyer did not become available within the grace period');
      }
    }
  }

  /**
   * Placeholder for notifying relevant parties (e.g., buyer, seller, transaction service).
   * In production this would emit an event or call a notification service.
   */
  private notifyCancellation(order: Order): void {
    // TODO: integrate with real notification system
    console.log(`[OrderTracking] NOTIFY: Order ${order.id} cancelled. Reason: ${order.cancellationReason}`);
  }
}