<script>
    let hasAccess = true; 
    
    function logout() {
        window.location.href = '/api/auth/logout';
    }
  </script>
  
  <!-- Main container -->
  <div class="min-h-screen bg-gray-100 flex flex-col">
  
    <!-- Top Navbar -->
    <nav class="bg-white border-gray-200 px-4 py-3 sm:px-6 lg:px-8 shadow-sm">
      <div class="flex justify-between items-center mx-auto max-w-7xl">
        <!-- Brand Name -->
        <a href="/" class="text-2xl font-bold text-blue-600">AutoLock</a>
  
        <!-- Logout Button -->
        <div>
          <button
            class="text-white bg-red-600 hover:bg-red-700 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-sm px-4 py-2 text-center"
            on:click={logout}>
            Log Out
          </button>
        </div>
      </div>
    </nav>
  
    <!-- Content Section -->
    <main class="flex-grow flex items-center justify-center">
      {#if hasAccess}
        <!-- Content when the user has access -->
        <div class="text-center">
          <h1 class="text-4xl font-bold text-green-600">Welcome to AutoLock!</h1>
          <p class="text-lg text-gray-600 mt-4">You now have full access to the system.</p>
        </div>
      {:else}
        <!-- Content when the user is on the waitlist -->
        <div class="text-center">
          <h1 class="text-3xl font-bold text-red-600">We're sorry!</h1>
          <p class="text-lg text-gray-600 mt-4">
            You're still on the waitlist and don't have access at the moment. We appreciate your patience.
          </p>
        </div>
      {/if}
    </main>
  
    <!-- Footer (optional) -->
    <footer class="bg-gray-800 text-gray-400 text-center py-4">
      © 2024 AutoLock. All rights reserved.
    </footer>
  </div>
  