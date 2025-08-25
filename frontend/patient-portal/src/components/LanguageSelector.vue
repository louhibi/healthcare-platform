<template>
  <div class="language-selector relative">
    <!-- Trigger Button -->
    <button
      @click="toggleDropdown"
      class="w-full flex items-center justify-between px-2 py-1 text-sm font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-50 rounded-md transition-colors"
      :class="{ 'bg-gray-50': isOpen }"
    >
      <div class="flex items-center space-x-2">
        <span class="text-base">{{ currentLocaleInfo?.flag || '🌐' }}</span>
        <span>{{ currentLocaleInfo?.nativeName || currentLocale }}</span>
      </div>
      <ChevronDownIcon 
        class="h-3 w-3 transition-transform"
        :class="{ 'rotate-180': isOpen }"
      />
    </button>

    <!-- Dropdown Menu -->
    <Transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute left-0 top-full mt-1 w-60 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 z-[60]"
        :class="{ 'right-0 left-auto': !isRTL }"
      >
        <div class="py-1" role="menu">
          <!-- Loading State -->
          <div v-if="isInitializing" class="px-3 py-2 text-sm text-gray-500">
            <div class="flex items-center space-x-2">
              <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
              <span>{{ t('common.loading') }}</span>
            </div>
          </div>

          <!-- Locale Options -->
          <template v-else>
            <button
              v-for="locale in availableLocales"
              :key="locale.code"
              @click="selectLocale(locale.code)"
              class="w-full text-left px-3 py-2 text-sm hover:bg-gray-100 transition-colors flex items-center justify-between"
              :class="{
                'bg-blue-50 text-blue-700': locale.code === currentLocale,
                'text-gray-700': locale.code !== currentLocale
              }"
              role="menuitem"
            >
              <div class="flex items-center space-x-2">
                <span class="text-base">{{ locale.flag }}</span>
                <div>
                  <div class="font-medium text-sm">{{ locale.nativeName }}</div>
                  <div class="text-xs text-gray-500">{{ locale.name }}</div>
                </div>
              </div>
              <CheckIcon 
                v-if="locale.code === currentLocale"
                class="h-4 w-4 text-blue-600"
              />
            </button>
          </template>

          <!-- Admin Translation Management -->
          <template v-if="canManageTranslations">
            <hr class="my-1 border-gray-200" />
            <router-link
              to="/admin/translations"
              class="block px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
              @click="closeDropdown"
            >
              <div class="flex items-center space-x-2">
                <CogIcon class="h-4 w-4" />
                <span>{{ t('admin.translations') }}</span>
              </div>
            </router-link>
          </template>
        </div>
      </div>
    </Transition>

    <!-- Overlay -->
    <div
      v-if="isOpen"
      class="fixed inset-0 z-[55]"
      @click="closeDropdown"
    ></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n as useVueI18n } from 'vue-i18n'
import { ChevronDownIcon, CheckIcon, CogIcon } from '@heroicons/vue/24/outline'
import { useI18n } from '../composables/useI18n'
import { useAuthStore } from '../stores/auth'
import { useToast } from 'vue-toastification'

// Composables
const { t } = useVueI18n()
const { 
  currentLocale, 
  availableLocales, 
  isInitializing, 
  changeLocale,
  isRTL 
} = useI18n()
const authStore = useAuthStore()
const toast = useToast()

// Local state
const isOpen = ref(false)

// Computed properties
const currentLocaleInfo = computed(() => {
  return availableLocales.find(locale => locale.code === currentLocale.value)
})

const canManageTranslations = computed(() => {
  return authStore.user?.role === 'admin'
})

// Methods
function toggleDropdown() {
  isOpen.value = !isOpen.value
}

function closeDropdown() {
  isOpen.value = false
}

async function selectLocale(localeCode) {
  if (localeCode === currentLocale.value) {
    closeDropdown()
    return
  }

  try {
    const success = await changeLocale(localeCode, true)
    if (success) {
      toast.success(t('admin.changeLocale') + ': ' + availableLocales.find(l => l.code === localeCode)?.nativeName)
      closeDropdown()
    } else {
      toast.error(t('errors.general'))
    }
  } catch (error) {
    console.error('Error changing locale:', error)
    toast.error(t('errors.general'))
  }
}

// Handle escape key
function handleEscape(event) {
  if (event.key === 'Escape' && isOpen.value) {
    closeDropdown()
  }
}

// Lifecycle
onMounted(() => {
  document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
})
</script>

<style scoped>
.language-selector {
  /* Ensure dropdown is positioned correctly in RTL layouts */
}

/* RTL adjustments */
html[dir="rtl"] .language-selector .space-x-2 > * + * {
  margin-left: 0;
  margin-right: 0.5rem;
}

html[dir="rtl"] .language-selector .space-x-3 > * + * {
  margin-left: 0;
  margin-right: 0.75rem;
}
</style>