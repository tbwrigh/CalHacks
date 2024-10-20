export async function GET({ url, cookies }) {
  const code = url.searchParams.get('code');

  if (!code) {
    return new Response(JSON.stringify({ error: 'No authorization code returned' }), { status: 400 });
  }

  const response = await fetch('https://github.com/login/oauth/access_token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
    body: JSON.stringify({
      client_id: import.meta.env.VITE_GITHUB_CLIENT_ID,
      client_secret: import.meta.env.VITE_GITHUB_CLIENT_SECRET,
      code,
      redirect_uri: 'http://localhost:5173/api/auth/callback',
    }),
  });

  const data = await response.json();

  if (!data.access_token) {
    return new Response(JSON.stringify({ error: 'No access token returned' }), { status: 400 });
  }

  cookies.set('at', data.access_token, {
    httpOnly: false,      // Prevents JavaScript access to the cookie
    secure: false,        // Set this to true for production (HTTPS)
    path: '/',
    maxAge: 60 * 60,     // 1 hour
    sameSite: 'lax',     // Protect against CSRF
  });

  return new Response(null, {
    status: 302,
    headers: {
      location: `/user`,
    },
  });
}
