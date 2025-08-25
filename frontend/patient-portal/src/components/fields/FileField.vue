<template>
  <div class="space-y-2">
    <!-- File Input -->
    <div class="relative">
      <input
        ref="fileInput"
        v-bind="$attrs"
        type="file"
        :class="fileInputClasses"
        :accept="acceptedTypes"
        :multiple="multiple"
        @input="handleFileInput"
        @change="handleFileChange"
        @dragover.prevent="dragover = true"
        @dragleave.prevent="dragover = false"
        @drop.prevent="handleDrop"
      />
      
      <!-- Custom file input styling -->
      <div 
        v-if="!hasFiles" 
        class="absolute inset-0 flex items-center justify-center pointer-events-none"
        :class="dragOverClasses"
      >
        <div class="text-center">
          <svg class="mx-auto h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 48 48">
            <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          <p class="text-sm text-gray-600 mt-1">
            Drop files here or <span class="font-medium text-indigo-600">browse</span>
          </p>
          <p v-if="acceptedTypes" class="text-xs text-gray-400 mt-1">
            Accepted: {{ acceptedTypes }}
          </p>
        </div>
      </div>
    </div>

    <!-- Selected Files Display -->
    <div v-if="hasFiles" class="space-y-2">
      <div class="text-sm font-medium text-gray-700">Selected Files:</div>
      <div 
        v-for="(file, index) in selectedFiles" 
        :key="`${file.name}-${index}`"
        class="flex items-center justify-between p-2 bg-gray-50 rounded border"
      >
        <div class="flex items-center space-x-2">
          <svg class="h-4 w-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4 4a2 2 0 00-2 2v8a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2H4zm12 2H4v8h12V6z" clip-rule="evenodd" />
          </svg>
          <span class="text-sm text-gray-900 truncate">{{ file.name }}</span>
          <span class="text-xs text-gray-500">({{ formatFileSize(file.size) }})</span>
        </div>
        <button
          v-if="!disabled && !loading"
          @click="removeFile(index)"
          class="text-red-500 hover:text-red-700"
          type="button"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- File Size/Count Limits Info -->
    <div v-if="showLimits" class="text-xs text-gray-500">
      <div v-if="maxFiles && maxFiles > 1">Max files: {{ maxFiles }}</div>
      <div v-if="maxSize">Max size per file: {{ formatFileSize(maxSize) }}</div>
    </div>

    <!-- Upload Progress (if uploading) -->
    <div v-if="uploading" class="w-full bg-gray-200 rounded-full h-2">
      <div 
        class="bg-indigo-600 h-2 rounded-full transition-all duration-300"
        :style="{ width: `${uploadProgress}%` }"
      ></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

// Props
const props = defineProps({
  modelValue: {
    type: [File, Array, FileList],
    default: null
  },
  accept: {
    type: String,
    default: ''
  },
  multiple: {
    type: Boolean,
    default: false
  },
  maxFiles: {
    type: Number,
    default: null
  },
  maxSize: {
    type: Number,
    default: 10 * 1024 * 1024 // 10MB default
  },
  hasError: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  uploading: {
    type: Boolean,
    default: false
  },
  uploadProgress: {
    type: Number,
    default: 0
  }
})

// Emits
const emit = defineEmits({
  'update:modelValue': (value) => true,
  'input': (value) => true,
  'change': (value) => true,
  'file-error': (error) => typeof error === 'string'
})

// Refs
const fileInput = ref(null)

// State
const dragover = ref(false)
const selectedFiles = ref([])

// Computed
const hasFiles = computed(() => selectedFiles.value.length > 0)

const acceptedTypes = computed(() => {
  return props.accept || '*/*'
})

const fileInputClasses = computed(() => {
  const base = 'block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-semibold file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100 cursor-pointer'
  const errorClass = props.hasError ? 'border-red-300 focus:border-red-500' : 'border-gray-300 focus:border-indigo-500'
  const disabledClass = (props.disabled || props.loading) ? 'opacity-50 cursor-not-allowed' : ''
  return `${base} ${errorClass} ${disabledClass}`
})

const dragOverClasses = computed(() => {
  return dragover.value ? 'bg-indigo-50 border-indigo-300' : 'bg-gray-50 border-gray-300'
})

const showLimits = computed(() => {
  return (props.maxFiles && props.maxFiles > 1) || props.maxSize
})

// Methods
const handleFileInput = (event) => {
  const files = Array.from(event.target.files || [])
  processFiles(files)
}

const handleFileChange = (event) => {
  const files = Array.from(event.target.files || [])
  processFiles(files)
  emit('change', props.multiple ? files : files[0])
}

const handleDrop = (event) => {
  dragover.value = false
  const files = Array.from(event.dataTransfer.files || [])
  processFiles(files)
}

const processFiles = (files) => {
  const validFiles = []
  
  for (const file of files) {
    // Check file size
    if (props.maxSize && file.size > props.maxSize) {
      emit('file-error', `File "${file.name}" is too large. Maximum size is ${formatFileSize(props.maxSize)}.`)
      continue
    }
    
    // Check file count
    if (props.maxFiles && validFiles.length >= props.maxFiles) {
      emit('file-error', `Maximum ${props.maxFiles} files allowed.`)
      break
    }
    
    validFiles.push(file)
  }
  
  if (props.multiple) {
    selectedFiles.value = validFiles
    emit('update:modelValue', validFiles)
    emit('input', validFiles)
  } else {
    selectedFiles.value = validFiles.slice(0, 1)
    emit('update:modelValue', validFiles[0] || null)
    emit('input', validFiles[0] || null)
  }
}

const removeFile = (index) => {
  selectedFiles.value.splice(index, 1)
  
  if (props.multiple) {
    emit('update:modelValue', selectedFiles.value)
  } else {
    emit('update:modelValue', selectedFiles.value[0] || null)
  }
  
  // Clear the file input
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Watch for external changes
watch(() => props.modelValue, (newValue) => {
  if (!newValue) {
    selectedFiles.value = []
  } else if (Array.isArray(newValue)) {
    selectedFiles.value = newValue
  } else {
    selectedFiles.value = [newValue]
  }
}, { immediate: true })
</script>