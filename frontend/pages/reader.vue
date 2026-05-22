<template>
  <div class="reader-page">
    <header class="reader-header">
      <button class="btn btn-secondary" @click="goBack">
        ← Volver
      </button>
      <div class="chapter-info">
        <h1>{{ chapterTitle }}</h1>
      </div>
      <div class="reader-controls">
        <button class="btn btn-secondary" @click="toggleFullscreen">
          {{ isFullscreen ? 'Salir' : 'Pantalla Completa' }}
        </button>
      </div>
    </header>

    <div class="reader-container" :class="{ fullscreen: isFullscreen }">
      <div v-if="loading" class="loading">Cargando...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else class="page-viewer">
        <img
          :src="currentPageSrc"
          :alt="`Página ${currentPage + 1}`"
          class="page-image"
          @load="onPageLoad"
        />
      </div>
    </div>

    <footer class="reader-footer">
      <button
        class="btn btn-secondary nav-btn"
        :disabled="currentPage === 0"
        @click="prevPage"
      >
        ← Anterior
      </button>

      <div class="page-indicator">
        <span>{{ currentPage + 1 }} / {{ totalPages }}</span>
        <input
          type="range"
          v-model.number="currentPage"
          :min="0"
          :max="totalPages - 1"
          class="page-slider"
        />
      </div>

      <button
        class="btn btn-secondary nav-btn"
        :disabled="currentPage >= totalPages - 1"
        @click="nextPage"
      >
        Siguiente →
      </button>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const chapterId = route.params.id as string

const loading = ref(true)
const error = ref('')
const chapterTitle = ref('')
const currentPage = ref(0)
const currentPageSrc = ref('')
const isFullscreen = ref(false)

onMounted(() => {
  fetchChapter()
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})

async function fetchChapter() {
  loading.value = true
  error.value = ''

  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/v1/reader/chapter/${chapterId}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      throw new Error('Error al cargar el capítulo')
    }

    const data = await response.json()
    chapterTitle.value = data.title || 'Capítulo'
    totalPages.value = data.page_count || 1
    currentPage.value = data.current_page || 0

    if (totalPages.value > 0) {
      loadPage(currentPage.value)
    }
  } catch (err: any) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const totalPages = ref(1)

async function loadPage(pageIndex: number) {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(
      `/api/v1/reader/chapter/${chapterId}/page/${pageIndex}`,
      {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      }
    )

    if (response.ok) {
      const blob = await response.blob()
      if (currentPageSrc.value) {
        URL.revokeObjectURL(currentPageSrc.value)
      }
      currentPageSrc.value = URL.createObjectURL(blob)
    }
  } catch (err) {
    console.error('Error loading page:', err)
  } finally {
    loading.value = false
  }
}

function prevPage() {
  if (currentPage.value > 0) {
    currentPage.value--
    loadPage(currentPage.value)
    saveProgress()
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value - 1) {
    currentPage.value++
    loadPage(currentPage.value)
    saveProgress()
  }
}

async function saveProgress() {
  try {
    const token = localStorage.getItem('token')
    await fetch('/api/v1/reader/progress', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        chapter_id: chapterId,
        page: currentPage.value,
        percentage: ((currentPage.value + 1) / totalPages.value) * 100
      })
    })
  } catch (err) {
    console.error('Error saving progress:', err)
  }
}

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}

function handleKeydown(e: KeyboardEvent) {
  switch (e.key) {
    case 'a':
    case 'ArrowLeft':
      prevPage()
      break
    case 'd':
    case 'ArrowRight':
      nextPage()
      break
    case 'f':
      toggleFullscreen()
      break
    case 'Escape':
      if (isFullscreen.value) {
        toggleFullscreen()
      } else {
        goBack()
      }
      break
  }
}

function goBack() {
  router.push('/library')
}
</script>

<style scoped>
.reader-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--bg-primary);
}

.reader-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
}

.chapter-info h1 {
  font-size: 1.25rem;
}

.reader-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.reader-container.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
}

.page-viewer {
  max-width: 100%;
  max-height: 100%;
}

.page-image {
  max-width: 100%;
  max-height: calc(100vh - 140px);
  object-fit: contain;
}

.reader-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: var(--bg-secondary);
  border-top: 1px solid var(--border);
}

.page-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.page-slider {
  width: 300px;
  cursor: pointer;
}

.nav-btn {
  min-width: 100px;
}

.loading,
.error {
  color: var(--text-secondary);
  padding: 2rem;
}
</style>