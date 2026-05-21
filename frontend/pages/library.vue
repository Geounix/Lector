<template>
  <div class="library-page">
    <header class="library-header">
      <div class="header-left">
        <h1>Mi Biblioteca</h1>
      </div>
      <div class="header-right">
        <input
          v-model="searchQuery"
          type="search"
          class="input search-input"
          placeholder="Buscar comics..."
          @input="handleSearch"
        />
        <button class="btn btn-primary" @click="showAddLibrary = true">
          + Agregar Biblioteca
        </button>
      </div>
    </header>

    <div v-if="loading" class="loading">
      <p>Cargando biblioteca...</p>
    </div>

    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <button class="btn btn-secondary" @click="fetchLibrary">Reintentar</button>
    </div>

    <div v-else>
      <div v-for="library in libraries" :key="library.id" class="library-section">
        <h2 class="section-title">{{ library.name }}</h2>
        <div v-if="library.series && library.series.length > 0" class="grid grid-4">
          <div
            v-for="series in library.series"
            :key="series.id"
            class="card series-card"
            @click="openSeries(series)"
          >
            <div class="cover-container">
              <img
                v-if="series.cover"
                :src="series.cover"
                :alt="series.title"
                class="cover-image"
              />
              <div v-else class="cover-placeholder">
                <span>📚</span>
              </div>
            </div>
            <div class="series-info">
              <h3 class="series-title">{{ series.title }}</h3>
              <p class="series-meta">
                {{ series.chapterCount || 0 }} capítulos
                <span v-if="series.year"> • {{ series.year }}</span>
              </p>
            </div>
          </div>
        </div>
        <div v-else class="empty-section">
          <p>No hay comics en esta biblioteca</p>
          <button class="btn btn-secondary" @click="scanLibrary(library.id)">
            Escanear ahora
          </button>
        </div>
      </div>

      <div v-if="libraries.length === 0" class="empty-state">
        <h2>No tienes bibliotecas</h2>
        <p>Crea tu primera biblioteca para empezar a organizar tus comics</p>
        <button class="btn btn-primary" @click="showAddLibrary = true">
          Crear Biblioteca
        </button>
      </div>
    </div>

    <AddLibraryModal
      v-if="showAddLibrary"
      @close="showAddLibrary = false"
      @created="handleLibraryCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const libraries = ref<any[]>([])
const searchQuery = ref('')
const loading = ref(true)
const error = ref('')
const showAddLibrary = ref(false)

onMounted(() => {
  fetchLibrary()
})

async function fetchLibrary() {
  loading.value = true
  error.value = ''

  try {
    const token = localStorage.getItem('token')

    const [libResponse, seriesResponse] = await Promise.all([
      fetch('/api/v1/library', {
        headers: { 'Authorization': `Bearer ${token}` }
      }),
      fetch('/api/v1/series', {
        headers: { 'Authorization': `Bearer ${token}` }
      })
    ])

    if (!libResponse.ok || !seriesResponse.ok) {
      throw new Error('Error al cargar la biblioteca')
    }

    const librariesData = await libResponse.json()
    const seriesData = await seriesResponse.json()

    libraries.value = librariesData.map((lib: any) => ({
      ...lib,
      series: seriesData.filter((s: any) => s.library_id === lib.id)
    }))
  } catch (err: any) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

function handleSearch() {
}

function openSeries(series: any) {
  if (series.chapters && series.chapters.length > 0) {
    router.push(`/reader/${series.chapters[0].id}`)
  }
}

async function scanLibrary(libraryId: number) {
  try {
    const token = localStorage.getItem('token')
    await fetch('/api/v1/series/scan', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ libraryId })
    })
    await fetchLibrary()
  } catch (err) {
    console.error('Scan failed:', err)
  }
}

function handleLibraryCreated() {
  showAddLibrary.value = false
  fetchLibrary()
}
</script>

<style scoped>
.library-page {
  max-width: 1400px;
  margin: 0 auto;
}

.library-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.header-left h1 {
  font-size: 1.75rem;
}

.header-right {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.search-input {
  width: 300px;
}

.loading,
.error,
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: var(--text-secondary);
}

.section-title {
  font-size: 1.25rem;
  margin-bottom: 1rem;
  color: var(--text-secondary);
}

.series-card {
  cursor: pointer;
  transition: transform 0.2s, border-color 0.2s;
}

.series-card:hover {
  transform: translateY(-4px);
  border-color: var(--accent);
}

.cover-container {
  aspect-ratio: 2/3;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 0.75rem;
  background-color: var(--bg-tertiary);
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3rem;
}

.series-title {
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 0.25rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.series-meta {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.empty-section {
  text-align: center;
  padding: 3rem;
  background-color: var(--bg-secondary);
  border-radius: 12px;
  border: 1px dashed var(--border);
}

.empty-section p {
  color: var(--text-secondary);
  margin-bottom: 1rem;
}
</style>