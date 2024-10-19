<script>
    export let data;
    const clientId = import.meta.env.VITE_GITHUB_CLIENT_ID;
    const redirectUri = 'http://localhost:5173/api/auth/callback'; // Replace with your real callback URL
    const githubLoginUrl = `https://github.com/login/oauth/authorize?client_id=${clientId}&redirect_uri=${encodeURIComponent(redirectUri)}&scope=read:user%20user:email`;
</script>
  
  <!-- Main container for the landing page -->
  <div class="min-h-screen bg-gray-100 flex flex-col">
  
    <!-- Top Navbar -->
    <nav class="bg-white border-gray-200 px-4 py-3 sm:px-6 lg:px-8 shadow-sm">
      <div class="flex justify-between items-center mx-auto max-w-7xl">
        <!-- Brand Name -->
        <a href="/" class="text-2xl font-bold text-blue-600">AutoLock</a>
        
        <!-- Login Button -->
        <div>
        {#if data.logged_in}
            <!-- Redirect to dashboard if logged in -->
            <a href="/user" class="text-white bg-green-600 hover:bg-green-700 focus:ring-4 focus:ring-green-300 font-medium rounded-lg text-sm px-4 py-2 text-center">
              Dashboard
            </a>
        {:else}  
          <button 
            class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 text-center mr-3 md:mr-0"
            on:click={() => window.location.href = githubLoginUrl}>
            Login
          </button>
        {/if}
        </div>
      </div>
    </nav>
  
    <!-- Main Content Section -->
    <section class="flex-grow flex items-center justify-center bg-gradient-to-r from-blue-600 to-indigo-600 py-12">
      <div class="text-center max-w-4xl mx-auto px-6 sm:px-8">
        
        <!-- Heading -->
        <h1 class="text-4xl font-extrabold text-white sm:text-5xl lg:text-6xl mb-6">
          Automatically Secure Your Codebase with AI
        </h1>
  
        <!-- Subtext / Description -->
        <p class="text-lg sm:text-xl text-gray-200 mb-8 leading-relaxed">
          AutoLock ensures the safety of your codebase by automatically identifying security vulnerabilities using advanced AI. It integrates with GitHub and helps you maintain your project's security effortlessly.
        </p>
  
        <!-- Details Section -->
        <div class="bg-white rounded-lg shadow-md p-6 sm:p-8 text-left">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">How It Works</h2>
  
          <ul class="list-disc list-inside text-gray-700 space-y-4">
            <li>
              <strong>Login with GitHub:</strong> Easily authenticate and integrate with your repositories.
            </li>
            <li>
              <strong>Use our CLI tool:</strong> Get started in seconds by using the AutoLock CLI to secure your codebase with just a few commands.
            </li>
            <li>
              <strong>CodeQL:</strong> AutoLock leverages GitHub’s powerful CodeQL analysis to identify vulnerabilities in your code.
            </li>
            <li>
              <strong>AI-driven Pull Requests:</strong> Our AI automatically generates pull requests to fix unsafe code for you, making your security process hands-free.
            </li>
          </ul>
        </div>
      </div>
    </section>
  
    <!-- Footer (optional) -->
    <footer class="bg-gray-800 text-gray-400 text-center py-4">
      © 2024 AutoLock. All rights reserved.
    </footer>
  
  </div>
  