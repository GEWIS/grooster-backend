import { jwtDecode, JwtPayload } from 'jwt-decode';

type Token = { token: string; expires: string };

export function clearTokenInStorage() {
  localStorage.clear();
}

export function parseToken(rawToken: string): Token {
  const expires = String(jwtDecode<JwtPayload>(rawToken).exp);
  return { token: rawToken, expires };
}

export function setTokenInStorage(jwtToken: string) {
  localStorage.setItem('jwt_token', JSON.stringify(parseToken(jwtToken)));
}

export function getTokenFromStorage(): Token {
  const rawToken = localStorage.getItem('access_token');

  let token = {} as Token;
  if (rawToken !== null) token = parseToken(rawToken);

  return {
    ...token,
  };
}

export function isTokenExpired(tokenEpochTimestamp: number): boolean {
  // If the expiry is somewhere around the year 33658 assume migration from old frontend and force re-login.
  if (tokenEpochTimestamp > 1000000000000) return true;

  const tokenExpirationTime = tokenEpochTimestamp * 1000;
  const currentTimestamp = new Date().getTime();
  return currentTimestamp > tokenExpirationTime;
}

/**
 * Returns True if there is a token in the LocalStorage and if it hasn't expired yet.
 */
export function isAuthenticated(): boolean {
  const token = getTokenFromStorage();
  if (!token.token || !token.expires) return false;
  return !isTokenExpired(Number(token.expires));
}

export function loginRedirect() {
  if (window.location.origin.includes('http://localhost')) {
    window.location.href = 'http://localhost:8080/api/v1/auth/redirect';
    return;
  }

  window.location.href = window.location.origin + '/api/v1/auth/redirect';
}

interface CustomJwtPayload extends Omit<JwtPayload, 'sub'> {
  preferred_username?: string;
  name?: string;
  sub?: number | string;
}

export function getGEWISId(): number {
  const rawToken = localStorage.getItem('access_token');
  if (!rawToken) return;

  const decoded = jwtDecode<CustomJwtPayload>(rawToken);

  let userId: number | undefined = undefined;

  if (decoded.preferred_username) {
    // @deprecated 2026-05-31 As this will then become old
    const match = decoded.preferred_username.match(/m(\d+)/);
    if (match) {
      userId = parseInt(match[1], 10);
    }
  } else if (decoded.sub) {
    userId = typeof decoded.sub === 'number' ? decoded.sub : parseInt(decoded.sub, 10);
  }

  return userId;
}
