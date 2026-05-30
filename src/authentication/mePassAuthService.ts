typescript
* Attempts QR code authentication as a fallback.
*
* @param request - The authentication request.
* @returns The result of the authentication attempt.
*/
private async attemptQrCode(request: MePassAuthRequest): Promise<MePassAuthResult> {
  try {
    this.logger.debug('Sending QR code authentication request', {
      userId: request.userId,
    });

    const response = await this.httpClient.post('/auth/qrcode', {
      userId: request.userId,
    }, { timeout: TIMEOUT_MS }); // Add timeout option

    // Validate response shape
    if (!response || typeof response.sessionToken !== 'string' || !response.expiresAt) {
      throw new RetryableAuthError('Invalid response structure from QR code endpoint');
    }

    return {
      success: true,
      sessionToken: response.sessionToken,
      expiresAt: response.expiresAt,
    };
  } catch (error: unknown) {
    const normalizedError = this.normalizeError(error);
    this.logger.warn('QR code attempt failed', {
      userId: request.userId,
      error: normalizedError,
    });

    // Re‑throw retryable errors; wrap non‑retryable as failure result
    if (error instanceof RetryableAuthError) {
      throw error;
    }
    if (error instanceof NonRetryableAuthError) {
      return {
        success: false,
        error: normalizedError,
      };
    }
    // Unknown errors (network, timeout, etc.) are treated as retryable
    throw new RetryableAuthError(normalizedError);
  }
}

/**
* Normalizes an unknown error into a structured AuthError.
*
* @param error - The error to normalize.
* @returns A structured AuthError object.
*/
private normalizeError(error: unknown): AuthError {
  if (error instanceof Error) {
    return {
      code: error instanceof RetryableAuthError || error instanceof NonRetryableAuthError ? error.code : 'UNKNOWN_ERROR',
      message: error.message,
      details: error.stack,
    };
  }
  return {
    code: 'UNKNOWN_ERROR',
    message: 'An unknown error occurred',
    details: String(error),
  };
}

/**
* Attempts an operation with retry logic and exponential backoff.
*
* @param operation - The operation to attempt.
* @param options - Retry options including maxRetries, operationName, and idempotencyKey.
* @returns The result of the operation.
* @throws {AuthError} If all retries fail.
*/
private async attemptWithRetry<T>(
  operation: () => Promise<T>,
  options: { maxRetries: number; operationName: string; idempotencyKey: string },
): Promise<T> {
  let attempt = 0;
  let lastError: AuthError | null = null;

  while (attempt <= options.maxRetries) {
    try {
      return await operation();
    } catch (error: unknown) {
      lastError = this.normalizeError(error);
      this.logger.warn(`Attempt ${attempt} of ${options.operationName} failed`, {
        idempotencyKey: options.idempotencyKey,
        error: lastError,
      });

      if (attempt === options.maxRetries || error instanceof NonRetryableAuthError) {
        break;
      }

      const backoff = Math.min(BASE_BACKOFF_MS * Math.pow(2, attempt), MAX_BACKOFF_MS);
      await new Promise((resolve) => setTimeout(resolve, backoff));
      attempt++;
    }
  }

  throw lastError!;
}