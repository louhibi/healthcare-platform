<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <!-- Header with Navigation -->
    <header v-if="showHeader" class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-4">
          <!-- Logo/Title -->
          <div class="flex items-center space-x-4">
            <div class="flex-shrink-0">
              <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 7.172V5L8 4z" />
              </svg>
            </div>
            <h1 class="text-2xl font-bold text-gray-900">{{ t('navigation.healthcarePortal') || 'Healthcare Portal' }}</h1>
          </div>

          <!-- Navigation -->
          <nav class="hidden md:flex space-x-8">
            <router-link 
              to="/" 
              class="text-gray-500 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              :class="{ 'text-blue-600 bg-blue-50': $route.path === '/' }"
            >
              {{ t('navigation.home') }}
            </router-link>
            <router-link 
              to="/patients" 
              class="text-gray-500 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              :class="{ 'text-blue-600 bg-blue-50': $route.path === '/patients' }"
            >
              {{ t('navigation.patients') }}
            </router-link>
            <router-link 
              to="/appointment-management" 
              class="text-gray-500 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              :class="{ 'text-blue-600 bg-blue-50': $route.path === '/appointment-management' }"
            >
              {{ t('navigation.appointments') }}
            </router-link>
            <router-link 
              to="/availability" 
              v-if="authStore.userRole === 'admin' || authStore.userRole === 'doctor'"
              class="text-gray-500 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              :class="{ 'text-blue-600 bg-blue-50': $route.path === '/availability' }"
            >
              {{ t('navigation.availability') }}
            </router-link>
          </nav>

          <!-- User Menu -->
          <div class="flex items-center space-x-4">
            <!-- User Info -->
            <div class="hidden md:block text-sm text-gray-700">
              <div class="text-right">
                <div class="font-medium">{{ authStore.userName }}</div>
                <div class="text-xs text-gray-500 mt-1">{{ authStore.healthcareEntityName }}</div>
              </div>
              <span class="ml-2 px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded-full capitalize">
                {{ authStore.userRole }}
              </span>
            </div>

            <!-- User Dropdown -->
            <div class="relative">
              <button
                @click="showUserMenu = !showUserMenu"
                class="flex items-center space-x-2 text-gray-500 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md px-2 py-1"
              >
                <div class="h-8 w-8 bg-gray-300 rounded-full flex items-center justify-center">
                  <svg class="h-5 w-5 text-gray-600" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
                  </svg>
                </div>
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>

              <!-- Dropdown Menu -->
              <div 
                v-if="showUserMenu"
                v-click-outside="() => showUserMenu = false"
                class="absolute right-0 mt-2 w-56 bg-white rounded-md shadow-lg py-1 z-50 border border-gray-200"
              >
                <div class="px-4 py-2 text-sm text-gray-700 border-b border-gray-100">
                  <p class="font-medium">{{ authStore.userName }}</p>
                  <p class="text-gray-500 capitalize">{{ authStore.userRole }}</p>
                </div>
                
                <!-- Language Selector in User Menu -->
                <div class="px-4 py-2 border-b border-gray-100">
                  <LanguageSelector />
                </div>
                
                <!-- Administration Panel -->
                <router-link 
                  to="/admin" 
                  v-if="authStore.userRole === 'admin'"
                  @click="showUserMenu = false"
                  class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
                  :class="{ 'bg-blue-50 text-blue-700': $route.path === '/admin' }"
                >
                  <svg class="h-4 w-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  {{ t('navigation.admin') }}
                </router-link>
                
                <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                  {{ t('auth.profileSettings') }}
                </a>
                <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                  {{ t('common.preferences') }}
                </a>
                <button
                  @click="handleLogout"
                  class="w-full text-left px-4 py-2 text-sm text-red-700 hover:bg-red-50 border-t border-gray-100"
                >
                  {{ t('auth.signOut') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main>
      <router-view />
    </main>

    <!-- Password Change Modal -->
    <PasswordChangeModal 
      :show="authStore.isAuthenticated && authStore.isTempPassword"
      @success="handlePasswordChangeSuccess"
    />

    <!-- Debug Toolbar -->
    <DebugToolbar v-if="debugStore.isEnabled" />
  </div>
</template>

<script>
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import DebugToolbar from './components/debug/DebugToolbar.vue'
import LanguageSelector from './components/LanguageSelector.vue'
import PasswordChangeModal from './components/PasswordChangeModal.vue'
import { useAuthStore } from './stores/auth'
import { useDebugStore } from './stores/debug'

export default {
  name: 'App',
  components: {
    PasswordChangeModal,
    LanguageSelector,
    DebugToolbar
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const { t } = useI18n()
    const authStore = useAuthStore()
    const debugStore = useDebugStore()
    const toast = useToast()
    
    const showUserMenu = ref(false)

    // Show header only when authenticated and not on auth pages
    const showHeader = computed(() => {
      return authStore.isAuthenticated && !['Login', 'Register'].includes(route.name)
    })

    const handleLogout = () => {
      authStore.logout()
      showUserMenu.value = false
      toast.success('Logged out successfully')
      router.push('/login')
    }

    const handlePasswordChangeSuccess = () => {
      toast.success('Password changed successfully! You can now use the application.')
    }

    // Initialize stores on app start
    onMounted(async () => {
      await authStore.initialize()
      await debugStore.initialize()
    })

    // Click outside directive for dropdown
    const vClickOutside = {
      beforeMount(el, binding) {
        el.clickOutsideEvent = function(event) {
          if (!(el === event.target || el.contains(event.target))) {
            binding.value()
          }
        }
        document.addEventListener('click', el.clickOutsideEvent)
      },
      unmounted(el) {
        document.removeEventListener('click', el.clickOutsideEvent)
      }
    }

    return {
      t,
      authStore,
      debugStore,
      showHeader,
      showUserMenu,
      handleLogout,
      handlePasswordChangeSuccess,
      vClickOutside
    }
  }
}
</script>

<style>
/* Global styles are handled by Tailwind CSS */
</style>