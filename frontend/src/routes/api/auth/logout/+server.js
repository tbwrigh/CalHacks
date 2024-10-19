// src/routes/api/auth/logout.js
export async function GET({ cookies }) {
    // Clear the JWT cookie by setting it with a past expiration date
    cookies.set('at', '', {
      httpOnly: true,
      secure: true, // Ensure this is true in production (HTTPS)
      path: '/',
      expires: new Date(0), // Expire the cookie immediately
      sameSite: 'lax',
    });
  
    // Redirect to home or login page after logout
    return new Response(null, {
      status: 302,
      headers: {
        location: '/',
      },
    });
  }
  