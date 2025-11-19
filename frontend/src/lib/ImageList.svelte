<script>
  import { createEventDispatcher } from 'svelte';
  
  export let images = [];
  export let currentImage = null;

  const dispatch = createEventDispatcher();
  const IMAGE_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080/images' : '/images';

  // State for collapse/expand
  let showUnannotated = true;
  let showAnnotated = false;

  // Load collapse state from localStorage
  if (typeof localStorage !== 'undefined') {
    const savedState = localStorage.getItem('imageListCollapseState');
    if (savedState) {
      try {
        const state = JSON.parse(savedState);
        showUnannotated = state.showUnannotated ?? true;
        showAnnotated = state.showAnnotated ?? false;
      } catch (e) {
        // Ignore parse errors
      }
    }
  }

  // Save collapse state to localStorage
  function saveCollapseState() {
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('imageListCollapseState', JSON.stringify({
        showUnannotated,
        showAnnotated
      }));
    }
  }

  function toggleUnannotated() {
    showUnannotated = !showUnannotated;
    saveCollapseState();
  }

  function toggleAnnotated() {
    showAnnotated = !showAnnotated;
    saveCollapseState();
  }

  function selectImage(image) {
    dispatch('select', image);
  }

  // Filter images
  $: unannotatedImages = images.filter(img => !img.annotated);
  $: annotatedImages = images.filter(img => img.annotated);
</script>

<div class="image-list">
  <!-- Unannotated section -->
  <div class="section">
    <div class="section-header" on:click={toggleUnannotated}>
      <span class="toggle-icon" class:collapsed={!showUnannotated}>▼</span>
      <h3>未标注 <span class="count">{unannotatedImages.length}</span></h3>
    </div>
    {#if showUnannotated}
      <div class="section-content">
        {#if unannotatedImages.length === 0}
          <p class="empty-message">暂无未标注图片</p>
        {:else}
          {#each unannotatedImages as image (image.id)}
            <div 
              class="image-item"
              class:active={currentImage && currentImage.id === image.id}
              on:click={() => selectImage(image)}
            >
              <div class="thumbnail">
                <img src="{IMAGE_BASE}/{image.filename}" alt={image.filename} />
              </div>
              <div class="info">
                <div class="filename" title={image.filename}>{image.filename}</div>
                <div class="status">
                  <span class="badge unannotated">未标注</span>
                </div>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    {/if}
  </div>

  <!-- Annotated section -->
  <div class="section">
    <div class="section-header" on:click={toggleAnnotated}>
      <span class="toggle-icon" class:collapsed={!showAnnotated}>▼</span>
      <h3>已标注 <span class="count">{annotatedImages.length}</span></h3>
    </div>
    {#if showAnnotated}
      <div class="section-content">
        {#if annotatedImages.length === 0}
          <p class="empty-message">暂无已标注图片</p>
        {:else}
          {#each annotatedImages as image (image.id)}
            <div 
              class="image-item"
              class:active={currentImage && currentImage.id === image.id}
              on:click={() => selectImage(image)}
            >
              <div class="thumbnail">
                <img src="{IMAGE_BASE}/{image.filename}" alt={image.filename} />
              </div>
              <div class="info">
                <div class="filename" title={image.filename}>{image.filename}</div>
                <div class="status">
                  <span class="badge annotated">已标注</span>
                </div>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    {/if}
  </div>
</div>

<style>
  .image-list {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .section {
    background: white;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid #e0e0e0;
  }

  .section-header {
    display: flex;
    align-items: center;
    padding: 12px 15px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    cursor: pointer;
    user-select: none;
    transition: background 0.2s;
  }

  .section-header:hover {
    background: linear-gradient(135deg, #5568d3 0%, #63408a 100%);
  }

  .section-header h3 {
    margin: 0;
    font-size: 14px;
    font-weight: 600;
    color: white;
    flex: 1;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .count {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 24px;
    height: 24px;
    padding: 0 8px;
    background: rgba(255, 255, 255, 0.25);
    border-radius: 12px;
    font-size: 12px;
    font-weight: 600;
  }

  .toggle-icon {
    margin-right: 8px;
    font-size: 12px;
    transition: transform 0.2s;
    color: white;
  }

  .toggle-icon.collapsed {
    transform: rotate(-90deg);
  }

  .section-content {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 10px;
  }

  .empty-message {
    text-align: center;
    color: #999;
    padding: 20px;
    margin: 0;
    font-size: 13px;
  }

  .image-item {
    display: flex;
    background: white;
    border: 2px solid transparent;
    border-radius: 8px;
    padding: 10px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .image-item:hover {
    border-color: #667eea;
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.2);
    transform: translateY(-1px);
  }

  .image-item.active {
    border-color: #667eea;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  }

  .thumbnail {
    width: 80px;
    height: 60px;
    flex-shrink: 0;
    overflow: hidden;
    border-radius: 6px;
    background: #f0f0f0;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }

  .thumbnail img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .info {
    flex: 1;
    margin-left: 12px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    min-width: 0;
  }

  .filename {
    font-size: 12px;
    color: #333;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    font-weight: 500;
  }

  .status {
    display: flex;
    align-items: center;
  }

  .badge {
    padding: 3px 10px;
    border-radius: 12px;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .badge.annotated {
    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
    color: white;
  }

  .badge.unannotated {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
    color: white;
  }
</style>
