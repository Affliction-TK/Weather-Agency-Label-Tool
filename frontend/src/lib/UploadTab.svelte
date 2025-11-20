<script>
  import { createEventDispatcher } from 'svelte';
  import { toasts } from './toastStore.js';
  
  const dispatch = createEventDispatcher();
  const API_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080/api' : '/api';

  let files = [];
  let uploading = false;
  let dragActive = false;

  function mergeFiles(fileList) {
    const incoming = Array.from(fileList).filter(file => file.type.startsWith('image/'));
    const seen = new Set(files.map(file => `${file.name}-${file.size}-${file.lastModified}`));
    const unique = [];

    for (const file of incoming) {
      const key = `${file.name}-${file.size}-${file.lastModified}`;
      if (!seen.has(key)) {
        seen.add(key);
        unique.push(file);
      }
    }

    if (unique.length > 0) {
      files = [...files, ...unique];
    }
  }

  function handleFileInput(event) {
    mergeFiles(event.target.files);
    event.target.value = '';
  }

  function handleDrop(event) {
    event.preventDefault();
    dragActive = false;
    mergeFiles(event.dataTransfer.files);
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
      role="region"
      aria-label="图片上传放置区"
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
    background: var(--content-bg);
  }

  .upload-container {
    max-width: 800px;
    margin: 0 auto;
  }

  h2 {
    color: var(--text-primary);
    margin-bottom: 24px;
    font-size: 24px;
    font-weight: 700;
    text-align: center;
    letter-spacing: -0.5px;
  }

  .drop-zone {
    border: 2px dashed var(--border-color);
    border-radius: var(--radius-m);
    padding: 60px 20px;
    text-align: center;
    background: var(--bg-color);
    transition: all 0.2s;
    cursor: pointer;
  }

  .drop-zone:hover {
    border-color: var(--primary-color);
    background: rgba(0, 122, 255, 0.02);
  }

  .drop-zone.active {
    border-color: var(--primary-color);
    background: rgba(0, 122, 255, 0.05);
    transform: scale(1.01);
  }

  .drop-zone-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
  }

  .upload-icon {
    width: 64px;
    height: 64px;
    color: var(--primary-color);
    opacity: 0.8;
  }

  .drop-zone p {
    color: var(--text-secondary);
    font-size: 15px;
    margin: 0;
    font-weight: 500;
  }

  #fileInput {
    display: none;
  }

  .file-button {
    padding: 10px 24px;
    background: white;
    color: var(--primary-color);
    border: 1px solid var(--primary-color);
    border-radius: 20px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: all 0.2s;
  }

  .file-button:hover {
    background: var(--primary-color);
    color: white;
  }

  .file-list {
    margin-top: 32px;
    background: transparent;
  }

  .file-list h3 {
    margin-top: 0;
    margin-bottom: 16px;
    color: var(--text-primary);
    font-size: 16px;
    font-weight: 600;
  }

  .file-item {
    display: flex;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid var(--divider-color);
    background: white;
    transition: background 0.2s;
  }

  .file-item:first-of-type {
    border-top-left-radius: var(--radius-s);
    border-top-right-radius: var(--radius-s);
  }

  .file-item:last-child {
    border-bottom: none;
    border-bottom-left-radius: var(--radius-s);
    border-bottom-right-radius: var(--radius-s);
  }

  .file-item:hover {
    background: #FAFAFA;
  }

  .file-name {
    flex: 1;
    color: var(--text-primary);
    font-size: 14px;
    font-weight: 500;
  }

  .file-size {
    color: var(--text-secondary);
    font-size: 13px;
    margin-right: 16px;
    font-weight: 400;
  }

  .remove-btn {
    width: 24px;
    height: 24px;
    border: none;
    background: rgba(0,0,0,0.1);
    color: white;
    border-radius: 50%;
    cursor: pointer;
    font-size: 16px;
    line-height: 1;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }

  .remove-btn:hover {
    background: var(--danger-color);
  }

  .upload-actions {
    margin-top: 32px;
    text-align: center;
  }

  .upload-btn {
    padding: 12px 40px;
    background: var(--success-color);
    color: white;
    border: none;
    border-radius: 24px;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 2px 8px rgba(52, 199, 89, 0.25);
  }

  .upload-btn:hover:not(:disabled) {
    background: #2DB84C;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(52, 199, 89, 0.35);
  }

  .upload-btn:active:not(:disabled) {
    transform: translateY(0);
  }

  .upload-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    box-shadow: none;
  }
</style>
