typescript
typescript
import { Router, Request, Response, NextFunction } from 'express';
import rateLimit from 'express-rate-limit';
import { v4 as uuidv4 } from 'uuid';
import winston from 'winston';
import { body, validationResult } from 'express-validator';
import { AbortController } from 'node-abort-controller'; // or globalThis if Node >=16

// ---------------------------------------------------------------------------
// Constants & Types
// ---------------------------------------------------------------------------

const AUTH_METHODS = ['mepass', 'qr'] as const;
type AuthMethod = (typeof AUTH_METHODS)[number];

interface MePassCredentials {
  readonly token?: string;
  readonly password?: string;
  readonly [key: string]: unknown;
}

interface MePassAuthResult {
  readonly success: boolean;
  readonly userId?: string;
  readonly error?: string;
}

interface QRChallenge {
  readonly challengeId: string;
  readonly qrDataUrl: string;
  readonly expiresAt: Date;
}

interface QRVerificationResult {
  readonly success: boolean;
  readonly userId?: string;
}

interface SessionResult {
  readonly sessionToken: string;
  readonly expiresIn: number;
}

interface ServiceInterfaces {
  mePassService: {
    authenticate(credentials: MePassCredentials, context?: object): Promise<MePassAuthResult>;
  };
  qrCodeService: {
    generateChallenge(): Promise<QRChallenge>;
    verifyChallenge(challengeId: string, response: string): Promise<QRVerificationResult>;
  };
  sessionService: {
    createSession(userId: string, idempotencyKey: string): Promise<SessionResult>;
  };
}

// In-memory store for idempotency (replace with Redis/DB in production)
const idempotencyStore = new Map<string, { sessionToken: string; expiresAt: number }>();

// ---------------------------------------------------------------------------
// Logger Configuration
// ---------------------------------------------------------------------------

const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.errors({ stack: true }),
    winston.format.json()
  ),
  defaultMeta: { service: 'auth-router' },
  transports: [
    new winston.transports.Console({
      format: process.env.NODE_ENV === 'production'
        ? winston.format.json()
        : winston.format.combine(
            winston.format.colorize(),
            winston.format.simple()
          ),
    }),
    // Uncomment for file logging in production
    // new winston.transports.File({ filename: 'logs/auth.log', level: 'warn' }),
  ],
});

// ---------------------------------------------------------------------------
// Custom Error Classes
// ---------------------------------------------------------------------------

class AuthValidationError extends Error {
  public readonly statusCode = 400;
  public readonly details: string[];

  constructor(message: string, details: string[] = []) {
    super(message);
    this.name = 'AuthValidationError';
    this.details = details;
  }
}

class AuthUnauthorizedError extends Error {
  public readonly statusCode = 401;

  constructor(message: string = 'Authentication failed') {
    super(message);
    this.name = 'AuthUnauthorizedError';
  }
}

// ---------------------------------------------------------------------------
// Helper: Create rate limiter
// ---------------------------------------------------------------------------

function createLoginLimiter(): rateLimit.RateLimitRequestHandler {
  return rateLimit({
    windowMs: 15 * 60 * 1000,          // 15 minutes
    max: 10,                            // limit each IP to 10 requests per windowMs
    message: { error: 'Too many login attempts, please try again later.' },
    standardHeaders: true,
    legacyHeaders: false,
    keyGenerator: (req: Request) => {
      // Use X-Forwarded-For if behind proxy, otherwise IP from connection
      return (req.headers['x-forwarded-for'] as string)?.split(',')[0]?.trim()
        || req.ip
        || req.socket.remoteAddress
        || 'unknown';
    },
  });
}

// ---------------------------------------------------------------------------
// Helper: Timeout wrapper for async functions
// ---------------------------------------------------------------------------

async function withTimeout<T>(promise: Promise<T>, timeout: number): Promise<T> {
  const controller = new AbortController();
  const id = setTimeout(() => controller.abort(), timeout);

  try {
    return await promise;
  } finally {
    clearTimeout(id);
  }
}

// ---------------------------------------------------------------------------
// Router Factory
// ---------------------------------------------------------------------------

/**
 * Creates an Express router for authentication endpoints.
 *
 * @param services - Object containing MePass, QR code, and session service implementations
 * @returns Configured Express Router
 */
export function createAuthRouter(services: ServiceInterfaces): Router {
  const router = Router();
  const loginLimiter = createLoginLimiter();

  // Apply rate limiter to all routes on this router
  router.use(loginLimiter);

  // -----------------------------------------------------------------------
  // POST /auth/login
  // -----------------------------------------------------------------------

  /**
   * POST /auth/login
   *
   * Authenticates a user via MePass or QR code. Supports fallback from MePass to QR.
   *
   * Request body (application/json):
   * {
   *   method: 'mepass' | 'qr',       // default: 'mepass'
   *   credentials?: object,           // required if method='mepass'
   *   qrChallengeId?: string,         // required if method='qr' and verifying
   *   qrResponse?: string             // required if method='qr' and verifying
   * }
   *
   * Success response (200):
   * {
   *   sessionToken: string,
   *   idempotencyKey: string,
   *   method: 'mepass' | 'qr',
   *   expiresIn: number,
   *   qrChallenge?: { challengeId, qrDataUrl, expiresAt } // only on first QR step
   * }
   *
   * Error responses:
   *   400 - Validation errors
   *   401 - Authentication failure
   *   429 - Rate limited
   *   500 - Internal server error
   */
  router.post(
    '/login',
    // Input validation middleware
    [
      body('method')
        .optional()
        .isIn(AUTH_METHODS)
        .withMessage(`Method must be one of: ${AUTH_METHODS.join(', ')}`),
      body('credentials')
        .if(body('method').equals('mepass').or(body('method').not().exists()))
        .exists({ checkFalsy: true })
        .withMessage('Credentials are required for MePass authentication')
        .isObject()
        .withMessage('Credentials must be an object'),
      body('qrChallengeId')
        .if(body('method').equals('qr'))
        .optional()
        .isString()
        .withMessage('qrChallengeId must be a string'),
      body('qrResponse')
        .if(body('method').equals('qr'))
        .optional()
        .isString()
        .withMessage('qrResponse must be a string'),
    ],
    async (req: Request, res: Response, next: NextFunction) => {
      const requestId = uuidv4();
      const ip = req.ip || req.socket.remoteAddress || 'unknown';

      logger.info('Login request received', {
        requestId,
        ip,
        method: req.body?.method,
        userAgent: req.headers['user-agent'],
      });

      // Handle validation errors
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        const details = errors.array().map((e) => e.msg);
        logger.warn('Login request validation failed', { requestId, details });
        return res.status(400).json({
          error: 'Invalid request',
          details,
          requestId,
        });
      }

      // Extract validated fields
      const method: AuthMethod = req.body.method || 'mepass';
      const credentials: MePassCredentials | undefined = req.body.credentials;
      const qrChallengeId: string | undefined = req.body.qrChallengeId;
      const qrResponse: string | undefined = req.body.qrResponse;

      try {
        let userId: string | null = null;
        let authMethod: AuthMethod = method;

        // --- MePass authentication attempt ---
        if (method === 'mepass') {
          logger.info('Attempting MePass authentication', { requestId });
          const mePassResult = await withTimeout(
            services.mePassService.authenticate(credentials!, { requestId }),
            5000 // 5 seconds timeout
          );

          if (!mePassResult.success) {
            throw new AuthUnauthorizedError(mePassResult.error || 'MePass authentication failed');
          }

          userId = mePassResult.userId;
        }

        // --- QR code authentication attempt ---
        if (method === 'qr') {
          if (!qrChallengeId && !qrResponse) {
            logger.info('Generating QR challenge', { requestId });
            const qrChallenge: QRChallenge = await services.qrCodeService.generateChallenge();
            return res.status(200).json({
              method,
              qrChallenge,
              requestId,
            });
          }

          if (qrChallengeId && qrResponse) {
            logger.info('Verifying QR response', { requestId, qrChallengeId });
            const qrVerificationResult = await withTimeout(
              services.qrCodeService.verifyChallenge(qrChallengeId, qrResponse),
              5000 // 5 seconds timeout
            );

            if (!qrVerificationResult.success) {
              throw new AuthUnauthorizedError('QR code verification failed');
            }

            userId = qrVerificationResult.userId;
          }
        }

        if (!userId) {
          throw new AuthUnauthorizedError('User ID not found after authentication');
        }

        // --- Create session ---
        const idempotencyKey = uuidv4();
        const sessionResult: SessionResult = await services.sessionService.createSession(userId, idempotencyKey);

        logger.info('Session created', { requestId, userId, sessionToken: sessionResult.sessionToken });

        res.status(200).json({
          sessionToken: sessionResult.sessionToken,
          idempotencyKey,
          method,
          expiresIn: sessionResult.expiresIn,
          requestId,
        });
      } catch (error) {
        if (error instanceof AuthValidationError || error instanceof AuthUnauthorizedError) {
          logger.warn('Authentication failed', { requestId, error: error.message });
          return res.status(error.statusCode).json({
            error: error.message,
            details: error.details || [],
            requestId,
          });
        }

        logger.error('Internal server error', { requestId, error });
        res.status(500).json({
          error: 'Internal server error',
          requestId,
        });
      }
    }
  );

  return router;
}