<template>
  <div class="min-h-screen bg-gray-50 py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="bg-white shadow rounded-lg p-6 mb-6">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Dynamic Form Configuration Demo</h1>
            <p class="mt-1 text-sm text-gray-500">
              This page demonstrates how forms adapt based on admin configuration settings.
              Visit the Admin Dashboard to configure form fields and see the changes here.
            </p>
          </div>
          <div class="mt-4 sm:mt-0 flex space-x-3">
            <button
              @click="$router.push('/admin')"
              class="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              Configure Forms
            </button>
            <button
              @click="$router.push('/')"
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              Back to Portal
            </button>
          </div>
        </div>
      </div>

      <!-- Instructions -->
      <div class="bg-blue-50 border border-blue-200 rounded-md p-4 mb-6">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-blue-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-blue-800">How to use this demo:</h3>
            <div class="mt-2 text-sm text-blue-700">
              <ol class="list-decimal list-inside space-y-1">
                <li>Click "Configure Forms" to access the admin form configuration panel</li>
                <li>Toggle fields on/off, change requirements, and reorder fields</li>
                <li>Return to this page to see how the form adapts to your configuration</li>
                <li>The form below will only show enabled fields and enforce requirements</li>
              </ol>
            </div>
          </div>
        </div>
      </div>

      <!-- Dynamic Patient Form -->
      <DynamicPatientForm 
        @submit="handleFormSubmit"
        @cancel="handleFormCancel"
      />
    </div>
  </div>
</template>

<script setup>
import DynamicPatientForm from '@/components/patients/DynamicPatientForm.vue'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'

// Composables
const toast = useToast()
const router = useRouter()

// Methods
const handleFormSubmit = (result) => {
  toast.success('Form submitted successfully! (This is just a demo)')
}

const handleFormCancel = () => {
  toast.info('Form cancelled')
  router.push('/')
}
</script>