<script>
    export let data;

    
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
    <main class="flex-grow flex flex-col items-center justify-center p-8">
      {#if data.error}
        <!-- Show error message if there's an error -->
        <div class="text-center">
          <h1 class="text-3xl font-bold text-red-600">Error</h1>
          <p class="text-lg text-gray-600 mt-4">{data.error}</p>
        </div>
      {:else if data.hasAccess}
      <h1 class="text-3xl font-bold text-green-600 mb-6">Welcome, {data.user}!</h1>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each data.repos as repo}
          <div class="bg-white shadow-md rounded-lg p-6">
            <h2 class="text-xl font-bold text-gray-900">{repo.name}</h2>
            <p class="text-sm text-gray-500">{repo.full_name}</p>
            <p class="text-sm text-gray-700 mt-2">{repo.description ? repo.description : 'No description provided.'}</p>
            <p class="text-sm text-gray-600 mt-2">Language: {repo.language}</p>
            <p class="text-sm text-gray-600 mt-2">Owner: {repo.owner.login}</p>
          </div>
        {/each}
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
  