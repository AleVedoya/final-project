<template>
  <section class="app">
    <section>
      <EmailSearchBar @search="handleSearch" />
    </section>

    <!-- Modal -->
    <div class="fixed inset-0 z-10 flex items-center justify-center" :class="{ hidden: !showMail }">
      <div class="fixed inset-0 bg-black opacity-50" @click="closeModal"></div>
      <div class="relative bg-white rounded-lg shadow-lg p-6 overflow-y-auto max-h-[80vh]">
        <h2 class="text-xl font-bold mb-4">{{ emailData.subject }}</h2>
        <p class="text-gray-800 mb-2"><strong>From:</strong> {{ emailData.from }}</p>
        <p class="text-gray-800 mb-2"><strong>To:</strong> {{ emailData.to }}</p>
        <p class="text-gray-800 mb-2"><strong>Date:</strong> {{ emailData.date }}</p>
        <p class="text-gray-800 mb-2"><strong>Content:</strong></p>
        <p class="text-gray-800 whitespace-pre-wrap">{{ emailData.content }}</p>
      </div>
    </div>

    <!-- Emails list -->
    <section class="table-wrapper max-h-[600px] overflow-y-auto">
      <table class="w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th
              class="bg-emerald-800 px-6 py-3 text-left text-xs font-medium text-gray-200 uppercase tracking-wider"
            >
              From
            </th>
            <th
              class="bg-emerald-800 px-6 py-3 text-left text-xs font-medium text-gray-200 uppercase tracking-wider"
            >
              To
            </th>
            <th
              class="bg-emerald-800 px-6 py-3 text-left text-xs font-medium text-gray-200 uppercase tracking-wider"
            >
              Subject
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="(email, index) in matchingEmails" :key="index" @click="openEmailModal(email)">
            <td class="px-6 py-4 whitespace-nowrap">{{ email.from }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              {{ email.to === '' ? 'No recipient' : email.to }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">{{ email.subject }}</td>
          </tr>
        </tbody>
      </table>
    </section>

    <!-- Matching emails component -->
    <section class="flex flex-col items-center space-x-4 overflow-x-auto max-w-full">
      <button
        v-for="(email, index) in matchingEmails"
        :key="index"
        class="text-sm text-gray-900 bg-gray-200 px-4 py-2 rounded-md whitespace-nowrap"
        @click="openEmailModal(email)"
      >
        {{ email.from }}
        {{ email.to }}
        {{ email.subject }}
      </button>
    </section>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
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

const matchingEmails = ref<IEmail[]>([])

const handleSearch = async (searchTerm: string) => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/search`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ searchTerm })
    })

    if (!response.ok) throw new Error('Something went wrong')

    const { emails } = await response.json()

    matchingEmails.value = emails
  } catch (err) {
    console.error(err)
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
</script>

<style scoped>
.table-wrapper {
  max-height: 600px;
}
</style>

<!-- <26953671.1075846655760.JavaMail.evans@thyme> -->
