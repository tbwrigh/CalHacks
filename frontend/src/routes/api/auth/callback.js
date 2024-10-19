export async function GET({ url }) {
    const code = url.searchParams.get('code');
  
    const response = await fetch('https://github.com/login/oauth/access_token', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      body: JSON.stringify({
        client_id: import.meta.env.VITE_GITHUB_CLIENT_ID,
        client_secret: import.meta.env.VITE_GITHUB_CLIENT_SECRET, // Add this to your .env file
        code,
        redirect_uri: 'http://localhost:3000/api/auth/callback', // Replace with your real callback URL
      }),
    });
  
    const data = await response.json();
    const accessToken = data.access_token;
  
    // Now you can use the access token to fetch user data
    const userResponse = await fetch('https://api.github.com/user', {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
  
    const user = await userResponse.json();
  
    // Store the user in session, or redirect to a protected route like /dashboard
    return new Response(null, {
      status: 302,
      headers: {
        location: `/dashboard?user=${encodeURIComponent(JSON.stringify(user))}`,
      },
    });
  }
  