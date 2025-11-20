<script>
  import { createEventDispatcher } from 'svelte';
  import { toasts } from './toastStore.js';
  import ConfirmModal from './ConfirmModal.svelte';
  
  export let image;
  export let annotation = null;
  export let stations = [];

  const dispatch = createEventDispatcher();
  const API_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080/api' : '/api';
  const IMAGE_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080/images' : '/images';

  let formData = {
    category: '大雾',
    severity: '轻度',
    observationTime: new Date().toISOString().slice(0, 16),
    location: '',
    longitude: '',
    latitude: '',
    stationId: ''
  };

  let saving = false;
  let deleting = false;
  let showDeleteConfirm = false;
  let suggestedStation = null;
  let allowAutoStationSelection = true; // Keeps auto-selection active until user manually overrides

  // Reset form when image changes
  $: if (image) {
    resetForm();
  }

  function resetForm() {
    if (annotation) {
      // Load existing annotation
      formData = {
        category: annotation.category || '大雾',
        severity: annotation.severity || '轻度',
        observationTime: annotation.observation_time ? new Date(annotation.observation_time).toISOString().slice(0, 16) : new Date().toISOString().slice(0, 16),
        location: annotation.location || '',
        longitude: annotation.longitude || '',
        latitude: annotation.latitude || '',
        stationId: annotation.station_id || ''
      };
    } else {
      // Reset to defaults for new annotation
      // 使用OCR识别的结果预填充表单（如果有的话）
      let defaultTime = new Date().toISOString().slice(0, 16);
      let defaultLocation = '';
      
      // 如果图片有OCR识别的时间，使用OCR时间
      if (image.ocr_time) {
        try {
          // OCR时间格式是 "YYYY-MM-DD HH:MM:SS"，需要转换为表单需要的 "YYYY-MM-DDTHH:MM"
          const ocrDate = new Date(image.ocr_time);
          if (!isNaN(ocrDate.getTime())) {
            defaultTime = ocrDate.toISOString().slice(0, 16);
          }
        } catch (e) {
          console.log('Failed to parse OCR time:', e);
        }
      }
      
      // 如果图片有OCR识别的地点，使用OCR地点
      if (image.ocr_location) {
        defaultLocation = image.ocr_location;
      }
      
      formData = {
        category: '大雾',
        severity: '轻度',
        observationTime: defaultTime,
        location: defaultLocation,
        longitude: '',
        latitude: '',
        stationId: ''
      };
    }
    suggestedStation = null;
    allowAutoStationSelection = !formData.stationId; // Existing annotations keep their station unless user clears it
  }

  // Watch for coordinate changes to suggest nearest station
  $: if (formData.longitude && formData.latitude) {
    findNearestStation();
  }

  async function findNearestStation() {
    try {
      const lon = parseFloat(formData.longitude);
      const lat = parseFloat(formData.latitude);
      
      if (isNaN(lon) || isNaN(lat)) return;
      
      const response = await fetch(`${API_BASE}/stations/nearest?longitude=${lon}&latitude=${lat}`);
      if (response.ok) {
        suggestedStation = await response.json();
        if (allowAutoStationSelection) {
          formData.stationId = suggestedStation.id;
        }
      }
    } catch (error) {
      console.error('Failed to find nearest station:', error);
    }
  }

  function handleStationChange(event) {
    formData.stationId = event.target.value;
    allowAutoStationSelection = event.target.value === '';
  }

  let fetchingCoordinates = false;

  async function fetchCoordinates() {
    if (!formData.location || !formData.location.trim()) {
      toasts.error('请先输入地点');
      return;
    }

    fetchingCoordinates = true;
    try {
      const response = await fetch(`${API_BASE}/geocode`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ address: formData.location })
      });

      if (response.ok) {
        const data = await response.json();
        formData.longitude = data.longitude.toString();
        formData.latitude = data.latitude.toString();
        toasts.success('经纬度获取成功！');
      } else {
        const errorText = await response.text();
        toasts.error('获取经纬度失败：' + (errorText || '请检查地址是否正确'));
      }
    } catch (error) {
      console.error('Failed to fetch coordinates:', error);
      toasts.error('获取经纬度失败：' + error.message);
    } finally {
      fetchingCoordinates = false;
    }
  }

  async function handleSubmit() {
    saving = true;
    try {
      const payload = {
        image_id: image.id,
        category: formData.category,
        severity: formData.severity,
        observation_time: new Date(formData.observationTime).toISOString(),
        location: formData.location,
        longitude: parseFloat(formData.longitude),
        latitude: parseFloat(formData.latitude),
        station_id: formData.stationId
      };

      const response = await fetch(`${API_BASE}/annotations`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      });

      if (response.ok) {
        toasts.success('标注保存成功！');
        dispatch('saved');
      } else {
        const errorText = await response.text();
        toasts.error('保存失败：' + (errorText || '请重试'));
      }
    } catch (error) {
      console.error('Failed to save annotation:', error);
      toasts.error('保存失败：' + error.message);
    } finally {
      saving = false;
    }
  }

  function handleDelete() {
    if (!annotation || deleting) {
      return;
    }
    showDeleteConfirm = true;
  }

  function requestImageDelete() {
    if (annotation || !image) {
      return;
    }
    dispatch('requestImageDelete', image);
  }

  async function executeDelete() {
    showDeleteConfirm = false;
    deleting = true;
    try {
      const response = await fetch(`${API_BASE}/annotations/${annotation.id}`, {
        method: 'DELETE'
      });

      if (response.ok) {
        toasts.success('标注已删除');
        dispatch('deleted');
      } else {
        const errorText = await response.text();
        toasts.error('删除失败：' + (errorText || '请重试'));
      }
    } catch (error) {
      console.error('Failed to delete annotation:', error);
      toasts.error('删除失败：' + error.message);
    } finally {
      deleting = false;
    }
  }
</script>

<div class="annotation-form">
  <div class="image-container">
    <div class="image-preview">
      <img src="{IMAGE_BASE}/{image.filename}" alt={image.filename} />
    </div>
    <div class="image-info">
      <strong>文件名：</strong>{image.filename}
    </div>
  </div>

  <form on:submit|preventDefault={handleSubmit}>
    <div class="form-row">
      <div class="form-group">
        <label for="category">天气类型 *</label>
        <select id="category" bind:value={formData.category} required>
          <option value="大雾">大雾</option>
          <option value="结冰">结冰</option>
          <option value="积涝">积涝</option>
        </select>
      </div>

      <div class="form-group">
        <label for="severity">严重等级 *</label>
        <select id="severity" bind:value={formData.severity} required>
          <option value="无">无</option>
          <option value="轻度">轻度</option>
          <option value="中度">中度</option>
          <option value="重度">重度</option>
        </select>
      </div>
    </div>

    <div class="form-group">
      <label for="observationTime">观测时间 *</label>
      <input 
        type="datetime-local" 
        id="observationTime" 
        bind:value={formData.observationTime} 
        required 
      />
    </div>

    <div class="form-group">
      <label for="location">地点 *</label>
      <input 
        type="text" 
        id="location" 
        bind:value={formData.location} 
        placeholder="输入地点名称"
        required 
      />
    </div>

    <div class="form-row">
      <div class="form-group">
        <label for="longitude">经度 *</label>
        <input 
          type="number" 
          id="longitude" 
          bind:value={formData.longitude} 
          step="0.000001"
          placeholder="例如：116.407396"
          required 
        />
      </div>

      <div class="form-group">
        <label for="latitude">维度 *</label>
        <input 
          type="number" 
          id="latitude" 
          bind:value={formData.latitude} 
          step="0.000001"
          placeholder="例如：39.904211"
          required 
        />
      </div>
    </div>

    <div class="geocode-action">
      <button 
        type="button" 
        class="geocode-btn"
        on:click={fetchCoordinates}
        disabled={fetchingCoordinates || !formData.location}
      >
        {#if fetchingCoordinates}
          获取中...
        {:else}
          📍 获取经纬度
        {/if}
      </button>
      <span class="geocode-hint">根据地点自动获取经纬度坐标</span>
    </div>

    <div class="form-group">
      <label for="station">监测点 *</label>
      <select id="station" bind:value={formData.stationId} on:change={handleStationChange} required>
        <option value="">请选择监测点</option>
        {#each stations as station (station.id)}
          <option value={station.id}>
            {station.name} ({station.longitude}, {station.latitude})
          </option>
        {/each}
      </select>
      {#if suggestedStation}
        <div class="suggestion">
          推荐的最近站点：<strong>{suggestedStation.name}</strong>
        </div>
      {/if}
    </div>

    <div class="form-actions">
      {#if annotation}
        <button type="button" class="danger-btn" on:click={handleDelete} disabled={saving || deleting}>
          {#if deleting}
            删除中...
          {:else}
            删除标注
          {/if}
        </button>
      {:else}
        <button 
          type="button" 
          class="danger-btn"
          on:click={requestImageDelete}
          disabled={saving}
        >
          删除图片
        </button>
      {/if}
      <button type="submit" disabled={saving || deleting}>
        {#if saving}
          保存中...
        {:else}
          保存标注
        {/if}
      </button>
    </div>
  </form>
</div>

{#if showDeleteConfirm}
  <ConfirmModal
    title="删除确认"
    message="确认删除该标注记录？此操作无法撤销。"
    confirmText="删除"
    type="danger"
    on:confirm={executeDelete}
    on:cancel={() => showDeleteConfirm = false}
  />
{/if}

<style>
  .annotation-form {
    height: 100%;
    overflow-y: auto;
    padding: 32px;
    background: white;
    max-width: 900px;
    margin: 0 auto;
  }

  .image-container {
    margin-bottom: 32px;
    background: white;
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0,0,0,0.08);
    overflow: hidden;
    border: 1px solid #eee;
  }

  .image-preview {
    background: #f8f9fa;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 200px;
    padding: 20px;
  }

  .image-preview img {
    max-width: 100%;
    max-height: 500px;
    object-fit: contain;
    display: block;
    border-radius: 4px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  }

  .image-info {
    padding: 16px 20px;
    font-size: 14px;
    color: #333;
    background: white;
    border-top: 1px solid #eee;
    font-family: monospace;
  }

  form {
    background: transparent;
    padding: 0;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;
  }

  .form-group {
    margin-bottom: 24px;
  }

  label {
    display: block;
    margin-bottom: 8px;
    font-weight: 600;
    color: #555;
    font-size: 13px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  input, select {
    width: 100%;
    padding: 12px 16px;
    border: 1px solid #e1e4e8;
    border-radius: 8px;
    font-size: 15px;
    box-sizing: border-box;
    transition: all 0.2s;
    background: #fcfcfd;
    color: #333;
    appearance: none;
  }

  /* Custom select arrow */
  select {
    background-image: url("data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%23007AFF%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E");
    background-repeat: no-repeat;
    background-position: right 16px top 50%;
    background-size: 10px auto;
    padding-right: 40px;
  }

  input:focus, select:focus {
    outline: none;
    background: white;
    border-color: #007aff;
    box-shadow: 0 0 0 4px rgba(0, 122, 255, 0.1);
  }

  .suggestion {
    margin-top: 8px;
    padding: 10px 12px;
    background: rgba(52, 199, 89, 0.1);
    border-radius: 8px;
    font-size: 13px;
    color: #2e7d32;
    font-weight: 500;
    display: flex;
    align-items: center;
  }
  
  .suggestion::before {
    content: '📍';
    margin-right: 6px;
  }

    .form-actions {
      margin-top: 40px;
      display: flex;
      justify-content: flex-end;
      gap: 12px;
      flex-wrap: wrap;
    }

  button {
    padding: 12px 32px;
    background: #007aff;
    color: white;
    border: none;
    border-radius: 24px; /* Pill shape */
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 4px 12px rgba(0, 122, 255, 0.25);
  }

  .danger-btn {
    background: #ff3b30;
    box-shadow: 0 4px 12px rgba(255, 59, 48, 0.25);
  }

  .danger-btn:hover:not(:disabled) {
    background: #d70015;
    box-shadow: 0 6px 16px rgba(255, 59, 48, 0.35);
  }

  button:hover:not(:disabled) {
    background: #0062cc;
    transform: translateY(-1px);
    box-shadow: 0 6px 16px rgba(0, 122, 255, 0.35);
  }

  button:active:not(:disabled) {
    transform: translateY(0);
  }

  button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    box-shadow: none;
    background: #ccc;
  }

  .geocode-action {
    margin-bottom: 24px;
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .geocode-btn {
    padding: 10px 24px;
    background: #34c759;
    color: white;
    border: none;
    border-radius: 20px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 4px 12px rgba(52, 199, 89, 0.25);
    white-space: nowrap;
  }

  .geocode-btn:hover:not(:disabled) {
    background: #2ea64a;
    transform: translateY(-1px);
    box-shadow: 0 6px 16px rgba(52, 199, 89, 0.35);
  }

  .geocode-btn:active:not(:disabled) {
    transform: translateY(0);
  }

  .geocode-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    box-shadow: none;
    background: #ccc;
  }

  .geocode-hint {
    font-size: 13px;
    color: #666;
    font-style: italic;
  }
</style>
