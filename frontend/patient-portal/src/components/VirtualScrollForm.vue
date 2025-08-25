<template>
  <div 
    ref="container"
    class="virtual-scroll-container"
    :style="{ height: containerHeight }"
    @scroll="handleScroll"
  >
    <!-- Spacer before visible items -->
    <div 
      v-if="startSpacer > 0"
      :style="{ height: `${startSpacer}px` }"
      class="spacer-start"
    ></div>

    <!-- Visible items -->
    <div
      v-for="(item, index) in visibleItems"
      :key="getItemKey(item, visibleStartIndex + index)"
      ref="itemRefs"
      :data-index="visibleStartIndex + index"
      class="virtual-scroll-item"
    >
      <slot 
        :item="item" 
        :index="visibleStartIndex + index"
        :is-virtual="true"
      >
        <div class="p-4 border-b border-gray-200">
          {{ item }}
        </div>
      </slot>
    </div>

    <!-- Spacer after visible items -->
    <div 
      v-if="endSpacer > 0"
      :style="{ height: `${endSpacer}px` }"
      class="spacer-end"
    ></div>

    <!-- Loading indicator -->
    <div 
      v-if="loading && visibleItems.length === 0"
      class="flex justify-center items-center h-32"
    >
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
    </div>

    <!-- Empty state -->
    <div 
      v-if="!loading && items.length === 0"
      class="flex justify-center items-center h-32 text-gray-500"
    >
      <slot name="empty">
        <div class="text-center">
          <p>No items to display</p>
        </div>
      </slot>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { trackRender, trackScroll, trackScrollSessionStart, trackScrollSessionEnd } from '@/utils/virtualScrollMetrics'

// Props
const props = defineProps({
  items: {
    type: Array,
    required: true,
    default: () => []
  },
  itemHeight: {
    type: [Number, Function],
    default: 80,
    validator: (value) => {
      return typeof value === 'number' || typeof value === 'function'
    }
  },
  containerHeight: {
    type: String,
    default: '400px'
  },
  overscan: {
    type: Number,
    default: 5,
    validator: (value) => value >= 0
  },
  getItemKey: {
    type: Function,
    default: (item, index) => {
      return item?.id ?? item?.key ?? index
    }
  },
  loading: {
    type: Boolean,
    default: false
  },
  threshold: {
    type: Number,
    default: 50,
    validator: (value) => value >= 10
  }
})

// Emits
const emit = defineEmits({
  'scroll': (scrollTop, scrollDirection) => true,
  'scroll-end': () => true,
  'visible-range': (startIndex, endIndex) => true
})

// Refs
const container = ref(null)
const itemRefs = ref([])

// State
const scrollTop = ref(0)
const clientHeight = ref(0)
const itemHeights = ref(new Map())
const lastScrollTop = ref(0)
const scrollDirection = ref('down')
const isScrolling = ref(false)
const scrollTimer = ref(null)

// Computed
const totalItems = computed(() => props.items.length)

const averageItemHeight = computed(() => {
  if (itemHeights.value.size === 0) {
    return typeof props.itemHeight === 'number' ? props.itemHeight : 80
  }
  
  const heights = Array.from(itemHeights.value.values())
  return heights.reduce((sum, height) => sum + height, 0) / heights.length
})

const totalHeight = computed(() => {
  if (itemHeights.value.size === props.items.length) {
    // All heights measured
    return Array.from(itemHeights.value.values()).reduce((sum, height) => sum + height, 0)
  }
  
  // Estimate total height
  const measuredCount = itemHeights.value.size
  const unmeasuredCount = totalItems.value - measuredCount
  const measuredHeight = Array.from(itemHeights.value.values()).reduce((sum, height) => sum + height, 0)
  
  return measuredHeight + (unmeasuredCount * averageItemHeight.value)
})

const visibleStartIndex = computed(() => {
  let accumulatedHeight = 0
  let startIndex = 0
  
  for (let i = 0; i < totalItems.value; i++) {
    const itemHeight = getItemHeightByIndex(i)
    
    if (accumulatedHeight + itemHeight > scrollTop.value) {
      startIndex = i
      break
    }
    
    accumulatedHeight += itemHeight
    startIndex = i + 1
  }
  
  return Math.max(0, startIndex - props.overscan)
})

const visibleEndIndex = computed(() => {
  const viewportBottom = scrollTop.value + clientHeight.value
  let accumulatedHeight = 0
  let endIndex = totalItems.value
  
  for (let i = 0; i < totalItems.value; i++) {
    const itemHeight = getItemHeightByIndex(i)
    accumulatedHeight += itemHeight
    
    if (accumulatedHeight >= viewportBottom) {
      endIndex = i
      break
    }
  }
  
  return Math.min(totalItems.value - 1, endIndex + props.overscan)
})

const visibleItems = computed(() => {
  const start = visibleStartIndex.value
  const end = visibleEndIndex.value + 1
  return props.items.slice(start, end)
})

const startSpacer = computed(() => {
  let height = 0
  for (let i = 0; i < visibleStartIndex.value; i++) {
    height += getItemHeightByIndex(i)
  }
  return height
})

const endSpacer = computed(() => {
  const start = visibleEndIndex.value + 1
  let height = 0
  
  for (let i = start; i < totalItems.value; i++) {
    height += getItemHeightByIndex(i)
  }
  
  return height
})

// Methods
const getItemHeightByIndex = (index) => {
  if (itemHeights.value.has(index)) {
    return itemHeights.value.get(index)
  }
  
  if (typeof props.itemHeight === 'function') {
    return props.itemHeight(props.items[index], index)
  }
  
  return props.itemHeight
}

const measureVisibleItems = async () => {
  const startTime = performance.now()
  
  await nextTick()
  
  if (!itemRefs.value || itemRefs.value.length === 0) return
  
  itemRefs.value.forEach((ref) => {
    if (ref && ref.dataset?.index) {
      const index = parseInt(ref.dataset.index)
      const height = ref.getBoundingClientRect().height
      
      if (height > 0) {
        itemHeights.value.set(index, height)
      }
    }
  })
  
  // Track render performance
  const renderTime = performance.now() - startTime
  trackRender(visibleItems.value.length, renderTime)
}

const handleScroll = (event) => {
  const newScrollTop = event.target.scrollTop
  
  // Track scroll performance
  const scrollDistance = Math.abs(newScrollTop - lastScrollTop.value)
  trackScroll(newScrollTop, scrollDistance)
  
  // Determine scroll direction
  scrollDirection.value = newScrollTop > lastScrollTop.value ? 'down' : 'up'
  lastScrollTop.value = newScrollTop
  scrollTop.value = newScrollTop
  
  // Track scroll session start
  if (!isScrolling.value) {
    trackScrollSessionStart()
  }
  
  isScrolling.value = true
  
  // Clear existing timer
  if (scrollTimer.value) {
    clearTimeout(scrollTimer.value)
  }
  
  // Set timer to detect scroll end
  scrollTimer.value = setTimeout(() => {
    isScrolling.value = false
    trackScrollSessionEnd()
    emit('scroll-end')
  }, 150)
  
  // Emit scroll event
  emit('scroll', newScrollTop, scrollDirection.value)
  
  // Measure items after scroll
  requestAnimationFrame(measureVisibleItems)
  
  // Check if near bottom for infinite loading
  const { scrollTop: st, scrollHeight, clientHeight: ch } = event.target
  if (scrollHeight - (st + ch) < props.threshold) {
    emit('scroll-end')
  }
}

const scrollToIndex = (index, alignment = 'auto') => {
  if (!container.value || index < 0 || index >= totalItems.value) return
  
  let targetScrollTop = 0
  for (let i = 0; i < index; i++) {
    targetScrollTop += getItemHeightByIndex(i)
  }
  
  const itemHeight = getItemHeightByIndex(index)
  const currentViewportTop = scrollTop.value
  const currentViewportBottom = scrollTop.value + clientHeight.value
  
  let finalScrollTop = targetScrollTop
  
  if (alignment === 'center') {
    finalScrollTop = targetScrollTop - (clientHeight.value - itemHeight) / 2
  } else if (alignment === 'end') {
    finalScrollTop = targetScrollTop - clientHeight.value + itemHeight
  } else if (alignment === 'auto') {
    // Only scroll if item is not visible
    if (targetScrollTop < currentViewportTop) {
      finalScrollTop = targetScrollTop
    } else if (targetScrollTop + itemHeight > currentViewportBottom) {
      finalScrollTop = targetScrollTop - clientHeight.value + itemHeight
    } else {
      return // Item is already visible
    }
  }
  
  container.value.scrollTop = Math.max(0, finalScrollTop)
}

const scrollToTop = () => {
  if (container.value) {
    container.value.scrollTop = 0
  }
}

const scrollToBottom = () => {
  if (container.value) {
    container.value.scrollTop = container.value.scrollHeight
  }
}

// Update client height when container mounts
const updateClientHeight = () => {
  if (container.value) {
    clientHeight.value = container.value.clientHeight
  }
}

// Watch for visible range changes
watch([visibleStartIndex, visibleEndIndex], ([start, end]) => {
  emit('visible-range', start, end)
})

// Watch for items changes
watch(() => props.items, () => {
  // Clear height cache when items change
  itemHeights.value.clear()
  nextTick(measureVisibleItems)
}, { deep: true })

// Lifecycle
onMounted(() => {
  updateClientHeight()
  nextTick(measureVisibleItems)
  
  // Add resize observer
  const resizeObserver = new ResizeObserver(updateClientHeight)
  if (container.value) {
    resizeObserver.observe(container.value)
  }
  
  onUnmounted(() => {
    resizeObserver.disconnect()
    if (scrollTimer.value) {
      clearTimeout(scrollTimer.value)
    }
  })
})

// Expose methods for parent components
defineExpose({
  scrollToIndex,
  scrollToTop,
  scrollToBottom,
  measureVisibleItems
})
</script>

<style scoped>
.virtual-scroll-container {
  overflow-y: auto;
  overflow-x: hidden;
}

.virtual-scroll-item {
  contain: layout style paint;
}

.spacer-start,
.spacer-end {
  flex-shrink: 0;
}
</style>