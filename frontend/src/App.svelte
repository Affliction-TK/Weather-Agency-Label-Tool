<script>
  import { onMount } from 'svelte';
  import ImageList from './lib/ImageList.svelte';
  import AnnotationForm from './lib/AnnotationForm.svelte';
  import UploadTab from './lib/UploadTab.svelte';
  import Toast from './lib/Toast.svelte';
  import { toasts } from './lib/toastStore.js';

  let images = [];
  let stations = [];
  let currentImage = null;
  let currentAnnotation = null;
  let activeTab = 'annotate'; // 'annotate' or 'upload'
  let loading = true;

  const API_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080/api' : '/api';

  onMount(async () => {
    await loadStations();
    await loadImages();
    loading = false;
  });

  async function loadStations() {
    try {
      const response = await fetch(`${API_BASE}/stations`);
      stations = await response.json();
    } catch (error) {
      console.error('Failed to load stations:', error);
      toasts.error('加载监测站点失败');
    }
  }

  async function loadImages() {
    try {
      const response = await fetch(`${API_BASE}/images`);
      images = await response.json();
      
      // Auto-select first unannotated image
      if (images.length > 0) {
        const firstUnannotated = images.find(img => !img.annotated);
        if (firstUnannotated) {
          await selectImage(firstUnannotated);
        } else if (images.length > 0) {
          // All images annotated, select first one
          await selectImage(images[0]);
        }
      }
    } catch (error) {
      console.error('Failed to load images:', error);
      toasts.error('加载图片列表失败');
    }
  }

  async function selectImage(image) {
    try {
      const response = await fetch(`${API_BASE}/images/${image.id}`);
      const data = await response.json();
      currentImage = data.image;
      currentAnnotation = data.annotation || null;
      activeTab = 'annotate';
    } catch (error) {
      console.error('Failed to load image details:', error);
      toasts.error('加载图片详情失败');
    }
  }

  async function handleAnnotationSaved() {
    await loadImages();
    // Move to next unannotated image
    const nextUnannotated = images.find(img => !img.annotated && img.id !== currentImage.id);
    if (nextUnannotated) {
      await selectImage(nextUnannotated);
    }
  }

  async function handleImageUploaded() {
    await loadImages();
    activeTab = 'annotate';
  }

  function handleTabChange(tab) {
    activeTab = tab;
  }

  $: allAnnotated = images.length > 0 && images.every(img => img.annotated);
</script>

<!-- Toast Notifications -->
{#each $toasts as toast (toast.id)}
  <Toast 
    message={toast.message} 
    type={toast.type} 
    duration={toast.duration}
    onClose={() => toasts.remove(toast.id)}
  />
{/each}

<main>
  <div class="container">
    <div class="sidebar">
      <h2>图片列表</h2>
      {#if loading}
        <p>加载中...</p>
      {:else if images.length === 0}
        <p>暂无图片</p>
      {:else}
        <ImageList {images} {currentImage} on:select={(e) => selectImage(e.detail)} />
      {/if}
    </div>

    <div class="content">
      <div class="tabs">
        <button 
          class:active={activeTab === 'annotate'} 
          on:click={() => handleTabChange('annotate')}
          disabled={!currentImage}
        >
          标注
        </button>
        <button 
          class:active={activeTab === 'upload'} 
          on:click={() => handleTabChange('upload')}
        >
          上传图片
        </button>
      </div>

      {#if activeTab === 'annotate'}
        {#if allAnnotated && !currentImage}
          <div class="all-done">
            <h2>✓ 所有图片已标注完成</h2>
            <p>请切换到"上传图片"标签页以添加新图片</p>
          </div>
        {:else if currentImage}
          <AnnotationForm 
            image={currentImage} 
            annotation={currentAnnotation}
            {stations}
            on:saved={handleAnnotationSaved}
          />
        {:else}
          <div class="no-image">
            <p>请从左侧列表选择图片进行标注</p>
          </div>
        {/if}
      {:else if activeTab === 'upload'}
        <UploadTab on:uploaded={handleImageUploaded} />
      {/if}
    </div>
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  main {
    width: 100%;
    height: 100vh;
    overflow: hidden;
  }

  .container {
    display: flex;
    height: 100%;
    margin: 0;
  }

  .sidebar {
    width: 320px;
    background: #f8f9fa;
    border-right: 1px solid #e0e0e0;
    overflow-y: auto;
    padding: 20px;
    box-shadow: 2px 0 10px rgba(0, 0, 0, 0.05);
  }

  .sidebar h2 {
    margin-top: 0;
    margin-bottom: 20px;
    font-size: 20px;
    font-weight: 700;
    color: #333;
    text-align: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background: white;
  }

  .tabs {
    display: flex;
    border-bottom: 2px solid #e0e0e0;
    background: white;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  }

  .tabs button {
    padding: 16px 32px;
    border: none;
    background: none;
    cursor: pointer;
    font-size: 15px;
    font-weight: 600;
    color: #666;
    border-bottom: 3px solid transparent;
    transition: all 0.3s;
    position: relative;
  }

  .tabs button:hover {
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
    color: #667eea;
  }

  .tabs button.active {
    color: #667eea;
    border-bottom-color: #667eea;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  }

  .tabs button:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .all-done {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    text-align: center;
    padding: 40px;
  }

  .all-done h2 {
    font-size: 36px;
    margin-bottom: 16px;
    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    font-weight: 700;
  }

  .all-done p {
    color: #666;
    font-size: 16px;
    line-height: 1.6;
  }

  .no-image {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #999;
    font-size: 16px;
  }
</style>
