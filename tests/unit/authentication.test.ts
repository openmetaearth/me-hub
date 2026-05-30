// tests/unit/authentication.test.ts

import { MePassAuthService } from '../../src/services/mePassAuthService';
import { QRCodeAuthService } from '../../src/services/qrCodeAuthService';
import { IdempotencyService } from '../../src/services/idempotencyService';

// -----------------------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------------------

const mockHttpClient = {
  post: jest.fn(),
};

const mockSessionStore = {
  set: jest.fn(),
  get: jest.fn(),
};

const mockIdempotencyService = {
  generateKey: jest.fn(),
  checkCache: jest.fn(),
  cacheResult: jest.fn(),
} as unknown as jest.Mocked<IdempotencyService>;

jest.mock('../../src/services/idempotencyService', () => ({
  IdempotencyService: jest.fn().mockImplementation(() => mockIdempotencyService),
}));

// -----------------------------------------------------------------------------
// Test suites
// -----------------------------------------------------------------------------

describe('MePassAuthService', () => {
  let mePassAuth: MePassAuthService;

  beforeEach(() => {
    jest.clearAllMocks();
    mePassAuth = new MePassAuthService(mockHttpClient, mockSessionStore);
  });

  describe('login', () => {
    it('should return session token on successful MePass authentication', async () => {
      const credentials = { username: 'testuser', password: 'correct-password' };
      const expectedToken = 'session-token-123';
      const expectedIdempotencyKey = 'idem-key-456';

      mockIdempotencyService.generateKey.mockReturnValue(expectedIdempotencyKey);
      mockHttpClient.post.mockResolvedValueOnce({
        data: { token: expectedToken, idempotencyKey: expectedIdempotencyKey },
      });

      const result = await mePassAuth.login(credentials);

      expect(result).toEqual({
        success: true,
        token: expectedToken,
        idempotencyKey: expectedIdempotencyKey,
      });
      expect(mockHttpClient.post).toHaveBeenCalledWith(
        '/auth/login',
        credentials,
        expect.any(Object)
      );
      expect(mockSessionStore.set).toHaveBeenCalledWith(
        expect.stringContaining('mePass'),
        expectedToken
      );
    });

    it('should throw MePassAuthError on invalid credentials', async () => {
      mockHttpClient.post.mockRejectedValueOnce(new Error('Invalid username or password'));

      await expect(mePassAuth.login({ username: 'bad', password: 'wrong' })).rejects.toThrow(
        'MePass authentication failed'
      );
    });

    it('should apply exponential backoff on consecutive failures', async () => {
      const credentials = { username: 'user', password: 'pass' };
      mockHttpClient.post
        .mockRejectedValueOnce(new Error('Network error'))
        .mockRejectedValueOnce(new Error('Network error'))
        .mockResolvedValueOnce({ data: { token: 'token-after-retry', idempotencyKey: 'key' } });

      const start = Date.now();
      const result = await mePassAuth.login(credentials, { maxRetries: 3 });
      const duration = Date.now() - start;

      expect(result.success).toBe(true);
      expect(mockHttpClient.post).toHaveBeenCalledTimes(3);
      // Backoff: first retry ~100ms, second ~200ms => total >= 300ms
      expect(duration).toBeGreaterThanOrEqual(300);
    });
  });

  describe('fallbackToQR', () => {
    it('should delegate to QRCodeAuthService when MePass fails', async () => {
      const credentials = { username: 'testuser', password: 'correct' };
      const qrMock = jest.spyOn(QRCodeAuthService.prototype, 'authenticate');
      const expectedQrResult = { success: true, qrData: 'qr-data-string' };

      mockHttpClient.post.mockRejectedValueOnce(new Error('MePass unavailable'));
      qrMock.mockResolvedValue(expectedQrResult);

      const result = await mePassAuth.fallbackToQR(credentials);

      expect(result).toEqual(expectedQrResult);
      expect(qrMock).toHaveBeenCalledWith(credentials);
    });
  });
});

describe('QRCodeAuthService', () => {
  let qrAuth: QRCodeAuthService;

  beforeEach(() => {
    jest.clearAllMocks();
    qrAuth = new QRCodeAuthService(mockHttpClient);
  });

  describe('authenticate', () => {
    it('should return qr data on successful authentication', async () => {
      const credentials = { username: 'qruser' };
      const expectedQrData = { qrCode: 'data:image/png;base64,...', sessionId: 'sess-001' };

      mockHttpClient.post.mockResolvedValueOnce({ data: expectedQrData });

      const result = await qrAuth.authenticate(credentials);

      expect(result).toEqual({ success: true, qrData: expectedQrData });
    });

    it('should throw QRCodeAuthError when QR code generation fails', async () => {
      mockHttpClient.post.mockRejectedValueOnce(new Error('QR service timeout'));

      await expect(qrAuth.authenticate({})).rejects.toThrow('QR code authentication failed');
    });

    it('should retry up to maxRetries times on network errors', async () => {
      mockHttpClient.post
        .mockRejectedValueOnce(new Error('Network error'))
        .mockRejectedValueOnce(new Error('Network error'))
        .mockResolvedValueOnce({ data: { qrCode: 'retry-success', sessionId: 'sess-retry' } });

      const result = await qrAuth.authenticate({}, { maxRetries: 3 });

      expect(result.success).toBe(true);
      expect(mockHttpClient.post).toHaveBeenCalledTimes(3);
    });
  });

  describe('pollForConfirmation', () => {
    it('should resolve when QR scan is confirmed', async () => {
      const sessionId = 'sess-001';
      mockHttpClient.post
        .mockResolvedValueOnce({ data: { status: 'pending' } })
        .mockResolvedValueOnce({ data: { status: 'pending' } })
        .mockResolvedValueOnce({ data: { status: 'confirmed', token: 'final-token' } });

      const result = await qrAuth.pollForConfirmation(sessionId);

      expect(result).toEqual({ success: true, token: 'final-token' });
      expect(mockHttpClient.post).toHaveBeenCalledTimes(3);
    });

    it('should timeout after configured duration', async () => {
      jest.useFakeTimers();
      mockHttpClient.post.mockResolvedValue({ data: { status: 'pending' } });

      const promise = qrAuth.pollForConfirmation('sess-timeout', { timeoutMs: 3000 });
      jest.advanceTimersByTime(3000);

      await expect(promise).rejects.toThrow('QR confirmation timed out');
      jest.useRealTimers();
    });
  });
});

// -----------------------------------------------------------------------------
// Integration: fallback flow with idempotency
// -----------------------------------------------------------------------------

describe('Authentication Fallback Flow (MePass -> QR)', () => {
  let mePassAuth: MePassAuthService;
  let qrAuth: QRCodeAuthService;

  beforeEach(() => {
    jest.clearAllMocks();
    qrAuth = new QRCodeAuthService(mockHttpClient);
    mePassAuth = new MePassAuthService(mockHttpClient, mockSessionStore);
    // inject qrAuth
    (mePassAuth as any).qrCodeAuthService = qrAuth;
  });

  it('should fall back to QR code authentication and return session with idempotency key', async () => {
    const credentials = { username: 'user', password: 'fail' };

    // MePass fails first two times, then succeeds via QR fallback
    mockHttpClient.post
      .mockRejectedValueOnce(new Error('MePass unavailable'))
      .mockRejectedValueOnce(new Error('MePass unavailable'))
      .mockResolvedValueOnce({ data: { qrCode: 'qr-data', sessionId: 'sess-002' } })
      .mockResolvedValueOnce({ data: { status: 'confirmed', token: 'final-session-token' } });

    mockIdempotencyService.generateKey.mockReturnValue('idem-key-flow');

    const result = await mePassAuth.loginWithFallback(credentials, { maxMePassRetries: 2 });

    expect(result).toMatchObject({
      success: true,
      token: 'final-session-token',
      idempotencyKey: 'idem-key-flow',
    });
    expect(mockHttpClient.post).toHaveBeenCalledTimes(4); // 2 mepass fails + 1 qr request + 1 poll
  });

  it('should return error if both MePass and QR fail', async () => {
    mockHttpClient.post.mockRejectedValue(new Error('Service unavailable'));

    const result = await mePassAuth.loginWithFallback({}, { maxMePassRetries: 1, maxQrRetries: 1 });

    expect(result.success).toBe(false);
    expect(result.error).toContain('Authentication failed');
  });
});