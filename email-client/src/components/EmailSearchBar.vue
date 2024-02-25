<template>
  <section class="flex justify-end gap-4 p-2">
    <input
      type="text"
      class="text-black rounded-md p-2 w-10/12 border"
      v-model="searchText"
      placeholder="Enter your search term..."
      @keyup.enter="onSearch"
    />
    <ActionButton @click="onSearch">Search</ActionButton>
  </section>
</template>

<script setup lang="ts">
import { ref, defineEmits } from 'vue'
import ActionButton from './ActionButton.vue'

const searchText = ref('')
const emit = defineEmits(['search'])

const onSearch = () => {
  emit('search', searchText.value)
}
</script>
