<script>
    export let data;

    console.log(data);
    function logout() {
        window.location.href = '/api/auth/logout';
    }
</script>

<div class="min-h-screen bg-gray-100 flex flex-col">

    <!-- Top Navbar -->
    <nav class="bg-white border-gray-200 px-4 py-3 sm:px-6 lg:px-8 shadow-sm">
        <div class="flex justify-between items-center mx-auto max-w-7xl">
            <!-- Brand Name -->
            <a href="/user" class="text-2xl font-bold text-blue-600">AutoLock</a>

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

    <main class="flex-grow flex flex-col items-center justify-center p-8">
        {#if data.error}
        <!-- Display error message if an error occurred -->
        <div class="text-center">
          <h1 class="text-3xl font-bold text-red-600">Error</h1>
          <p class="text-lg text-gray-600 mt-4">{data.error}</p>
        </div>
      {:else if data.scanComplete}
        <!-- Display scan result -->
        <div class="text-center">
          <h1 class="text-3xl font-bold text-green-600 mb-4">Success: Scan Complete</h1>

          <!-- Display repository name -->
          <h2 class="text-2xl font-semibold text-gray-800 mb-6">Repository: {data.name}</h2>

          <!-- Detected technologies -->
          <div class="w-full max-w-lg text-left">
            <h3 class="text-xl font-semibold text-gray-700 mb-4">Technologies Detected:</h3>
            {#if data.resultData.languages.length > 0}
              <ul class="list-disc list-inside">
                {#each data.resultData.languages as language}
                  <li class="text-gray-700">
                    {language.Name} <span class="text-sm text-gray-500">({language.Extension})</span>
                  </li>
                {/each}
              </ul>
            {:else}
              <p class="text-gray-600">No technologies detected.</p>
            {/if}
          </div>

          <!-- Install button (does nothing for now) -->
          <button
            class="mt-6 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:ring-4 focus:ring-blue-300">
            Install
          </button>
        </div>
      {:else}
        <!-- Display loading bar while waiting for scan completion -->
        <div class="w-full max-w-md">
          <p class="text-xl font-semibold text-gray-700 mb-4">Scanning in progress...</p>
          <div class="w-full bg-gray-300 rounded-full h-4 mb-4">
            <div class="bg-blue-600 h-4 rounded-full animate-pulse" style="width: 75%"></div>
          </div>
          <p class="text-sm text-gray-600">Please wait while we scan your repository.</p>
        </div>
      {/if}
    </main>

    <!-- Footer (optional) -->
    <footer class="bg-gray-800 text-gray-400 text-center py-4">
        © 2024 AutoLock. All rights reserved.
    </footer>
</div>
      