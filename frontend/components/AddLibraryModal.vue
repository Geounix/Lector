<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <h2>Agregar Biblioteca</h2>
      
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="name">Nombre</label>
          <input
            id="name"
            v-model="name"
            type="text"
            class="input"
            placeholder="Mi Biblioteca"
            required
          />
        </div>

        <div class="form-group">
          <label for="type">Tipo</label>
          <select id="type" v-model="type" class="input">
            <option value="comics">Cómics</option>
            <option value="manga">Manga</option>
            <option value="webtoon">Webtoon</option>
          </select>
        </div>

        <div class="form-group">
          <label for="path">Ruta</label>
          <input
            id="path"
            v-model="path"
            type="text"
            class="input"
            placeholder="/path/to/comics"
            required
          />
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="$emit('close')">
            Cancelar
          </button>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? 'Creando...' : 'Crear' }}
          </button>
        </div>

        <p v-if="error" class="error-message">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits(['close', 'created'])

const name = ref('')
const type = ref('comics')
const path = ref('')
const loading = ref(false)
const error = ref('')

async function handleSubmit() {
  loading.value = true
  error.value = ''

  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/v1/library', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: name.value,
        type: type.value,
        path: path.value
      })
    })

    if (!response.ok) {
      const data = await response.json()
      throw new Error(data.error || 'Error al crear biblioteca')
    }

    emit('created')
  } catch (err: any) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.75);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 2rem;
  width: 100%;
  max-width: 450px;
}

.modal-content h2 {
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.error-message {
  color: var(--error);
  font-size: 0.875rem;
  margin-top: 1rem;
}
</style>