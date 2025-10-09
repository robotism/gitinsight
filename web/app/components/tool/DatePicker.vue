<template>
    <div>
      <q-input
        dense
        flat
        v-model="displayDate"
        :prefix="prefix"
        :label="label"
        :placeholder="placeholder"
        @click="showDate = true"
        :clearable="clearable"
        @clear="clearDate"
      >
        <template #append>
          <q-icon name="event" class="cursor-pointer" @click="showDate = true" />
        </template>
      </q-input>
  
      <q-popup-proxy v-model="showDate" transition-show="scale" transition-hide="scale">
        <q-date
          v-model="innerValue"
          mask="YYYY-MM-DD"
          :default-year-month="defaultYearMonth"
          @update:model-value="onDateSelect"
        />
      </q-popup-proxy>
    </div>
  </template>
  
  <script setup>
  import { ref, watch, computed } from 'vue'
  
  /**
   * Props
   */
  const props = defineProps({
    /** v-model 绑定值 */
    modelValue: {
      type: String,
      default: ''
    },
    /** 显示的 label */
    label: {
      type: String,
      default: ''
    },
    /** 前缀文字 */
    prefix: {
      type: String,
      default: ''
    },
    /** 占位符（默认空） */
    placeholder: {
      type: String,
      default: ''
    },
    /** 默认显示的年月 */
    defaultYearMonth: {
      type: String,
      default: ''
    },
    /** 是否可清除 */
    clearable: {
      type: Boolean,
      default: true
    }
  })
  
  /**
   * Emits
   */
  const emit = defineEmits(['update:modelValue'])
  
  /**
   * State
   */
  const innerValue = ref(props.modelValue || '')
  const displayDate = ref(props.modelValue || '')
  const showDate = ref(false)
  
  /**
   * Watchers
   */
  watch(
    () => props.modelValue,
    (val) => {
      innerValue.value = val
      displayDate.value = val
    }
  )
  
  /**
   * Methods
   */
  function onDateSelect(val) {
    innerValue.value = val
    displayDate.value = val
    emit('update:modelValue', val)
    showDate.value = false
  }
  
  function clearDate() {
    innerValue.value = ''
    displayDate.value = ''
    emit('update:modelValue', '')
  }
  </script>
  