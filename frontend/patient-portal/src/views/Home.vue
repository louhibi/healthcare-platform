<template>
  <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
    <div class="px-4 py-6 sm:px-0">
      <div class="border-4 border-dashed border-gray-200 rounded-lg h-96 flex items-center justify-center">
        <div class="text-center">
          <h2 class="text-3xl font-bold text-gray-900 mb-4">Welcome to Healthcare Portal</h2>
          <p class="text-lg text-gray-600 mb-8">Manage your patients and appointments efficiently</p>
          <div class="flex space-x-4 justify-center">
            <router-link 
              to="/appointments" 
              class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
            >
              View Appointments
            </router-link>
            <router-link 
              to="/patients" 
              class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
            >
              Manage Patients
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Home'
}
</script>