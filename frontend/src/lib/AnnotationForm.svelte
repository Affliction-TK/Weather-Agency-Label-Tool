<script>
  import { createEventDispatcher } from 'svelte';
  
  export let image;
  export let annotation = null;
  export let stations = [];

  const dispatch = createEventDispatcher();
  const API_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080/api' : '/api';
  const IMAGE_BASE = window.location.hostname === 'localhost' ? 'http://localhost:8080/images' : '/images';

  let formData = {
    category: annotation?.category || '大雾',
    severity: annotation?.severity || '轻度',
    observationTime: annotation?.observation_time ? new Date(annotation.observation_time).toISOString().slice(0, 16) : new Date().toISOString().slice(0, 16),
    location: annotation?.location || '',
    longitude: annotation?.longitude || '',
    latitude: annotation?.latitude || '',
    stationId: annotation?.station_id || ''
  };

  let saving = false;
  let suggestedStation = null;

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
        if (!formData.stationId) {
          formData.stationId = suggestedStation.id;
        }
      }
    } catch (error) {
      console.error('Failed to find nearest station:', error);
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
        station_id: parseInt(formData.stationId)
      };

      const response = await fetch(`${API_BASE}/annotations`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      });

      if (response.ok) {
        dispatch('saved');
        alert('标注保存成功！');
      } else {
        alert('保存失败，请重试');
      }
    } catch (error) {
      console.error('Failed to save annotation:', error);
      alert('保存失败：' + error.message);
    } finally {
      saving = false;
    }
  }
</script>

<div class="annotation-form">
  <div class="image-preview">
    <img src="{IMAGE_BASE}/{image.filename}" alt={image.filename} />
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
          <option value="积劳">积劳</option>
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

    <div class="form-group">
      <label for="station">监测点 *</label>
      <select id="station" bind:value={formData.stationId} required>
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
      <button type="submit" disabled={saving}>
        {#if saving}
          保存中...
        {:else}
          保存标注
        {/if}
      </button>
    </div>
  </form>
</div>

<style>
  .annotation-form {
    height: 100%;
    overflow-y: auto;
    padding: 20px;
  }

  .image-preview {
    margin-bottom: 20px;
    border: 1px solid #ddd;
    border-radius: 8px;
    overflow: hidden;
    background: #f9f9f9;
  }

  .image-preview img {
    width: 100%;
    max-height: 400px;
    object-fit: contain;
    background: white;
  }

  .image-info {
    padding: 10px;
    font-size: 14px;
    color: #666;
  }

  form {
    background: white;
    padding: 20px;
    border-radius: 8px;
    border: 1px solid #ddd;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 15px;
  }

  .form-group {
    margin-bottom: 20px;
  }

  label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
    color: #333;
    font-size: 14px;
  }

  input, select {
    width: 100%;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
    box-sizing: border-box;
  }

  input:focus, select:focus {
    outline: none;
    border-color: #1976d2;
    box-shadow: 0 0 0 2px rgba(25, 118, 210, 0.1);
  }

  .suggestion {
    margin-top: 8px;
    padding: 8px;
    background: #e3f2fd;
    border-radius: 4px;
    font-size: 12px;
    color: #1976d2;
  }

  .form-actions {
    margin-top: 30px;
    text-align: right;
  }

  button {
    padding: 12px 30px;
    background: #1976d2;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
  }

  button:hover:not(:disabled) {
    background: #1565c0;
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
