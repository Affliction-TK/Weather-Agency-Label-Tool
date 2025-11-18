<script>
  import { createEventDispatcher } from 'svelte';
  
  export let images = [];
  export let currentImage = null;

  const dispatch = createEventDispatcher();
  const IMAGE_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080/images' : '/images';

  function selectImage(image) {
    dispatch('select', image);
  }
</script>

<div class="image-list">
  {#each images as image (image.id)}
    <div 
      class="image-item"
      class:active={currentImage && currentImage.id === image.id}
      class:annotated={image.annotated}
      on:click={() => selectImage(image)}
    >
      <div class="thumbnail">
        <img src="{IMAGE_BASE}/{image.filename}" alt={image.filename} />
      </div>
      <div class="info">
        <div class="filename" title={image.filename}>{image.filename}</div>
        <div class="status">
          {#if image.annotated}
            <span class="badge annotated">已标注</span>
          {:else}
            <span class="badge unannotated">未标注</span>
          {/if}
        </div>
      </div>
    </div>
  {/each}
</div>

<style>
  .image-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
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
    border-color: #1976d2;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }

  .image-item.active {
    border-color: #1976d2;
    background: #e3f2fd;
  }

  .thumbnail {
    width: 80px;
    height: 60px;
    flex-shrink: 0;
    overflow: hidden;
    border-radius: 4px;
    background: #f0f0f0;
  }

  .thumbnail img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .info {
    flex: 1;
    margin-left: 10px;
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
  }

  .status {
    display: flex;
    align-items: center;
  }

  .badge {
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 11px;
    font-weight: 500;
  }

  .badge.annotated {
    background: #4caf50;
    color: white;
  }

  .badge.unannotated {
    background: #ff9800;
    color: white;
  }
</style>
