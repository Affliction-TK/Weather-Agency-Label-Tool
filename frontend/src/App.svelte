<script>
  import { onMount } from 'svelte';
  import ImageList from './lib/ImageList.svelte';
  import AnnotationForm from './lib/AnnotationForm.svelte';
  import UploadTab from './lib/UploadTab.svelte';
  import Toast from './lib/Toast.svelte';
  import ConfirmModal from './lib/ConfirmModal.svelte';
  import { toasts } from './lib/toastStore.js';

  let images = [];
  let stations = [];
  let currentImage = null;
  let currentAnnotation = null;
  let activeTab = 'annotate'; // 'annotate' or 'upload'
  let loading = true;
  let searchQuery = '';
  let showImageDeleteConfirm = false;
  let imagePendingDelete = null;
  let deletingImage = false;

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

  async function handleAnnotationDeleted() {
    const previousImageId = currentImage?.id;
    await loadImages();

    if (!previousImageId) {
      currentImage = null;
      currentAnnotation = null;
      return;
    }

    const refreshedImage = images.find(img => img.id === previousImageId);
    if (refreshedImage) {
      await selectImage(refreshedImage);
    } else {
      currentImage = null;
      currentAnnotation = null;
    }
  }

  $: allAnnotated = images.length > 0 && images.every(img => img.annotated);
  $: normalizedSearch = searchQuery.trim().toLowerCase();
  $: filteredImages = normalizedSearch
    ? images.filter(img => (img.filename || '').toLowerCase().includes(normalizedSearch))
    : images;

  function handleImageDeleteRequest(image) {
    if (!image) {
      return;
    }
    imagePendingDelete = image;
    showImageDeleteConfirm = true;
  }

  function closeImageDeleteModal() {
    showImageDeleteConfirm = false;
    imagePendingDelete = null;
  }

  async function confirmImageDelete() {
    if (!imagePendingDelete || deletingImage) {
      return;
    }

    deletingImage = true;
    try {
      const response = await fetch(`${API_BASE}/images/${imagePendingDelete.id}`, {
        method: 'DELETE'
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || '请重试');
      }

      const deletedId = imagePendingDelete.id;
      if (currentImage && currentImage.id === deletedId) {
        currentImage = null;
        currentAnnotation = null;
      }

      toasts.success('图片已删除');
      await loadImages();
    } catch (error) {
      console.error('Failed to delete image:', error);
      toasts.error('删除失败：' + (error.message || '请重试'));
    } finally {
      deletingImage = false;
      closeImageDeleteModal();
    }
  }
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
      <div class="app-header">
        <div class="logo">🌤️</div>
        <h1>无锡气象局<br>图像标注工具</h1>
      </div>
      <div class="list-header">
        <h2>图片列表</h2>
      </div>
      <div class="search-box">
        <input
          type="text"
          placeholder="搜索文件名"
          bind:value={searchQuery}
        />
        {#if searchQuery}
          <button
            type="button"
            class="clear-btn"
            on:click={() => searchQuery = ''}
            aria-label="清除搜索"
          >&times;</button>
        {/if}
      </div>
      {#if loading}
        <div class="loading-state">
          <div class="spinner"></div>
          <p>加载中...</p>
        </div>
      {:else if images.length === 0}
        <div class="empty-list">
          <p>暂无图片</p>
        </div>
      {:else if filteredImages.length === 0}
        <div class="empty-list">
          <p>未找到匹配的图片</p>
          <button class="action-btn" on:click={() => searchQuery = ''}>清除搜索</button>
        </div>
      {:else}
        <ImageList 
          images={filteredImages} 
          {currentImage} 
          on:select={(e) => selectImage(e.detail)}
        />
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
            <div class="status-icon">🎉</div>
            <h2>所有图片已标注完成</h2>
            <p>太棒了！您已经完成了所有任务。</p>
            <button class="action-btn" on:click={() => handleTabChange('upload')}>上传新图片</button>
          </div>
        {:else if currentImage}
          <AnnotationForm 
            image={currentImage} 
            annotation={currentAnnotation}
            {stations}
            on:saved={handleAnnotationSaved}
            on:deleted={handleAnnotationDeleted}
            on:requestImageDelete={(e) => handleImageDeleteRequest(e.detail)}
          />
        {:else}
          <div class="no-image">
            <div class="empty-state">
              <div class="empty-icon">👈</div>
              <h3>准备开始</h3>
              <p>请从左侧列表选择一张图片或上传图片进行标注</p>
            </div>
          </div>
        {/if}
      {:else if activeTab === 'upload'}
        <UploadTab on:uploaded={handleImageUploaded} />
      {/if}
    </div>
  </div>
</main>

{#if showImageDeleteConfirm && imagePendingDelete}
  <ConfirmModal
    title="删除图片"
    message={`确认删除图片 ${imagePendingDelete.filename} 吗？此操作无法撤销。`}
    confirmText={deletingImage ? '删除中...' : '删除'}
    type="danger"
    on:confirm={confirmImageDelete}
    on:cancel={closeImageDeleteModal}
  />
{/if}

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    background-color: #f0f2f5;
    color: #333;
  }

  main {
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(120deg, #e0c3fc 0%, #8ec5fc 100%); /* 清新渐变背景 */
    padding: 20px;
    box-sizing: border-box;
  }

  .container {
    display: flex;
    width: 100%;
    max-width: 1400px;
    height: 90vh;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border-radius: 16px;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    overflow: hidden;
    border: 1px solid rgba(255, 255, 255, 0.6);
  }

  .sidebar {
    width: 300px;
    background: rgba(250, 250, 252, 0.8);
    border-right: 1px solid rgba(0, 0, 0, 0.06);
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
  }

  .app-header {
    padding: 24px 20px;
    display: flex;
    align-items: center;
    gap: 12px;
    border-bottom: 1px solid rgba(0,0,0,0.04);
  }

  .logo {
    font-size: 32px;
    filter: drop-shadow(0 2px 4px rgba(0,0,0,0.1));
  }

  .app-header h1 {
    margin: 0;
    font-size: 16px;
    line-height: 1.4;
    font-weight: 700;
    color: #1a1a1a;
  }

  .list-header {
    padding: 16px 20px 8px;
  }

  .list-header h2 {
    margin: 0;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: #888;
    font-weight: 600;
  }

  .content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background: white;
    position: relative;
  }

  .tabs {
    display: flex;
    justify-content: center;
    padding: 12px;
    background: white;
    border-bottom: 1px solid #eee;
    z-index: 10;
  }

  .tabs button {
    padding: 8px 24px;
    border: none;
    background: transparent;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    color: #666;
    transition: all 0.2s ease;
    border-radius: 8px;
    margin: 0 4px;
  }

  .tabs button:hover {
    color: #333;
    background-color: #f5f5f5;
  }

  .tabs button.active {
    color: #007aff;
    background-color: rgba(0, 122, 255, 0.08);
    font-weight: 600;
  }

  .tabs button:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  /* Empty States */
  .all-done, .no-image {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    text-align: center;
    padding: 40px;
    color: #666;
  }

  .status-icon, .empty-icon {
    font-size: 64px;
    margin-bottom: 24px;
    opacity: 0.8;
  }

  .empty-state h3 {
    margin: 0 0 8px;
    font-size: 20px;
    color: #333;
  }

  .all-done h2 {
    font-size: 24px;
    margin-bottom: 12px;
    color: #333;
  }

  .action-btn {
    margin-top: 24px;
    padding: 10px 24px;
    background: #007aff;
    color: white;
    border: none;
    border-radius: 20px;
    font-size: 15px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 4px 12px rgba(0, 122, 255, 0.2);
  }

  .action-btn:hover {
    background: #0062cc;
    transform: translateY(-1px);
    box-shadow: 0 6px 16px rgba(0, 122, 255, 0.3);
  }

  .loading-state, .empty-list {
    padding: 40px 20px;
    text-align: center;
    color: #999;
  }

  .search-box {
    position: relative;
    padding: 0 20px 16px;
  }

  .search-box input {
    width: 100%;
    padding: 10px 36px 10px 12px;
    border-radius: 10px;
    border: 1px solid rgba(0, 0, 0, 0.1);
    font-size: 13px;
    background: rgba(255, 255, 255, 0.8);
    box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.04);
  }

  .search-box input:focus {
    outline: none;
    border-color: #007aff;
    box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.1);
  }

  .clear-btn {
    position: absolute;
    right: 28px;
    top: 50%;
    transform: translateY(-50%);
    border: none;
    background: transparent;
    font-size: 18px;
    line-height: 1;
    color: #999;
    cursor: pointer;
    padding: 0;
  }

  .clear-btn:hover {
    color: #333;
  }

  .spinner {
    width: 24px;
    height: 24px;
    border: 2px solid #eee;
    border-top-color: #007aff;
    border-radius: 50%;
    margin: 0 auto 16px;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
