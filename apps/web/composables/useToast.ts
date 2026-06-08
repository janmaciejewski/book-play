const toast = reactive({
  message: '',
  type: '' as '' | 'success' | 'error',
})

const toastClass = computed(() => {
  if (toast.type === 'success') {
    return 'fixed bottom-4 right-4 px-4 py-3 rounded-lg shadow-lg z-50 text-sm bg-green-500 text-white'
  }
  return 'fixed bottom-4 right-4 px-4 py-3 rounded-lg shadow-lg z-50 text-sm bg-red-500 text-white'
})

let timer: ReturnType<typeof setTimeout> | null = null

function showToast(message: string, type: 'success' | 'error' = 'success', duration = 3000) {
  if (timer) clearTimeout(timer)
  toast.message = message
  toast.type = type
  timer = setTimeout(() => {
    toast.message = ''
  }, duration)
}

export function useToast() {
  return { toast, toastClass, showToast }
}
