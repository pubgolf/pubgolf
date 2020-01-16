export const authHelper = {
  isAuthorized (session) {
    return Boolean(session.user && session.user.authtoken);
  },

  async preserveSession ({ eventKey, user }) {
    return fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        eventKey,
        user,
      }),
    });
  },
  async restoreSession () {
    const response = await fetch('/api/login', {
      method: 'GET',
    });
    if (!response.ok) {
      throw new Error(`HTTP error, status = ${response.status}`);
    }
    return response.text();
  },
};
