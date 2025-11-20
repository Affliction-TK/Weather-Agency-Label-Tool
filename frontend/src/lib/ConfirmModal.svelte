<script>
  import { createEventDispatcher } from 'svelte';
  import { fade, scale } from 'svelte/transition';

  export let title = '确认';
  export let message = '您确定要执行此操作吗？';
  export let confirmText = '确认';
  export let cancelText = '取消';
  export let type = 'danger'; // 'danger' | 'primary'

  const dispatch = createEventDispatcher();

  function handleConfirm() {
    dispatch('confirm');
  }

  function handleCancel() {
    dispatch('cancel');
  }

  // 新增：处理键盘事件
  function handleKeydown(event) {
    if (event.key === 'Enter' || event.key === ' ') {
      handleCancel();
    }
  }
</script>

<!-- 修改：添加 role, tabindex 和 on:keydown -->
<div 
  class="modal-backdrop" 
  transition:fade={{ duration: 200 }} 
  on:click|self={handleCancel}
  on:keydown|self={handleKeydown}
  role="button"
  tabindex="0"
  aria-label="关闭模态框"
>
  <div class="modal-content" transition:scale={{ duration: 200, start: 0.95 }} role="document">
    <div class="modal-header">
      <h3>{title}</h3>
    </div>
    <div class="modal-body">
      <p>{message}</p>
    </div>
    <div class="modal-footer">
      <button class="btn-cancel" on:click={handleCancel}>{cancelText}</button>
      <button class="btn-confirm {type}" on:click={handleConfirm}>{confirmText}</button>
    </div>
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
  }

  .modal-content {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    width: 90%;
    max-width: 400px;
    border-radius: 16px;
    box-shadow: 0 20px 50px rgba(0, 0, 0, 0.2);
    overflow: hidden;
    border: 1px solid rgba(255, 255, 255, 0.5);
  }

  .modal-header {
    padding: 20px 24px 10px;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #1a1a1a;
  }

  .modal-body {
    padding: 0 24px 24px;
  }

  .modal-body p {
    margin: 0;
    color: #666;
    font-size: 15px;
    line-height: 1.5;
  }

  .modal-footer {
    padding: 16px 24px;
    background: rgba(0, 0, 0, 0.02);
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    border-top: 1px solid rgba(0, 0, 0, 0.05);
  }

  button {
    padding: 10px 20px;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    border: none;
    transition: all 0.2s;
  }

  .btn-cancel {
    background: white;
    color: #666;
    border: 1px solid #e5e5e5;
  }

  .btn-cancel:hover {
    background: #f5f5f5;
    color: #333;
  }

  .btn-confirm {
    color: white;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .btn-confirm.primary {
    background: #007aff;
  }

  .btn-confirm.primary:hover {
    background: #0062cc;
  }

  .btn-confirm.danger {
    background: #ff3b30;
  }

  .btn-confirm.danger:hover {
    background: #d70015;
  }
</style>
