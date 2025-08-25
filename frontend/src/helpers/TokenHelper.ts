import {jwtDecode, JwtPayload} from 'jwt-decode';

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
    const rawToken = localStorage.getItem('access_token') as string;

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
    // ApiService.auth.authRedirectGet().then((response) => {
    //     window.location.href = response.data;
    // }).catch((err) => {
    //     console.error("Failed redirecting:", err);
    // });
    window.location.href = "http://localhost:8080/api/v1/auth/redirect";
}

export function getGEWISId(): number {
    const rawToken = localStorage.getItem('access_token') as string
    const decoded = jwtDecode<JwtPayload & { preferred_username?: string }>(rawToken);

    let userId: number | undefined = undefined;
    if (decoded.preferred_username) {
        const match = decoded.preferred_username.match(/m(\d+)/);
        if (match) {
            userId = parseInt(match[1], 10);
        }
    }

    return userId;
}