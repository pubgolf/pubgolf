export const authHelper = {
  isAuthorized (session) {
    return Boolean(session.user && session.user.authtoken);
  },

  async logIn ({ eventKey, phoneNumber, authCode }, _fetch = fetch) {
    return _fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        eventKey,
        phoneNumber,
        authCode,
      }),
    });
  },
  async restoreSession (_fetch = fetch) {
    const response = await _fetch('/api/login', {
      method: 'GET',
    });
    if (!response.ok) {
      throw new Error(`HTTP error, status = ${response.status}`);
    }
    return response.text();
  },
};
