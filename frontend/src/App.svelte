<script>
  import { onMount } from 'svelte';
  import ImageList from './lib/ImageList.svelte';
  import AnnotationForm from './lib/AnnotationForm.svelte';
  import UploadTab from './lib/UploadTab.svelte';

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
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  }

  main {
    width: 100%;
    height: 100vh;
    overflow: hidden;
  }

  .container {
    display: flex;
    height: 100%;
  }

  .sidebar {
    width: 300px;
    background: #f5f5f5;
    border-right: 1px solid #ddd;
    overflow-y: auto;
    padding: 20px;
  }

  .sidebar h2 {
    margin-top: 0;
    font-size: 18px;
    color: #333;
  }

  .content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .tabs {
    display: flex;
    border-bottom: 1px solid #ddd;
    background: white;
  }

  .tabs button {
    padding: 15px 30px;
    border: none;
    background: none;
    cursor: pointer;
    font-size: 16px;
    color: #666;
    border-bottom: 3px solid transparent;
    transition: all 0.2s;
  }

  .tabs button:hover {
    background: #f5f5f5;
  }

  .tabs button.active {
    color: #1976d2;
    border-bottom-color: #1976d2;
  }

  .tabs button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .all-done {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    text-align: center;
    color: #4caf50;
  }

  .all-done h2 {
    font-size: 32px;
    margin-bottom: 10px;
  }

  .all-done p {
    color: #666;
    font-size: 16px;
  }

  .no-image {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #999;
  }
</style>
