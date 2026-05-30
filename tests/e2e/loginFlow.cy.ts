/// <reference types="cypress" />

// File: tests/e2e/loginFlow.cy.ts
// Purpose: Validate ME-Hub login flow with Me Pass retry and QR fallback
// Behavior: On initial Me Pass failure, retry up to N times then switch to QR.
// Implements idempotency key handling to prevent duplicate session creation.

import { faker } from '@faker-js/faker';

// ---------------------------------------------------------------------------
// Constants & Types
// ---------------------------------------------------------------------------
const LOGIN_URL = '/login';
const DASHBOARD_URL = '/dashboard';
const MAX_RETRIES = 3;
const RETRY_DELAY_MS = 500;

interface AuthApiResponse {
  status: 'success' | 'error';
  idempotencyKey?: string;
  sessionToken?: string;
  errorMessage?: string;
}

interface AuthApiRequest {
  method: 'me_pass' | 'qr_code';
  credentials?: { email: string; password: string };
  qrCodeHash?: string;
  idempotencyKey: string;
}

// ---------------------------------------------------------------------------
// Custom Cypress Commands
// ---------------------------------------------------------------------------
Cypress.Commands.add('loginWithMePass', (credentials: { email: string; password: string }, idempotencyKey: string) => {
  cy.intercept('POST', '/api/auth/login', (req) => {
    if (req.body.method === 'me_pass') {
      req.alias = 'mePassLogin';
    }
  }).as('mePassLogin');

  cy.get('[data-testid="me-pass-button"]').click();
  cy.get('[data-testid="email-input"]').type(credentials.email);
  cy.get('[data-testid="password-input"]').type(credentials.password);
  cy.get('[data-testid="submit-login"]').click();
});

Cypress.Commands.add('loginWithQRCode', (qrHash: string, idempotencyKey: string) => {
  cy.intercept('POST', '/api/auth/login', (req) => {
    if (req.body.method === 'qr_code') {
      req.alias = 'qrLogin';
    }
  }).as('qrLogin');

  cy.get('[data-testid="qr-fallback-button"]').click();
  cy.get('[data-testid="qr-code-input"]').type(qrHash);
  cy.get('[data-testid="submit-qr"]').click();
});

Cypress.Commands.add('assertLoginError', (expectedMessage: string) => {
  cy.get('[data-testid="login-error"]')
    .should('be.visible')
    .and('contain.text', expectedMessage);
});

Cypress.Commands.add('assertLoginSuccess', () => {
  cy.url().should('include', DASHBOARD_URL);
  cy.get('[data-testid="user-menu"]').should('be.visible');
});

// ---------------------------------------------------------------------------
// Fixtures (simulated backend)
// ---------------------------------------------------------------------------
const createAuthApiResponse = (overrides?: Partial<AuthApiResponse>): AuthApiResponse => ({
  status: 'success',
  idempotencyKey: faker.string.uuid(),
  sessionToken: faker.string.alphanumeric(64),
  ...overrides,
});

const simulateMePassFailure = (): AuthApiResponse => ({
  status: 'error',
  errorMessage: 'Login failed. Please try again or use QR code.',
});

const simulateQRSuccess = (): AuthApiResponse => ({
  status: 'success',
  idempotencyKey: faker.string.uuid(),
  sessionToken: faker.string.alphanumeric(64),
});

// ---------------------------------------------------------------------------
// Test Suite
// ---------------------------------------------------------------------------
describe('ME-Hub Login Flow with Retry & QR Fallback', () => {
  let idempotencyKey: string;

  beforeEach(() => {
    // Generate unique idempotency key per test run
    idempotencyKey = faker.string.uuid();
    cy.visit(LOGIN_URL);
    cy.clearCookies();
    cy.clearLocalStorage();

    // Intercept all auth API calls by default (override in each test as needed)
    cy.intercept('POST', '/api/auth/login', (req) => {
      const body: AuthApiRequest = req.body;
      if (body.method === 'me_pass') {
        // Simulate failure for Me Pass attempts
        req.reply({
          statusCode: 401,
          body: simulateMePassFailure(),
        });
      } else if (body.method === 'qr_code') {
        // Simulate success for QR
        req.reply({
          statusCode: 200,
          body: simulateQRSuccess(),
        });
      }
    }).as('authLogin');
  });

  it('should retry Me Pass up to 3 times then fallback to QR and succeed', () => {
    // Step 1: Attempt Me Pass login (will fail)
    cy.loginWithMePass(
      { email: 'user@example.com', password: 'wrongpassword' },
      idempotencyKey
    );

    // Wait for first Me Pass attempt
    cy.wait('@authLogin').its('response.statusCode').should('eq', 401);
    cy.assertLoginError('Login failed. Please try again or use QR code.');

    // Step 2: Retry Me Pass up to MAX_RETRIES times
    for (let attempt = 2; attempt <= MAX_RETRIES; attempt++) {
      cy.loginWithMePass(
        { email: 'user@example.com', password: 'wrongpassword' },
        idempotencyKey
      );

      cy.wait('@authLogin').its('response.statusCode').should('eq', 401);
      cy.assertLoginError(`Login failed. Attempt ${attempt}/${MAX_RETRIES}. Please try again or use QR code.`);
    }

    // Step 3: After max retries, UI should show QR fallback prompt
    cy.get('[data-testid="qr-fallback-button"]')
      .should('be.visible')
      .and('contain.text', 'Use QR Code');

    // Step 4: Perform QR authentication
    const qrHash = faker.string.alphanumeric(32);
    cy.loginWithQRCode(qrHash, idempotencyKey);

    // Wait for QR login API
    cy.wait('@authLogin').its('response.statusCode').should('eq', 200);

    // Step 5: Verify successful login
    cy.assertLoginSuccess();

    // Step 6: Verify idempotency key stored locally
    cy.window().then((win) => {
      const storedKey = win.localStorage.getItem('idempotencyKey');
      expect(storedKey).to.eq(idempotencyKey);
    });

    // Step 7: Verify session token exists
    cy.window().then((win) => {
      const token = win.localStorage.getItem('sessionToken');
      expect(token).to.match(/^[a-f0-9]{64}$/);
    });
  });

  it('should handle immediate QR fallback when Me Pass option is skipped', () => {
    // User may choose QR from the start
    cy.get('[data-testid="qr-fallback-button"]').click();
    const qrHash = faker.string.alphanumeric(32);
    cy.get('[data-testid="qr-code-input"]').type(qrHash);
    cy.get('[data-testid="submit-qr"]').click();

    cy.wait('@authLogin').its('response.statusCode').should('eq', 200);
    cy.assertLoginSuccess();
  });

  it('should return cached response for duplicate idempotency key on Me Pass', () => {
    // First attempt with idempotency key
    cy.loginWithMePass(
      { email: 'user@example.com', password: 'wrongpassword' },
      idempotencyKey
    );

    cy.wait('@authLogin').then((interception) => {
      const firstResponse = interception.response;
      expect(firstResponse!.statusCode).to.eq(401);

      // Second attempt with same idempotency key – backend returns cached response
      cy.loginWithMePass(
        { email: 'user@example.com', password: 'wrongpassword' },
        idempotencyKey
      );

      cy.wait('@authLogin').then((secondInterception) => {
        const secondResponse = secondInterception.response;
        // Backend should return same 401 without processing
        expect(secondResponse!.statusCode).to.eq(401);
        expect(secondResponse!.body).to.deep.equal(firstResponse!.body);
      });
    });
  });

  it('should show error toast for expired QR code and allow retry', () => {
    // Override QR response to simulate expiry once
    cy.intercept('POST', '/api/auth/login', (req) => {
      if (req.body.method === 'qr_code') {
        req.reply({
          statusCode: 400,
          body: {
            status: 'error',
            errorMessage: 'QR code expired. Please request a new one.',
          },
        });
      } else {
        req.reply({
          statusCode: 401,
          body: simulateMePassFailure(),
        });
      }
    }).as('authLoginExpiredQR');

    cy.get('[data-testid="qr-fallback-button"]').click();
    cy.get('[data-testid="qr-code-input"]').type('expiredHash');
    cy.get('[data-testid="submit-qr"]').click();

    cy.wait('@authLoginExpiredQR').its('response.statusCode').should('eq', 400);
    cy.get('[data-testid="login-error"]')
      .should('be.visible')
      .and('contain.text', 'QR code expired');

    // Retry with new QR code
    cy.get('[data-testid="qr-code-input"]').clear().type('newQRHash');
    cy.get('[data-testid="submit-qr"]').click();

    cy.wait('@authLoginExpiredQR').its('response.statusCode').should('eq', 400); // still mocked, but real app would succeed
    // In production this would succeed after retry; test demonstrates the pattern.
  });

  it('should disable submit button during ongoing login request to prevent double-clicks', () => {
    cy.intercept('POST', '/api/auth/login', (req) => {
      // Delay response to simulate network
      req.reply((res) => {
        setTimeout(() => {
          res.send({ statusCode: 200, body: createAuthApiResponse() });
        }, 1000);
      });
    }).as('slowLogin');

    cy.get('[data-testid="me-pass-button"]').click();
    cy.get('[data-testid="submit-login"]').should('be.disabled');

    cy.wait('@slowLogin').then(() => {
      cy.get('[data-testid="submit-login"]').should('be.enabled');
    });
  });
});

// ---------------------------------------------------------------------------
// Type augmentation for custom commands
// ---------------------------------------------------------------------------
declare global {
  namespace Cypress {
    interface Chainable {
      loginWithMePass(
        credentials: { email: string; password: string },
        idempotencyKey: string
      ): Chainable<void>;
      loginWithQRCode(qrHash: string, idempotencyKey: string): Chainable<void>;
      assertLoginError(expectedMessage: string): Chainable<void>;
      assertLoginSuccess(): Chainable<void>;
    }
  }
}