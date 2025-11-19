import { writable } from 'svelte/store';

function createToastStore() {
  const { subscribe, update } = writable([]);
  let id = 0;

  return {
    subscribe,
    show: (message, type = 'success', duration = 3000) => {
      const toastId = id++;
      update(toasts => [...toasts, { id: toastId, message, type, duration }]);
      
      if (duration > 0) {
        setTimeout(() => {
          update(toasts => toasts.filter(t => t.id !== toastId));
        }, duration + 300); // Add extra time for animation
      }
    },
    success: (message, duration = 3000) => {
      const toastId = id++;
      update(toasts => [...toasts, { id: toastId, message, type: 'success', duration }]);
      if (duration > 0) {
        setTimeout(() => {
          update(toasts => toasts.filter(t => t.id !== toastId));
        }, duration + 300);
      }
    },
    error: (message, duration = 4000) => {
      const toastId = id++;
      update(toasts => [...toasts, { id: toastId, message, type: 'error', duration }]);
      if (duration > 0) {
        setTimeout(() => {
          update(toasts => toasts.filter(t => t.id !== toastId));
        }, duration + 300);
      }
    },
    info: (message, duration = 3000) => {
      const toastId = id++;
      update(toasts => [...toasts, { id: toastId, message, type: 'info', duration }]);
      if (duration > 0) {
        setTimeout(() => {
          update(toasts => toasts.filter(t => t.id !== toastId));
        }, duration + 300);
      }
    },
    warning: (message, duration = 3000) => {
      const toastId = id++;
      update(toasts => [...toasts, { id: toastId, message, type: 'warning', duration }]);
      if (duration > 0) {
        setTimeout(() => {
          update(toasts => toasts.filter(t => t.id !== toastId));
        }, duration + 300);
      }
    },
    remove: (toastId) => {
      update(toasts => toasts.filter(t => t.id !== toastId));
    }
  };
}

export const toasts = createToastStore();
