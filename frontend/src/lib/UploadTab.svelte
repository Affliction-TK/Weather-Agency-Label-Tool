<script>
  import { createEventDispatcher } from 'svelte';
  import { toasts } from './toastStore.js';
  
  const dispatch = createEventDispatcher();
  const API_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080/api' : '/api';

  let files = [];
  let uploading = false;
  let dragActive = false;

  function handleFileInput(event) {
    files = Array.from(event.target.files);
  }

  function handleDrop(event) {
    event.preventDefault();
    dragActive = false;
    files = Array.from(event.dataTransfer.files).filter(file => file.type.startsWith('image/'));
  }

  function handleDragOver(event) {
    event.preventDefault();
    dragActive = true;
  }

  function handleDragLeave() {
    dragActive = false;
  }

  async function uploadFiles() {
    if (files.length === 0) {
      toasts.warning('请选择要上传的图片');
      return;
    }

    uploading = true;
    let successCount = 0;
    let failCount = 0;

    for (const file of files) {
      try {
        const formData = new FormData();
        formData.append('image', file);

        const response = await fetch(`${API_BASE}/upload`, {
          method: 'POST',
          body: formData
        });

        if (response.ok) {
          successCount++;
        } else {
          failCount++;
        }
      } catch (error) {
        console.error('Upload failed:', error);
        failCount++;
      }
    }

    uploading = false;
    files = [];

    if (successCount > 0) {
      toasts.success(`成功上传 ${successCount} 张图片${failCount > 0 ? `，${failCount} 张失败` : ''}`);
      dispatch('uploaded');
    } else {
      toasts.error('上传失败，请重试');
    }
  }

  function removeFile(index) {
    files = files.filter((_, i) => i !== index);
  }
</script>

<div class="upload-tab">
  <div class="upload-container">
    <h2>上传新图片</h2>
    
    <div 
      class="drop-zone"
      class:active={dragActive}
      on:drop={handleDrop}
      on:dragover={handleDragOver}
      on:dragleave={handleDragLeave}
    >
      <div class="drop-zone-content">
        <svg class="upload-icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>
        <p>拖拽图片到此处，或点击选择文件</p>
        <input 
          type="file" 
          accept="image/*" 
          multiple 
          on:change={handleFileInput}
          id="fileInput"
        />
        <label for="fileInput" class="file-button">选择文件</label>
      </div>
    </div>

    {#if files.length > 0}
      <div class="file-list">
        <h3>选中的文件 ({files.length})</h3>
        {#each files as file, index (file.name + index)}
          <div class="file-item">
            <span class="file-name">{file.name}</span>
            <span class="file-size">{(file.size / 1024).toFixed(1)} KB</span>
            <button 
              type="button" 
              class="remove-btn"
              on:click={() => removeFile(index)}
            >
              ×
            </button>
          </div>
        {/each}
      </div>

      <div class="upload-actions">
        <button 
          type="button" 
          class="upload-btn"
          disabled={uploading}
          on:click={uploadFiles}
        >
          {#if uploading}
            上传中...
          {:else}
            上传 {files.length} 张图片
          {/if}
        </button>
      </div>
    {/if}
  </div>
</div>

<style>
  .upload-tab {
    height: 100%;
    overflow-y: auto;
    padding: 32px;
    background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  }

  .upload-container {
    max-width: 800px;
    margin: 0 auto;
  }

  h2 {
    color: #333;
    margin-bottom: 28px;
    font-size: 28px;
    font-weight: 700;
    text-align: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .drop-zone {
    border: 3px dashed #ccc;
    border-radius: 16px;
    padding: 60px 20px;
    text-align: center;
    background: white;
    transition: all 0.3s;
    cursor: pointer;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  }

  .drop-zone.active {
    border-color: #667eea;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
    transform: scale(1.02);
  }

  .drop-zone-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
  }

  .upload-icon {
    width: 80px;
    height: 80px;
    color: #667eea;
  }

  .drop-zone p {
    color: #666;
    font-size: 16px;
    margin: 0;
    font-weight: 500;
  }

  #fileInput {
    display: none;
  }

  .file-button {
    padding: 14px 32px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border-radius: 12px;
    cursor: pointer;
    font-size: 15px;
    font-weight: 600;
    transition: all 0.3s;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  }

  .file-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
  }

  .file-list {
    margin-top: 32px;
    background: white;
    border-radius: 16px;
    padding: 24px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  }

  .file-list h3 {
    margin-top: 0;
    margin-bottom: 20px;
    color: #333;
    font-size: 18px;
    font-weight: 700;
  }

  .file-item {
    display: flex;
    align-items: center;
    padding: 14px;
    border-bottom: 1px solid #f0f0f0;
    transition: background 0.2s;
    border-radius: 8px;
  }

  .file-item:hover {
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  }

  .file-item:last-child {
    border-bottom: none;
  }

  .file-name {
    flex: 1;
    color: #333;
    font-size: 14px;
    font-weight: 500;
  }

  .file-size {
    color: #999;
    font-size: 13px;
    margin-right: 16px;
    font-weight: 500;
  }

  .remove-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);
    color: white;
    border-radius: 50%;
    cursor: pointer;
    font-size: 20px;
    line-height: 1;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    box-shadow: 0 2px 8px rgba(235, 51, 73, 0.3);
  }

  .remove-btn:hover {
    transform: scale(1.1);
    box-shadow: 0 4px 12px rgba(235, 51, 73, 0.4);
  }

  .upload-actions {
    margin-top: 24px;
    text-align: center;
  }

  .upload-btn {
    padding: 16px 48px;
    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
    color: white;
    border: none;
    border-radius: 12px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
    box-shadow: 0 4px 12px rgba(17, 153, 142, 0.3);
  }

  .upload-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(17, 153, 142, 0.4);
  }

  .upload-btn:active:not(:disabled) {
    transform: translateY(0);
  }

  .upload-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }
</style>
