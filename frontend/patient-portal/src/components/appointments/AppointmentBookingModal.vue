<template>
  <TransitionRoot as="template" :show="show">
    <Dialog as="div" class="relative z-50" @close="close">
      <TransitionChild
        as="template"
        enter="ease-out duration-300"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="ease-in duration-200"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
      </TransitionChild>

      <div class="fixed inset-0 z-10 overflow-y-auto">
        <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
          <TransitionChild
            as="template"
            enter="ease-out duration-300"
            enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            enter-to="opacity-100 translate-y-0 sm:scale-100"
            leave="ease-in duration-200"
            leave-from="opacity-100 translate-y-0 sm:scale-100"
            leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          >
            <DialogPanel class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-4xl sm:p-6">
              <div>
                <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-green-100">
                  <CalendarDaysIcon class="h-6 w-6 text-green-600" />
                </div>
                <div class="mt-3 text-center sm:mt-5">
                  <DialogTitle as="h3" class="text-base font-semibold leading-6 text-gray-900">
                    Book New Appointment
                  </DialogTitle>
                  <div class="mt-2">
                    <p class="text-sm text-gray-500">
                      Schedule a new appointment with conflict detection and alternative suggestions
                    </p>
                  </div>
                </div>
              </div>

              <!-- Step indicator -->
              <div class="mt-6">
                <nav aria-label="Progress">
                  <ol class="flex items-center justify-center">
                    <li v-for="(step, stepIdx) in steps" :key="step.name" :class="stepIdx !== steps.length - 1 ? 'pr-8 sm:pr-20' : ''" class="relative">
                      <div v-if="step.status === 'complete'" class="absolute inset-0 flex items-center" aria-hidden="true">
                        <div class="h-0.5 w-full bg-indigo-600" />
                      </div>
                      <div class="relative flex h-8 w-8 items-center justify-center rounded-full" 
                           :class="step.status === 'complete' ? 'bg-indigo-600' : step.status === 'current' ? 'border-2 border-indigo-600 bg-white' : 'border-2 border-gray-300 bg-white'">
                        <CheckIcon v-if="step.status === 'complete'" class="h-5 w-5 text-white" />
                        <span v-else-if="step.status === 'current'" class="h-2.5 w-2.5 rounded-full bg-indigo-600" />
                        <span v-else class="h-2.5 w-2.5 rounded-full bg-transparent" />
                      </div>
                      <span class="absolute -bottom-6 left-1/2 -translate-x-1/2 text-xs font-medium text-gray-500">
                        {{ step.name }}
                      </span>
                    </li>
                  </ol>
                </nav>
              </div>

              <div class="mt-12">
                <!-- Step 1: Basic Information -->
                <div v-if="currentStep === 1">
                  <!-- Dynamic Form Configuration Loading -->
                  <div v-if="appointmentFormConfig.isLoading.value" class="text-center py-8">
                    <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
                    <p class="mt-4 text-sm text-gray-500">Loading form configuration...</p>
                  </div>
                  
                  <!-- Dynamic form sections based on configuration -->
                  <template v-else>
                    <div class="space-y-6">
                      <!-- Core appointment fields -->
                      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
                        <!-- Patient Selection - special handling -->
                        <div v-if="isFieldVisible('patient_id')" class="sm:col-span-2">
                          <SearchableSelect
                            v-model="bookingForm.patient_id"
                            :label="getFieldDisplayName('patient_id') + (isFieldRequired('patient_id') ? ' *' : '')"
                            placeholder="Search patients by name, ID, phone..."
                            :required="isFieldRequired('patient_id')"
                            :items="patients"
                            item-type="patients"
                            :display-key="patient => `${patient.first_name} ${patient.last_name}`"
                            :secondary-key="patient => formatPatientInfo(patient)"
                            :search-keys="['first_name', 'last_name', 'patient_id', 'phone', 'email']"
                            :loading="loadingPatients"
                            loading-text="Loading patients..."
                            item-icon="UserIcon"
                            @search="onPatientSearch"
                            @select="onPatientSelected"
                          />
                          <div v-if="appointmentFormConfig.hasFieldError('patient_id')" class="mt-1">
                            <p v-for="error in appointmentFormConfig.getFieldError('patient_id')" :key="error" class="text-xs text-red-600">
                              {{ error }}
                            </p>
                          </div>
                        </div>

                        <!-- Doctor Selection - special handling -->
                        <div v-if="isFieldVisible('doctor_id')" class="sm:col-span-2">
                          <label for="doctor-select" class="block text-sm font-medium text-gray-700">
                            {{ getFieldDisplayName('doctor_id') }}
                            <span v-if="isFieldRequired('doctor_id')" class="text-red-500 ml-1">*</span>
                          </label>
                          <select
                            id="doctor-select"
                            v-model="bookingForm.doctor_id"
                            :required="isFieldRequired('doctor_id')"
                            @change="onDoctorChange"
                            :class="[
                              'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                              appointmentFormConfig.hasFieldError('doctor_id') ? 'border-red-500' : ''
                            ]"
                          >
                            <option value="">Select Doctor</option>
                            <option v-for="doctor in doctors" :key="doctor.id" :value="doctor.id">
                              Dr. {{ doctor.first_name }} {{ doctor.last_name }} - {{ doctor.specialization }}
                            </option>
                          </select>
                          <div v-if="appointmentFormConfig.hasFieldError('doctor_id')" class="mt-1">
                            <p v-for="error in appointmentFormConfig.getFieldError('doctor_id')" :key="error" class="text-xs text-red-600">
                              {{ error }}
                            </p>
                          </div>
                        </div>
                      </div>
                      
                      <!-- Dynamic form fields by category -->
                      <div v-for="(categoryFields, categoryName) in getFieldsByCategory()" :key="categoryName" class="space-y-4">
                        <div v-if="categoryFields.filter(f => f.is_enabled && !['patient_id', 'doctor_id'].includes(f.name)).length > 0" class="border-t pt-4 first:border-t-0 first:pt-0">
                          <h4 v-if="categoryName !== 'Other'" class="text-md font-medium text-gray-900 mb-4">{{ categoryName }}</h4>
                          
                          <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
                            <template v-for="field in categoryFields" :key="field.name">
                              <div v-if="field.is_enabled && !['patient_id', 'doctor_id'].includes(field.name)" :class="field.field_type === 'textarea' ? 'sm:col-span-2' : ''">
                                <label :for="'booking-' + field.name" class="block text-sm font-medium text-gray-700">
                                  {{ field.display_name }}
                                  <span v-if="field.is_required" class="text-red-500 ml-1">*</span>
                                </label>
                                
                                <!-- Special handling for appointment type with duration loading -->
                                <select
                                  v-if="field.name === 'type'"
                                  :id="'booking-' + field.name"
                                  v-model="bookingForm[field.name]"
                                  :required="field.is_required"
                                  :class="[
                                    'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                    appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                                  ]"
                                >
                                  <option value="">Select {{ field.display_name }}</option>
                                  <option v-for="option in field.options" :key="option.value" :value="option.value">
                                    {{ option.label }}
                                  </option>
                                </select>
                                
                                <!-- Special handling for duration with dynamic options -->
                                <select
                                  v-else-if="field.name === 'duration'"
                                  :id="'booking-' + field.name"
                                  v-model="bookingForm[field.name]"
                                  :required="field.is_required"
                                  :disabled="!bookingForm.type || loadingDurationOptions"
                                  :class="[
                                    'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm disabled:opacity-50',
                                    appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                                  ]"
                                >
                                  <option value="">{{ !bookingForm.type ? 'Select appointment type first' : loadingDurationOptions ? 'Loading options...' : 'Select duration' }}</option>
                                  <option 
                                    v-for="option in availableDurationOptions" 
                                    :key="option.id" 
                                    :value="option.duration_minutes"
                                    :class="{ 'font-semibold': option.is_default }"
                                  >
                                    {{ option.duration_minutes }} minutes{{ option.is_default ? ' (Default)' : '' }}
                                  </option>
                                </select>
                                
                                <!-- Text input fields -->
                                <input
                                  v-else-if="['text', 'email', 'tel', 'url'].includes(field.field_type)"
                                  :id="'booking-' + field.name"
                                  v-model="bookingForm[field.name]"
                                  :type="field.field_type"
                                  :required="field.is_required"
                                  :class="[
                                    'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                    appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                                  ]"
                                  :placeholder="field.description || ''"
                                />
                                
                                <!-- Textarea fields -->
                                <textarea
                                  v-else-if="field.field_type === 'textarea'"
                                  :id="'booking-' + field.name"
                                  v-model="bookingForm[field.name]"
                                  :required="field.is_required"
                                  :rows="field.name === 'reason' ? 3 : 2"
                                  :class="[
                                    'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                    appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                                  ]"
                                  :placeholder="field.description || (field.name === 'reason' ? 'Please describe the reason for this appointment...' : 'Any additional information...')"
                                ></textarea>
                                
                                <!-- Select fields (other than type and duration) -->
                                <select
                                  v-else-if="field.field_type === 'select' && !['type', 'duration'].includes(field.name)"
                                  :id="'booking-' + field.name"
                                  v-model="bookingForm[field.name]"
                                  :required="field.is_required"
                                  :class="[
                                    'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                    appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                                  ]"
                                >
                                  <option value="">Select {{ field.display_name }}</option>
                                  <option v-for="option in field.options" :key="option.value" :value="option.value">
                                    {{ option.label }}
                                  </option>
                                </select>
                                
                                <!-- Number fields -->
                                <input
                                  v-else-if="field.field_type === 'number'"
                                  :id="'booking-' + field.name"
                                  v-model="bookingForm[field.name]"
                                  type="number"
                                  :required="field.is_required"
                                  :class="[
                                    'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                    appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                                  ]"
                                />
                                
                                <!-- Boolean/checkbox fields -->
                                <div v-else-if="field.field_type === 'boolean'" class="mt-1">
                                  <label class="inline-flex items-center">
                                    <input
                                      :id="'booking-' + field.name"
                                      v-model="bookingForm[field.name]"
                                      type="checkbox"
                                      class="form-checkbox h-4 w-4 text-indigo-600"
                                    />
                                    <span class="ml-2 text-sm text-gray-600">{{ field.description || field.display_name }}</span>
                                  </label>
                                </div>
                                
                                <!-- Fallback for other field types -->
                                <input
                                  v-else
                                  :id="'booking-' + field.name"
                                  v-model="bookingForm[field.name]"
                                  type="text"
                                  :required="field.is_required"
                                  :class="[
                                    'mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm',
                                    appointmentFormConfig.hasFieldError(field.name) ? 'border-red-500' : ''
                                  ]"
                                  :placeholder="field.description || ''"
                                />
                                
                                <!-- Field errors -->
                                <div v-if="appointmentFormConfig.hasFieldError(field.name)" class="mt-1">
                                  <p v-for="error in appointmentFormConfig.getFieldError(field.name)" :key="error" class="text-xs text-red-600">
                                    {{ error }}
                                  </p>
                                </div>
                                
                                <!-- Special message for duration field -->
                                <p v-if="field.name === 'duration' && bookingForm.type && availableDurationOptions.length === 0 && !loadingDurationOptions" class="mt-1 text-sm text-red-600">
                                  No duration options configured for this appointment type
                                </p>
                                
                                <!-- Field description -->
                                <p v-else-if="field.description && field.field_type !== 'boolean'" class="mt-1 text-xs text-gray-500">
                                  {{ field.description }}
                                </p>
                              </div>
                            </template>
                          </div>
                        </div>
                      </div>
                    </div>
                  </template>
                </div>

                <!-- Step 2: Date & Time Selection -->
                <div v-else-if="currentStep === 2" class="space-y-6">
                  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    <!-- Available Dates Selection -->
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">
                        Available Dates *
                      </label>
                      
                      <div v-if="loadingAvailableDates" class="text-center py-4 border border-gray-200 rounded-md">
                        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600 mx-auto"></div>
                        <p class="mt-2 text-sm text-gray-500">Loading available dates...</p>
                      </div>
                      
                      <div v-else-if="availableDates.length === 0" class="text-center py-8 border border-gray-200 rounded-md">
                        <CalendarDaysIcon class="mx-auto h-12 w-12 text-gray-400" />
                        <p class="mt-2 text-sm text-gray-500">No available dates</p>
                        <p class="text-xs text-gray-400">{{ selectedDoctorName }} has no availability</p>
                      </div>
                      
                      <div v-else class="space-y-2 max-h-64 overflow-y-auto border border-gray-200 rounded-md p-2">
                        <div v-for="dateGroup in groupedAvailableDates" :key="dateGroup.month" class="mb-4">
                          <h4 class="font-medium text-gray-900 mb-2 px-2">{{ dateGroup.month }}</h4>
                          <div class="grid grid-cols-7 gap-1">
                            <button
                              v-for="date in dateGroup.dates"
                              :key="date.date"
                              @click="selectDate(date.date)"
                              :class="[
                                'p-2 text-sm rounded-md border transition-colors',
                                bookingForm.date === date.date
                                  ? 'bg-indigo-600 text-white border-indigo-600'
                                  : 'bg-white text-gray-900 border-gray-200 hover:bg-gray-50',
                                date.isToday ? 'ring-2 ring-indigo-200' : ''
                              ]"
                            >
                              <div class="font-medium">{{ date.day }}</div>
                              <div class="text-xs opacity-75">{{ date.weekday }}</div>
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>

                    <!-- Time Slot Selection -->
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">
                        Available Time Slots
                      </label>
                      
                      <div v-if="loadingSlots" class="text-center py-4">
                        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600 mx-auto"></div>
                        <p class="mt-2 text-sm text-gray-500">Loading available slots...</p>
                      </div>

                      <div v-else-if="availableSlots.length === 0 && bookingForm.date" class="text-center py-4">
                        <ClockIcon class="mx-auto h-12 w-12 text-gray-400" />
                        <p class="mt-2 text-sm text-gray-500">No available slots for this date</p>
                        <p class="text-xs text-gray-400">Try selecting a different date</p>
                      </div>

                      <div v-else-if="availableSlots.length > 0" class="space-y-4">
                        <div v-for="(slotGroup, timeOfDay) in groupedSlots" :key="timeOfDay">
                          <h4 class="font-medium text-gray-900 mb-2 capitalize">{{ timeOfDay }}</h4>
                          <div class="grid grid-cols-3 gap-2">
                            <button
                              v-for="slot in slotGroup"
                              :key="slot.date_time"
                              @click="selectTimeSlot(slot)"
                              :disabled="!slot.is_available"
                              :class="[
                                'px-3 py-2 text-xs font-medium rounded-md border',
                                slot.is_available
                                  ? bookingForm.time_slot === formatTime24(slot.date_time)
                                    ? 'bg-indigo-600 text-white border-indigo-600'
                                    : 'bg-white text-gray-900 border-gray-300 hover:bg-gray-50'
                                  : 'bg-gray-100 text-gray-400 border-gray-200 cursor-not-allowed'
                              ]"
                            >
                              {{ formatTime(slot.date_time) }}
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  
                  <!-- Room Selection (only show if date and time are selected) -->
                  <div v-if="bookingForm.date && bookingForm.time_slot" class="mt-6">
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      <span>Available Rooms for {{ formatDate(bookingForm.date) }} at {{ bookingForm.time_slot }}</span>
                      <span v-if="requireRoomAssignment" class="ml-2 text-red-600 text-xs font-normal">*Required</span>
                      <span v-else class="ml-2 text-gray-500 text-xs font-normal">(Optional)</span>
                    </label>
                    
                    <div v-if="loadingRooms" class="text-center py-4">
                      <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600 mx-auto"></div>
                      <p class="mt-2 text-sm text-gray-500">Loading available rooms...</p>
                    </div>
                    
                    <div v-else-if="availableRooms.length === 0" class="text-center py-4 text-gray-500">
                      <p v-if="requireRoomAssignment" class="text-sm text-red-600">No rooms available for this time slot. Please select a different time.</p>
                      <p v-else class="text-sm">No specific room requirements - room will be assigned automatically</p>
                    </div>
                    
                    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                      <!-- Option for no specific room preference (only if not required) -->
                      <button
                        v-if="!requireRoomAssignment"
                        @click="bookingForm.room_id = ''"
                        :class="[
                          'p-3 text-left border rounded-lg transition-colors',
                          bookingForm.room_id === ''
                            ? 'border-indigo-500 bg-indigo-50 ring-2 ring-indigo-500'
                            : 'border-gray-300 hover:border-gray-400 hover:bg-gray-50'
                        ]"
                      >
                        <div class="flex items-center">
                          <div class="flex-shrink-0">
                            <div class="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center">
                              <span class="text-sm">🏥</span>
                            </div>
                          </div>
                          <div class="ml-3 flex-1">
                            <p class="text-sm font-medium text-gray-900">Any Available Room</p>
                            <p class="text-xs text-gray-500">System will assign automatically</p>
                          </div>
                          <div v-if="bookingForm.room_id === ''" class="flex-shrink-0">
                            <CheckIcon class="h-5 w-5 text-indigo-600" />
                          </div>
                        </div>
                      </button>
                      
                      <!-- Room options (available and unavailable) -->
                      <button
                        v-for="room in availableRooms"
                        :key="room.id"
                        @click="room.is_available ? bookingForm.room_id = room.id : null"
                        :disabled="!room.is_available"
                        :class="[
                          'p-3 text-left border rounded-lg transition-colors',
                          !room.is_available
                            ? 'border-gray-200 bg-gray-50 cursor-not-allowed opacity-60'
                            : bookingForm.room_id === room.id
                              ? 'border-indigo-500 bg-indigo-50 ring-2 ring-indigo-500'
                              : 'border-gray-300 hover:border-gray-400 hover:bg-gray-50'
                        ]"
                      >
                        <div class="flex items-center">
                          <div class="flex-shrink-0">
                            <div :class="[
                              'w-8 h-8 rounded-full flex items-center justify-center',
                              room.is_available ? 'bg-purple-100' : 'bg-gray-200'
                            ]">
                              <span class="text-sm">{{ getRoomTypeIcon(room.room_type) }}</span>
                            </div>
                          </div>
                          <div class="ml-3 flex-1">
                            <p :class="[
                              'text-sm font-medium',
                              room.is_available ? 'text-gray-900' : 'text-gray-500'
                            ]">
                              {{ room.room_number }} - {{ room.room_name || room.room_type }}
                            </p>
                            <p class="text-xs text-gray-500">
                              Floor {{ room.floor }} • {{ room.department }}
                            </p>
                            <p v-if="room.is_available" class="text-xs text-green-600">✓ Available</p>
                            <p v-else class="text-xs text-red-600">✗ Not Available</p>
                          </div>
                          <div v-if="room.is_available && bookingForm.room_id === room.id" class="flex-shrink-0">
                            <CheckIcon class="h-5 w-5 text-indigo-600" />
                          </div>
                          <div v-if="!room.is_available" class="flex-shrink-0">
                            <XMarkIcon class="h-5 w-5 text-red-400" />
                          </div>
                        </div>
                      </button>
                    </div>
                  </div>
                </div>

                <!-- Step 3: Confirmation -->
                <div v-else-if="currentStep === 3">
                  <!-- Conflict Display -->
                  <div v-if="bookingResult && !bookingResult.success" class="mb-6">
                    <div class="rounded-md bg-yellow-50 p-4">
                      <div class="flex">
                        <ExclamationTriangleIcon class="h-5 w-5 text-yellow-400" />
                        <div class="ml-3">
                          <h3 class="text-sm font-medium text-yellow-800">
                            Scheduling Conflict Detected
                          </h3>
                          <div class="mt-2 text-sm text-yellow-700">
                            <p>{{ bookingResult.message }}</p>
                            <ul v-if="bookingResult.conflicts" class="mt-2 list-disc list-inside">
                              <li v-for="conflict in bookingResult.conflicts" :key="conflict.conflict_type">
                                {{ conflict.description }}
                              </li>
                            </ul>
                          </div>
                        </div>
                      </div>
                    </div>

                    <!-- Alternative Slots -->
                    <div v-if="bookingResult.alternative_slots && bookingResult.alternative_slots.length > 0" class="mt-4">
                      <h4 class="text-sm font-medium text-gray-900 mb-3">Suggested Alternative Times:</h4>
                      <div class="grid grid-cols-2 lg:grid-cols-3 gap-3">
                        <button
                          v-for="alt in bookingResult.alternative_slots.slice(0, 6)"
                          :key="alt.date_time"
                          @click="selectAlternativeSlot(alt)"
                          class="p-3 text-left border border-gray-300 rounded-lg hover:border-indigo-500 hover:bg-indigo-50"
                        >
                          <div class="text-sm font-medium text-gray-900">
                            {{ formatDate(alt.date_time) }}
                          </div>
                          <div class="text-sm text-gray-500">
                            {{ formatTime(alt.date_time) }}
                          </div>
                          <div class="text-xs text-gray-400 capitalize">
                            {{ alt.slot_type }}
                          </div>
                        </button>
                      </div>
                    </div>
                  </div>

                  <!-- Success Display -->
                  <div v-else-if="bookingResult && bookingResult.success" class="mb-6">
                    <div class="rounded-md bg-green-50 p-4">
                      <div class="flex">
                        <CheckCircleIcon class="h-5 w-5 text-green-400" />
                        <div class="ml-3">
                          <h3 class="text-sm font-medium text-green-800">
                            Appointment Booked Successfully
                          </h3>
                          <div class="mt-2 text-sm text-green-700">
                            <p>{{ bookingResult.message }}</p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Booking Summary -->
                  <div class="bg-gray-50 rounded-lg p-4">
                    <h4 class="text-base font-medium text-gray-900 mb-4">Appointment Summary</h4>
                    <div class="grid grid-cols-2 gap-4 text-sm">
                      <div>
                        <span class="font-medium text-gray-700">Patient:</span>
                        <p class="text-gray-900">{{ selectedPatientName }}</p>
                      </div>
                      <div>
                        <span class="font-medium text-gray-700">Doctor:</span>
                        <p class="text-gray-900">{{ selectedDoctorName }}</p>
                      </div>
                      <div>
                        <span class="font-medium text-gray-700">Date:</span>
                        <p class="text-gray-900">{{ formatDate(bookingForm.date) }}</p>
                      </div>
                      <div>
                        <span class="font-medium text-gray-700">Time:</span>
                        <p class="text-gray-900">{{ bookingForm.time_slot }}</p>
                      </div>
                      <div>
                        <span class="font-medium text-gray-700">Type:</span>
                        <p class="text-gray-900 capitalize">{{ bookingForm.type }}</p>
                      </div>
                      <div>
                        <span class="font-medium text-gray-700">Duration:</span>
                        <p class="text-gray-900">{{ bookingForm.duration }} minutes</p>
                      </div>
                      <div>
                        <span class="font-medium text-gray-700">Room:</span>
                        <p class="text-gray-900">{{ selectedRoomName }}</p>
                      </div>
                      <div class="col-span-2">
                        <span class="font-medium text-gray-700">Reason:</span>
                        <p class="text-gray-900">{{ bookingForm.reason }}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div class="mt-8 sm:mt-6 sm:flex sm:flex-row-reverse">
                <button
                  v-if="currentStep < 3"
                  @click="nextStep"
                  :disabled="!canProceedToNextStep"
                  type="button"
                  class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:ml-3 sm:w-auto disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ currentStep === 2 ? 'Review & Book' : 'Continue' }}
                </button>
                
                <button
                  v-if="currentStep === 3 && bookingResult && bookingResult.success"
                  @click="$emit('appointment-booked')"
                  type="button"
                  class="inline-flex w-full justify-center rounded-md bg-green-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-green-500 sm:ml-3 sm:w-auto"
                >
                  Done
                </button>

                <button
                  v-if="currentStep === 3 && bookingResult && !bookingResult.success"
                  @click="bookWithAlternative"
                  :disabled="!selectedAlternative"
                  type="button"
                  class="inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 sm:ml-3 sm:w-auto disabled:opacity-50"
                >
                  Book Alternative
                </button>

                <button
                  v-if="currentStep > 1"
                  @click="previousStep"
                  type="button"
                  class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
                >
                  Back
                </button>

                <button
                  @click="close"
                  type="button"
                  class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
                >
                  Cancel
                </button>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'
import {
  CalendarDaysIcon,
  CheckIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
  ClockIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'
import { appointmentsApi } from '@/api/appointments'
import { patientsApi } from '@/api/patients'
import SearchableSelect from '@/components/SearchableSelect.vue'
import { 
  getTodayString, 
  parseDate, 
  formatDate, 
  formatDateForInput,
  getShortDayOfWeek,
  calculateAge as calcAge,
  formatDateTime,
  formatTimeFromDateTime,
  formatTime24
} from '@/utils/dateUtils'
import { convertUTCToEntityDate, createUTCDateTime } from '@/utils/timezoneUtils'
import { useAuthStore } from '@/stores/auth'
import { useEntityStore } from '@/stores/entity'
import { useFormConfig } from '@/composables/useFormConfig'

// Props & Emits
defineProps({
  show: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['close', 'appointment-booked'])

// Store access
const authStore = useAuthStore()
const entityStore = useEntityStore()

// Form configuration
const appointmentFormConfig = useFormConfig('appointment')

// Reactive state
const currentStep = ref(1)
const patients = ref([])
const doctors = ref([])
const availableSlots = ref([])
const availableDates = ref([])
const bookingResult = ref(null)
const selectedAlternative = ref(null)
const loadingSlots = ref(false)
const loadingAvailableDates = ref(false)
const loadingPatients = ref(false)
const loadingRooms = ref(false)
const loadingDurationOptions = ref(false)
const availableRooms = ref([])
const availableDurationOptions = ref([])
const requireRoomAssignment = ref(false)
const loadingRoomRequirement = ref(false)
const entityTimezone = ref('UTC')

const bookingForm = reactive({})

// Initialize booking form data based on configuration
const initializeBookingForm = () => {
  if (appointmentFormConfig.enabledFields.value.length > 0) {
    const initialData = appointmentFormConfig.initializeFormData()
    Object.assign(bookingForm, {
      ...initialData,
      // Always include these appointment-specific fields for booking process
      date: '',
      time_slot: '',
      room_id: '',
      check_conflicts: true
    })
  } else {
    // Fallback for when configuration is not loaded
    Object.assign(bookingForm, {
      patient_id: '',
      doctor_id: '',
      date: '',
      time_slot: '',
      duration: 30,
      type: '',
      reason: '',
      notes: '',
      priority: 'normal',
      room_id: '',
      check_conflicts: true
    })
  }
}

// Helper functions for form configuration
const isFieldVisible = (fieldName) => {
  return appointmentFormConfig.isFieldEnabled(fieldName)
}

const isFieldRequired = (fieldName) => {
  return appointmentFormConfig.isFieldRequired(fieldName)
}

const getFieldDisplayName = (fieldName) => {
  return appointmentFormConfig.getFieldDisplayName(fieldName) || fieldName.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const getFieldsByCategory = () => {
  return appointmentFormConfig.fieldsByCategory.value
}

// Computed
const today = computed(() => getTodayString())

const steps = computed(() => [
  { name: 'Details', status: currentStep.value > 1 ? 'complete' : currentStep.value === 1 ? 'current' : 'upcoming' },
  { name: 'Schedule & Room', status: currentStep.value > 2 ? 'complete' : currentStep.value === 2 ? 'current' : 'upcoming' },
  { name: 'Confirm', status: currentStep.value === 3 ? 'current' : 'upcoming' }
])

const canProceedToNextStep = computed(() => {
  if (currentStep.value === 1) {
    // Use dynamic form validation for step 1
    if (appointmentFormConfig.enabledFields.value.length > 0) {
      const requiredFields = appointmentFormConfig.requiredFields.value
      const step1RequiredFields = requiredFields.filter(field => 
        ['patient_id', 'doctor_id', 'type', 'reason', 'duration', 'priority', 'notes'].includes(field.name)
      )
      
      return step1RequiredFields.every(field => {
        const value = bookingForm[field.name]
        return value !== null && value !== undefined && value !== ''
      })
    } else {
      // Fallback validation
      return bookingForm.patient_id && bookingForm.doctor_id && bookingForm.type && bookingForm.reason
    }
  }
  if (currentStep.value === 2) {
    const hasDateAndTime = bookingForm.date && bookingForm.time_slot
    
    // If room assignment is required, check if a room is selected AND available
    if (requireRoomAssignment.value) {
      const selectedRoom = availableRooms.value.find(room => room.id == bookingForm.room_id)
      return hasDateAndTime && bookingForm.room_id && selectedRoom?.is_available
    }
    
    // If room is selected (optional), ensure it's available
    if (bookingForm.room_id) {
      const selectedRoom = availableRooms.value.find(room => room.id == bookingForm.room_id)
      return hasDateAndTime && selectedRoom?.is_available
    }
    
    return hasDateAndTime
  }
  return false
})

const selectedPatientName = computed(() => {
  const patient = patients.value.find(p => p.id == bookingForm.patient_id)
  return patient ? `${patient.first_name} ${patient.last_name}` : ''
})

const selectedDoctorName = computed(() => {
  const doctor = doctors.value.find(d => d.id == bookingForm.doctor_id)
  return doctor ? `Dr. ${doctor.first_name} ${doctor.last_name}` : ''
})

const selectedRoomName = computed(() => {
  if (!bookingForm.room_id) {
    return 'Any available room'
  }
  const room = availableRooms.value.find(r => r.id == bookingForm.room_id)
  return room ? `${room.room_number} - ${room.room_name || room.room_type}` : 'Selected room'
})

const groupedSlots = computed(() => {
  return availableSlots.value.reduce((groups, slot) => {
    const type = slot.slot_type || 'morning'
    if (!groups[type]) groups[type] = []
    groups[type].push(slot)
    return groups
  }, {})
})

const groupedAvailableDates = computed(() => {
  const todayString = getTodayString()
  const groups = {}
  
  availableDates.value.forEach(dateStr => {
    const date = parseDate(dateStr)
    const monthKey = date.toLocaleDateString('en-US', { year: 'numeric', month: 'long' })
    
    if (!groups[monthKey]) {
      groups[monthKey] = []
    }
    
    groups[monthKey].push({
      date: dateStr,
      day: date.getDate(),
      weekday: getShortDayOfWeek(dateStr),
      isToday: dateStr === todayString
    })
  })
  
  // Convert to array and sort by month
  return Object.entries(groups).map(([month, dates]) => ({
    month,
    dates: dates.sort((a, b) => a.date.localeCompare(b.date))
  })).sort((a, b) => a.dates[0].date.localeCompare(b.dates[0].date))
})

// Methods
const close = () => {
  emit('close')
}

const nextStep = async () => {
  if (currentStep.value === 1) {
    // Moving to step 2 - load available dates if doctor is selected
    if (bookingForm.doctor_id) {
      await loadAvailableDates()
    }
  } else if (currentStep.value === 2) {
    // Book the appointment
    await bookAppointment()
  }
  currentStep.value++
}

const previousStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

const onDoctorChange = async () => {
  // Reset date and time when doctor changes
  bookingForm.date = ''
  bookingForm.time_slot = ''
  availableSlots.value = []
  availableDates.value = []
  
  // Load available dates for the selected doctor
  if (bookingForm.doctor_id) {
    await loadAvailableDates()
  }
}

const onDateChange = async () => {
  if (bookingForm.date && bookingForm.doctor_id) {
    await loadAvailableSlots()
  }
}

const loadAvailableSlots = async () => {
  try {
    loadingSlots.value = true
    const response = await appointmentsApi.getAvailableTimeSlots(
      bookingForm.doctor_id,
      bookingForm.date,
      bookingForm.duration
    )
    availableSlots.value = response.data.slots || []
  } catch (err) {
    console.error('Load slots error:', err)
    availableSlots.value = []
  } finally {
    loadingSlots.value = false
  }
}

const loadAvailableDates = async () => {
  try {
    loadingAvailableDates.value = true
    const response = await appointmentsApi.getDoctorAvailability({
      doctor_id: bookingForm.doctor_id,
      status: 'available'
    })
    
    // Extract unique dates from availability records using entity timezone
    if (response.data && Array.isArray(response.data)) {
      const entityTimezone = entityStore.entityTimezone
      if (!entityTimezone) {
        console.error('Healthcare entity timezone not available')
        availableDates.value = []
        return
      }
      
      
      const dates = response.data
        .filter(availability => availability.start_datetime) // Only process records with datetime
        .map(availability => {
          // Convert UTC datetime to entity timezone date
          try {
            return convertUTCToEntityDate(availability.start_datetime, entityTimezone)
          } catch (error) {
            console.warn('Failed to convert date for availability:', availability, error)
            return null
          }
        })
        .filter(date => date) // Remove any invalid dates
      
      const uniqueDates = [...new Set(dates)]
      
      
      // Filter out past dates
      const todayString = getTodayString()
      availableDates.value = uniqueDates
        .filter(date => date >= todayString)
        .sort((a, b) => a.localeCompare(b))
        
    } else {
      availableDates.value = []
    }
  } catch (err) {
    console.error('Load available dates error:', err)
    availableDates.value = []
  } finally {
    loadingAvailableDates.value = false
  }
}

const selectDate = async (date) => {
  bookingForm.date = date
  bookingForm.time_slot = ''
  availableSlots.value = []
  
  // Load time slots for the selected date
  if (bookingForm.doctor_id && bookingForm.date) {
    await loadAvailableSlots()
  }
}

const selectTimeSlot = async (slot) => {
  if (slot.is_available) {
    bookingForm.time_slot = formatTime24(slot.date_time)
    // Reset room selection when time changes
    bookingForm.room_id = ''
    // Load available rooms for the selected time
    await loadAvailableRooms()
  }
}

const selectAlternativeSlot = (alt) => {
  selectedAlternative.value = alt
  bookingForm.date = formatDateForInput(new Date(alt.date_time))
  bookingForm.time_slot = formatTime24(alt.date_time)
}

const bookAppointment = async () => {
  try {
    // Get healthcare entity timezone
    const entityTimezone = entityStore.entityTimezone
    if (!entityTimezone) {
      throw new Error('Healthcare entity timezone not available')
    }
    
    // Properly convert entity time to UTC (same pattern as availability)
    const dateTime = createUTCDateTime(bookingForm.date, bookingForm.time_slot, entityTimezone)
    
    const bookingData = {
      patient_id: parseInt(bookingForm.patient_id),
      doctor_id: parseInt(bookingForm.doctor_id),
      date_time: dateTime, // Send UTC ISO string like availability does
      duration: parseInt(bookingForm.duration),
      type: bookingForm.type,
      reason: bookingForm.reason,
      notes: bookingForm.notes,
      priority: bookingForm.priority,
      room_id: parseInt(bookingForm.room_id) || 0,
      check_conflicts: bookingForm.check_conflicts
    }

    const response = await appointmentsApi.bookAppointment(bookingData)
    bookingResult.value = response.data
  } catch (err) {
    console.error('Booking error:', err)
    bookingResult.value = {
      success: false,
      message: 'Failed to book appointment: ' + (err.message || 'Unknown error'),
      conflicts: [],
      alternative_slots: []
    }
  }
}

const bookWithAlternative = async () => {
  if (selectedAlternative.value) {
    await bookAppointment()
  }
}

const loadPatients = async () => {
  try {
    loadingPatients.value = true
    const response = await patientsApi.getPatients({ limit: 100 })
    patients.value = response.patients || []
  } catch (err) {
    console.error('Load patients error:', err)
    patients.value = []
  } finally {
    loadingPatients.value = false
  }
}

const formatPatientInfo = (patient) => {
  const info = []
  if (patient.patient_id) info.push(`ID: ${patient.patient_id}`)
  
  // Use pre-calculated age from backend if available, otherwise calculate from date_of_birth
  if (patient.age && patient.age > 0) {
    info.push(`${patient.age} years old`)
  } else if (patient.date_of_birth) {
    // Extract date part from ISO string (YYYY-MM-DDTHH:MM:SSZ -> YYYY-MM-DD)
    const dateOnly = patient.date_of_birth.split('T')[0]
    const age = calculateAge(dateOnly)
    if (!isNaN(age) && age >= 0) {
      info.push(`${age} years old`)
    }
  }
  
  if (patient.phone) info.push(formatPhone(patient.phone))
  return info.join(' • ')
}

const calculateAge = (dateOfBirth) => {
  return calcAge(dateOfBirth)
}

const formatPhone = (phone) => {
  if (!phone) return ''
  
  const cleaned = phone.replace(/\D/g, '')
  
  if (cleaned.length === 10) {
    return `(${cleaned.slice(0, 3)}) ${cleaned.slice(3, 6)}-${cleaned.slice(6)}`
  } else if (cleaned.length === 11 && cleaned.startsWith('1')) {
    return `+1 (${cleaned.slice(1, 4)}) ${cleaned.slice(4, 7)}-${cleaned.slice(7)}`
  }
  
  return phone
}

const onPatientSearch = (query) => {
  // Handle patient search if needed (could trigger API search)
}

const onPatientSelected = (patient) => {
  // Handle patient selection if additional logic is needed
}

const loadAvailableRooms = async () => {
  if (!bookingForm.date || !bookingForm.time_slot || !bookingForm.duration) {
    availableRooms.value = []
    return
  }

  try {
    loadingRooms.value = true
    
    // Get healthcare entity timezone (same pattern as appointment booking)
    const entityTimezone = entityStore.entityTimezone
    if (!entityTimezone) {
      throw new Error('Healthcare entity timezone not available')
    }
    
    // Convert entity time to UTC (same pattern as appointment booking)
    const dateTime = createUTCDateTime(bookingForm.date, bookingForm.time_slot, entityTimezone)
    
    const response = await appointmentsApi.getAvailableRooms(dateTime, bookingForm.duration)
    availableRooms.value = response.data?.available_rooms || []
  } catch (err) {
    console.error('Load available rooms error:', err)
    availableRooms.value = []
  } finally {
    loadingRooms.value = false
  }
}


const getRoomTypeIcon = (roomType) => {
  const icons = {
    'consultation': '💬',
    'examination': '🔍',
    'procedure': '⚕️',
    'operating': '🏥',
    'emergency': '🚨'
  }
  return icons[roomType] || '🏠'
}

const loadDurationOptions = async (appointmentType) => {
  if (!appointmentType) {
    availableDurationOptions.value = []
    return
  }

  try {
    loadingDurationOptions.value = true
    const response = await appointmentsApi.getDurationOptions(appointmentType)
    availableDurationOptions.value = response.data || []
    
    // Auto-select default duration if available
    const defaultOption = availableDurationOptions.value.find(option => option.is_default)
    if (defaultOption && !bookingForm.duration) {
      bookingForm.duration = defaultOption.duration_minutes
    }
  } catch (err) {
    console.error('Load duration options error:', err)
    availableDurationOptions.value = []
  } finally {
    loadingDurationOptions.value = false
  }
}

const loadDoctors = async () => {
  try {
    const response = await appointmentsApi.getDoctorsByEntity()
    doctors.value = response.data || []
  } catch (err) {
    console.error('Load doctors error:', err)
  }
}


const loadRoomRequirement = () => {
  try {
    loadingRoomRequirement.value = true
    // Use entity store instead of separate API call
    requireRoomAssignment.value = entityStore.entityRequireRoomAssignment
  } catch (err) {
    console.error('Load room requirement error:', err)
    requireRoomAssignment.value = false // Default to not required
  } finally {
    loadingRoomRequirement.value = false
  }
}

// Utility functions
const formatDateLocal = (dateString) => {
  return formatDate(dateString)
}

const formatTime = (dateTimeString) => {
  // For healthcare applications, display times in the facility's timezone
  // rather than the user's browser timezone for consistency
  const date = new Date(dateTimeString)
  
  // Use the healthcare entity's timezone from the store/API
  const timezone = entityTimezone.value || 'UTC'
  
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit', 
    hour12: true,
    timeZone: timezone
  })
}

// Watchers
watch(() => bookingForm.type, (newType) => {
  // Reset duration when appointment type changes
  bookingForm.duration = ''
  loadDurationOptions(newType)
})

watch(() => requireRoomAssignment.value, (isRequired) => {
  // If room assignment becomes required and no room is selected, 
  // ensure the room field is cleared to trigger validation
  if (isRequired && !bookingForm.room_id) {
    bookingForm.room_id = ''
  }
})


// Lifecycle
onMounted(async () => {
  try {
    // Initialize form configuration first
    await appointmentFormConfig.initialize()
    initializeBookingForm()
  } catch (err) {
    console.warn('Failed to load appointment form configuration, using defaults:', err)
    initializeBookingForm()
  }
  
  await Promise.all([loadPatients(), loadDoctors(), loadRoomRequirement()])
})
</script>