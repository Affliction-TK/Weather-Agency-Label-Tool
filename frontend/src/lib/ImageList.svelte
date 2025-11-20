<script>
  import { createEventDispatcher } from 'svelte';
  
  /** @type {{id: number, filename: string, annotated: boolean}[]} */
  export let images = [];
  /** @type {{id: number} | null} */
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

  function handleKeyActivate(event, callback) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      callback();
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
    <div 
      class="section-header" 
      on:click={toggleUnannotated}
      role="button"
      tabindex="0"
      on:keydown={(event) => handleKeyActivate(event, toggleUnannotated)}
    >
      <span class="toggle-icon" class:collapsed={!showUnannotated}>▼</span>
      <h3>未标注 <span class="count">{unannotatedImages.length}</span></h3>
    </div>
    {#if showUnannotated}
      <div class="section-content">
        {#if unannotatedImages.length === 0}
          <div class="empty-message-container">
            <span class="empty-icon">✨</span>
            <p class="empty-message">暂无未标注图片</p>
          </div>
        {:else}
          {#each unannotatedImages as image (image.id)}
            <div 
              class="image-item unannotated-item"
              class:active={currentImage && currentImage.id === image.id}
              on:click={() => selectImage(image)}
              role="button"
              tabindex="0"
              on:keydown={(event) => handleKeyActivate(event, () => selectImage(image))}
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
    <div 
      class="section-header" 
      on:click={toggleAnnotated}
      role="button"
      tabindex="0"
      on:keydown={(event) => handleKeyActivate(event, toggleAnnotated)}
    >
      <span class="toggle-icon" class:collapsed={!showAnnotated}>▼</span>
      <h3>已标注 <span class="count">{annotatedImages.length}</span></h3>
    </div>
    {#if showAnnotated}
      <div class="section-content">
        {#if annotatedImages.length === 0}
          <div class="empty-message-container">
            <span class="empty-icon">📭</span>
            <p class="empty-message">暂无已标注图片</p>
          </div>
        {:else}
          {#each annotatedImages as image (image.id)}
            <div 
              class="image-item annotated-item"
              class:active={currentImage && currentImage.id === image.id}
              on:click={() => selectImage(image)}
              role="button"
              tabindex="0"
              on:keydown={(event) => handleKeyActivate(event, () => selectImage(image))}
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
    gap: 16px;
  }

  .section {
    display: flex;
    flex-direction: column;
  }

  .section-header {
    display: flex;
    align-items: center;
    padding: 8px 4px;
    cursor: pointer;
    user-select: none;
    color: var(--text-secondary);
    transition: color 0.2s;
  }

  .section-header:hover {
    color: var(--text-primary);
  }

  .section-header h3 {
    margin: 0;
    font-size: 13px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    flex: 1;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .count {
    color: var(--text-secondary);
    font-weight: 400;
  }

  .toggle-icon {
    margin-right: 6px;
    font-size: 10px;
    transition: transform 0.2s;
    opacity: 0.7;
  }

  .toggle-icon.collapsed {
    transform: rotate(-90deg);
  }

  .section-content {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .empty-message-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 24px 16px;
    color: var(--text-secondary);
    background: rgba(0,0,0,0.02);
    border-radius: 8px;
    margin: 4px 0;
    border: 1px dashed rgba(0,0,0,0.05);
  }

  .empty-icon {
    font-size: 24px;
    margin-bottom: 8px;
    opacity: 0.6;
  }

  .empty-message {
    text-align: center;
    margin: 0;
    font-size: 12px;
    font-weight: 500;
  }

  .image-item {
    display: flex;
    align-items: center;
    padding: 8px;
    border-radius: var(--radius-s);
    cursor: pointer;
    transition: all 0.2s;
    background: transparent;
  }

  .image-item:hover {
    background: rgba(0, 0, 0, 0.04);
  }

  .image-item.active {
    background: var(--primary-color);
    color: white;
    box-shadow: var(--shadow-sm);
  }

  .thumbnail {
    width: 48px;
    height: 48px;
    flex-shrink: 0;
    overflow: hidden;
    border-radius: 6px;
    background: #E5E5EA;
    border: 1px solid rgba(0,0,0,0.1);
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
    justify-content: center;
    min-width: 0;
    gap: 4px;
  }

  .filename {
    font-size: 13px;
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    color: var(--text-primary);
  }

  .image-item.active .filename {
    color: white;
  }

  .status {
    display: flex;
    align-items: center;
  }

  .badge {
    font-size: 11px;
    font-weight: 500;
    color: var(--text-secondary);
  }

  .image-item.active .badge {
    color: rgba(255, 255, 255, 0.8);
  }

  /* Hide badges in list to keep it clean, or make them very subtle dots */
  .badge {
    display: none; 
  }
  
  /* Add a status dot instead */
  .info::after {
    content: '';
    display: block;
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: var(--warning-color); /* Default unannotated */
    margin-top: 4px;
  }

  .image-item.active .info::after {
    background-color: rgba(255, 255, 255, 0.8);
  }

  /* Status dots */
  .unannotated-item .info::after {
    background-color: var(--warning-color);
  }

  .annotated-item .info::after {
    background-color: var(--success-color);
  }

  /* We need to target the parent to change the dot color based on section, 
     but CSS doesn't support parent selectors easily based on children classes without :has.
     However, we know which section we are in.
  */
</style>
