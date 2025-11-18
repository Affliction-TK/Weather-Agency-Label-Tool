<script>
  import { createEventDispatcher } from 'svelte';
  
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
      alert('请选择要上传的图片');
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
      alert(`成功上传 ${successCount} 张图片${failCount > 0 ? `，${failCount} 张失败` : ''}`);
      dispatch('uploaded');
    } else {
      alert('上传失败，请重试');
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
    padding: 20px;
    background: #f9f9f9;
  }

  .upload-container {
    max-width: 800px;
    margin: 0 auto;
  }

  h2 {
    color: #333;
    margin-bottom: 20px;
  }

  .drop-zone {
    border: 3px dashed #ddd;
    border-radius: 8px;
    padding: 60px 20px;
    text-align: center;
    background: white;
    transition: all 0.3s;
    cursor: pointer;
  }

  .drop-zone.active {
    border-color: #1976d2;
    background: #e3f2fd;
  }

  .drop-zone-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 15px;
  }

  .upload-icon {
    width: 64px;
    height: 64px;
    color: #999;
  }

  .drop-zone p {
    color: #666;
    font-size: 16px;
    margin: 0;
  }

  #fileInput {
    display: none;
  }

  .file-button {
    padding: 12px 30px;
    background: #1976d2;
    color: white;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: background 0.2s;
  }

  .file-button:hover {
    background: #1565c0;
  }

  .file-list {
    margin-top: 30px;
    background: white;
    border-radius: 8px;
    padding: 20px;
    border: 1px solid #ddd;
  }

  .file-list h3 {
    margin-top: 0;
    margin-bottom: 15px;
    color: #333;
    font-size: 16px;
  }

  .file-item {
    display: flex;
    align-items: center;
    padding: 10px;
    border-bottom: 1px solid #f0f0f0;
  }

  .file-item:last-child {
    border-bottom: none;
  }

  .file-name {
    flex: 1;
    color: #333;
    font-size: 14px;
  }

  .file-size {
    color: #999;
    font-size: 12px;
    margin-right: 10px;
  }

  .remove-btn {
    width: 24px;
    height: 24px;
    border: none;
    background: #ff5252;
    color: white;
    border-radius: 50%;
    cursor: pointer;
    font-size: 18px;
    line-height: 1;
    padding: 0;
  }

  .remove-btn:hover {
    background: #ff1744;
  }

  .upload-actions {
    margin-top: 20px;
    text-align: center;
  }

  .upload-btn {
    padding: 15px 40px;
    background: #4caf50;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
  }

  .upload-btn:hover:not(:disabled) {
    background: #45a049;
  }

  .upload-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
