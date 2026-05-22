<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <h2>Eliminar Biblioteca</h2>
      <p class="warning-text">
        ¿Estás seguro de que deseas eliminar la biblioteca "{{ library?.name }}"?
      </p>
      <p class="warning-text secondary">
        Esta acción eliminará todas las series y capítulos asociados. Los archivos de cómics no serán eliminados.
      </p>

      <form @submit.prevent="handleSubmit">
        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="$emit('close')">
            Cancelar
          </button>
          <button type="submit" class="btn btn-danger" :disabled="loading">
            {{ loading ? 'Eliminando...' : 'Eliminar' }}
          </button>
        </div>

        <p v-if="error" class="error-message">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  library: any
}>()

const emit = defineEmits(['close', 'deleted'])

const loading = ref(false)
const error = ref('')

async function handleSubmit() {
  loading.value = true
  error.value = ''

  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/v1/library/${props.library.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      const data = await response.json()
      throw new Error(data.error || 'Error al eliminar biblioteca')
    }

    emit('deleted')
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
  margin-bottom: 1rem;
}

.warning-text {
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.warning-text.secondary {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.btn-danger {
  background-color: var(--error);
  color: white;
  border: none;
}

.btn-danger:hover {
  opacity: 0.9;
}

.error-message {
  color: var(--error);
  font-size: 0.875rem;
  margin-top: 1rem;
}
</style>