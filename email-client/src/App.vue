<template>
  <section class="app w-4/5 m-auto border p-2 mt-2 rounded-md">
    <section class="flex flex-col items-center">
      <EmailSearchBar @search="handleSearch" class="w-full" />
    </section>

    <!-- Modal to show email details -->
    <div class="fixed inset-0 z-10 flex items-center justify-center" :class="{ hidden: !showMail }">
      <div class="fixed inset-0 bg-black opacity-50" @click="closeModal"></div>
      <div
        class="relative bg-white rounded-lg shadow-lg p-6 overflow-y-auto max-h-[80vh] w-full sm:max-w-lg"
      >
        <h2 class="text-lg md:text-xl font-bold mb-4">{{ emailData.subject }}</h2>
        <p class="text-gray-800 mb-2"><strong>From:</strong> {{ emailData.from }}</p>
        <p class="text-gray-800 mb-2"><strong>To:</strong> {{ emailData.to }}</p>
        <p class="text-gray-800 mb-2"><strong>Date:</strong> {{ emailData.date }}</p>
        <p class="text-gray-800 mb-2"><strong>Content:</strong></p>
        <p class="text-gray-800 whitespace-pre-wrap">{{ emailData.content }}</p>
      </div>
    </div>

    <!-- Modal if the term is not found -->
    <div
      class="fixed inset-0 z-10 flex items-center justify-center"
      :class="{ hidden: errorMessage === '' }"
    >
      <div class="fixed inset-0 bg-black opacity-50" @click="clearError"></div>
      <div class="relative bg-white rounded-lg shadow-lg p-6 text-center w-full sm:max-w-md">
        <p>{{ errorMessage }}</p>
        <button
          class="bg-red-700 mt-4 text-white rounded-md p-2 opacity-90 hover:opacity-100"
          @click="clearError"
        >
          Close
        </button>
      </div>
    </div>

    <!-- Emails list -->
    <section class="w-full overflow-x-auto">
      <table class="w-full divide-y divide-gray-200 border">
        <thead class="bg-emerald-800 rounded-md border">
          <tr class="rounded-md p-2 border">
            <th
              class="px-4 py-3 text-left text-xs sm:text-sm md:text-base font-medium text-gray-200 uppercase tracking-wider"
            >
              Date
            </th>
            <th
              class="px-4 py-3 text-left text-xs sm:text-sm md:text-base font-medium text-gray-200 uppercase tracking-wider"
            >
              From
            </th>
            <th
              class="px-4 py-3 text-left text-xs sm:text-sm md:text-base font-medium text-gray-200 uppercase tracking-wider"
            >
              To
            </th>
            <th
              class="px-4 py-3 text-left text-xs sm:text-sm md:text-base font-medium text-gray-200 uppercase tracking-wider"
            >
              Subject
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="(email, index) in paginatedEmails" :key="index" @click="openEmailModal(email)">
            <td class="text-sm px-4 py-2 whitespace-nowrap">{{ email.date }}</td>
            <td class="text-sm px-4 py-2 whitespace-nowrap">{{ email.from }}</td>
            <td class="text-sm px-4 py-2 whitespace-wrap">
              {{ email.to === '' ? 'No recipient' : email.to }}
            </td>
            <td class="text-sm px-4 py-2 whitespace-nowrap">{{ email.subject }}</td>
          </tr>
        </tbody>
      </table>
    </section>
    <!-- Pagination Controls -->
    <section class="flex justify-end p-2">
      <div class="pagination-controls flex items-center py-4 justify-end rounded-md p-2 border">
        <button @click="prevPage" :disabled="currentPage <= 1">
          <img src="./assets/chevron-left.svg" alt="Anterior" />
        </button>
        <span>Page {{ currentPage }} of {{ totalPages }}</span>
        <button @click="nextPage" :disabled="currentPage >= totalPages">
          <img src="./assets/chevron-right.svg" alt="Siguiente" />
        </button>
        <!-- Add buttons for each page -->
        <template v-if="totalPages > 1">
          <button
            v-for="pageNumber in displayedPages"
            :key="pageNumber"
            @click="goToPage(pageNumber)"
            class="pagination-button flex items-center justify-end rounded-md p-2 border"
          >
            {{ pageNumber }}
          </button>
        </template>
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { IEmail, IEmailResponse } from './interfaces/IEmail'
import EmailSearchBar from './components/EmailSearchBar.vue'

const showMail = ref(false)
const emailData = ref<IEmailResponse>({
  from: '',
  to: '',
  date: '',
  subject: '',
  content: ''
})

const openEmailModal = (email: IEmail) => {
  showMail.value = true
  emailData.value = {
    from: email.from,
    to: email.to,
    date: email.date,
    subject: email.subject,
    content: email.content
  }
}

const resultsPerPage = 10
const currentPage = ref(1)
const matchingEmails = ref<IEmail[]>([])
const errorMessage = ref('')

const clearError = () => {
  errorMessage.value = ''
}

const handleSearch = async (searchTerm: string, sortField: string) => {
  if (searchTerm === '') {
    matchingEmails.value = []
    currentPage.value = 1
    errorMessage.value = ''
  } else {
    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/search`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ searchTerm, sortField })
      })

      if (!response.ok) throw new Error('Something went wrong')

      const { emails } = await response.json()

      if (emails.length === 0) {
        errorMessage.value = 'Term not found'
      } else {
        matchingEmails.value = emails
        errorMessage.value = ''
      }
    } catch (err) {
      console.error(err)
    }
  }
}

const closeModal = () => {
  showMail.value = false
}

onMounted(() => {
  document.addEventListener('keydown', handleEscapeKey)
})

const handleEscapeKey = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && showMail.value) {
    closeModal()
  }
}

const paginatedEmails = computed(() => {
  const start = (currentPage.value - 1) * resultsPerPage
  const end = start + resultsPerPage
  return matchingEmails.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(matchingEmails.value.length / resultsPerPage)
})

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}
function goToPage(page: number) {
  currentPage.value = page
}

const displayedPages = computed(() => {
  const numPagesToShow = 5
  const currentPageIndex = currentPage.value
  const lastPage = totalPages.value
  let startPage = Math.max(1, currentPageIndex - Math.floor(numPagesToShow / 2))
  let endPage = Math.min(lastPage, startPage + numPagesToShow - 1)

  if (lastPage - endPage < Math.floor(numPagesToShow / 2)) {
    startPage = Math.max(1, lastPage - numPagesToShow + 1)
  }

  return Array.from({ length: endPage - startPage + 1 }, (_, i) => startPage + i)
})
</script>

<!-- <style scoped>
/* Emails list style */
.table-wrapper {
  max-height: 600px;
}
</style> -->
