<template>
  <Transition name="modal" appear>
    <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto">
      <Transition name="backdrop" appear>
        <div 
          class="fixed inset-0 backdrop-blur-sm bg-black/30 transition-all duration-300" 
          @click="closeModal"
        ></div>
      </Transition>
      
      <div class="flex min-h-full items-center justify-center p-4">
        <Transition name="modal-content" appear>
          <div class="relative w-full max-w-2xl bg-gray-800/95 backdrop-blur-md rounded-xl shadow-2xl border border-gray-700/50 transform transition-all duration-300">

            <div class="flex items-center justify-between p-6 border-b border-gray-700/50">
              <h3 class="text-xl font-semibold text-white">
                <slot name="header">Modal Title</slot>
              </h3>
              <button
                @click="closeModal"
                class="text-gray-400 hover:text-white transition-colors p-1 rounded-full hover:bg-gray-700/50"
              >
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <div class="p-6">
              <slot name="content">Modal content goes here</slot>
            </div>
            
            <div v-if="$slots.footer" class="flex items-center justify-end p-6 border-t border-gray-700/50">
              <slot name="footer"></slot>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
interface Props {
  isOpen: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  close: []
}>()

const closeModal = () => {
  emit('close')
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.backdrop-enter-active,
.backdrop-leave-active {
  transition: all 0.3s ease;
}

.backdrop-enter-from,
.backdrop-leave-to {
  opacity: 0;
  backdrop-filter: blur(0px);
}

.modal-content-enter-active,
.modal-content-leave-active {
  transition: all 0.3s ease;
}

.modal-content-enter-from,
.modal-content-leave-to {
  opacity: 0;
  transform: scale(0.9) translateY(-20px);
}

.modal-content-enter-to,
.modal-content-leave-from {
  opacity: 1;
  transform: scale(1) translateY(0);
}
</style> 