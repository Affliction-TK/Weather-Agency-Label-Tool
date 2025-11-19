<script>
  import { onMount } from 'svelte';
  
  export let message = '';
  export let type = 'success'; // success, error, info, warning
  export let duration = 3000;
  export let onClose = () => {};
  
  let visible = false;
  
  onMount(() => {
    visible = true;
    
    if (duration > 0) {
      setTimeout(() => {
        close();
      }, duration);
    }
  });
  
  function close() {
    visible = false;
    setTimeout(() => {
      onClose();
    }, 300); // Wait for animation to finish
  }
  
  const icons = {
    success: '✓',
    error: '✕',
    info: 'ℹ',
    warning: '⚠'
  };
</script>

<div class="toast {type}" class:visible on:click={close}>
  <div class="toast-icon">{icons[type]}</div>
  <div class="toast-message">{message}</div>
  <button class="toast-close" on:click={close}>×</button>
</div>

<style>
  .toast {
    position: fixed;
    top: 20px;
    right: 20px;
    min-width: 300px;
    max-width: 500px;
    padding: 16px 20px;
    border-radius: 12px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    display: flex;
    align-items: center;
    gap: 12px;
    opacity: 0;
    transform: translateX(400px);
    transition: all 0.3s cubic-bezier(0.68, -0.55, 0.265, 1.55);
    cursor: pointer;
    z-index: 10000;
    backdrop-filter: blur(10px);
  }

  .toast.visible {
    opacity: 1;
    transform: translateX(0);
  }

  .toast-icon {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    font-weight: bold;
    flex-shrink: 0;
  }

  .toast-message {
    flex: 1;
    font-size: 14px;
    font-weight: 500;
    line-height: 1.4;
  }

  .toast-close {
    width: 24px;
    height: 24px;
    border: none;
    background: rgba(0, 0, 0, 0.1);
    color: inherit;
    border-radius: 50%;
    cursor: pointer;
    font-size: 18px;
    line-height: 1;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    transition: background 0.2s;
  }

  .toast-close:hover {
    background: rgba(0, 0, 0, 0.2);
  }

  .toast.success {
    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
    color: white;
  }

  .toast.success .toast-icon {
    background: rgba(255, 255, 255, 0.25);
  }

  .toast.error {
    background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);
    color: white;
  }

  .toast.error .toast-icon {
    background: rgba(255, 255, 255, 0.25);
  }

  .toast.info {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .toast.info .toast-icon {
    background: rgba(255, 255, 255, 0.25);
  }

  .toast.warning {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
    color: white;
  }

  .toast.warning .toast-icon {
    background: rgba(255, 255, 255, 0.25);
  }

  @media (max-width: 640px) {
    .toast {
      right: 10px;
      left: 10px;
      min-width: auto;
    }
  }
</style>
