export const authHelper = {
  isAuthorized ({ user }, eventKey) {
    return Boolean(user && user.authToken && eventKey === user.eventKey);
  },

  async preserveSession (user, _fetch = fetch) {
    return _fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(user),
    });
  },

  async restoreSession (_fetch = fetch) {
    const response = await _fetch('/api/login', {
      method: 'GET',
    });
    if (!response.ok) {
      throw new Error(`HTTP error, status = ${response.status}`);
    }
    return response.json();
  },
};
