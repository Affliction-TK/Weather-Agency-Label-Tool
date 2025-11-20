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
    top: 24px;
    right: 24px;
    min-width: 320px;
    max-width: 480px;
    padding: 16px;
    border-radius: 14px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15), 0 0 0 1px rgba(0,0,0,0.05);
    display: flex;
    align-items: flex-start;
    gap: 12px;
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
    transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
    cursor: pointer;
    z-index: 10000;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
  }

  .toast.visible {
    opacity: 1;
    transform: translateY(0) scale(1);
  }

  .toast.success {
    border-left: 4px solid #34C759;
  }

  .toast.error {
    border-left: 4px solid #FF3B30;
  }

  .toast.warning {
    border-left: 4px solid #FF9500;
  }

  .toast.info {
    border-left: 4px solid #007AFF;
  }

  .toast-icon {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: bold;
    flex-shrink: 0;
    margin-top: 2px;
  }

  .toast.success .toast-icon { background: #34C759; color: white; }
  .toast.error .toast-icon { background: #FF3B30; color: white; }
  .toast.warning .toast-icon { background: #FF9500; color: white; }
  .toast.info .toast-icon { background: #007AFF; color: white; }

  .toast-message {
    flex: 1;
    font-size: 14px;
    font-weight: 500;
    line-height: 1.4;
    color: #1D1D1F;
    padding-top: 2px;
  }

  .toast-close {
    width: 20px;
    height: 20px;
    border: none;
    background: transparent;
    color: #86868B;
    border-radius: 50%;
    cursor: pointer;
    font-size: 18px;
    line-height: 1;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    margin-top: 2px;
  }

  .toast-close:hover {
    background: rgba(0, 0, 0, 0.05);
    color: #1D1D1F;
  }
</style>
