<script setup lang="ts">
import { useAppStore } from '../../stores/app'
import { useSettingsStore } from '../../stores/settings'

const app = useAppStore()
const settings = useSettingsStore()

const emit = defineEmits<{
  close: []
}>()
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" @click.self="emit('close')">
    <div class="bg-[var(--color-surface)] rounded-xl shadow-xl border border-[var(--color-border)] w-full max-w-md p-6">
      <h2 class="text-lg font-semibold mb-4">设置</h2>

      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <label class="text-sm">字号</label>
          <input
            v-model.number="settings.settings.fontSize"
            type="number"
            min="10"
            max="48"
            class="w-20 px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm text-center"
          />
        </div>

        <div class="flex items-center justify-between">
          <label class="text-sm">下载源</label>
          <select
            v-model="settings.settings.downloadSource"
            class="px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm"
          >
            <option value="best">sekai.best</option>
            <option value="unipjsk">unipjsk.com</option>
            <option value="haruki">haruki (CN)</option>
          </select>
        </div>

        <div class="flex items-center justify-between">
          <label class="text-sm">保存 \\N 换行符</label>
          <input v-model="settings.settings.saveN" type="checkbox" class="accent-[var(--color-primary)]" />
        </div>

        <div class="flex items-center justify-between">
          <label class="text-sm">SSL 验证</label>
          <input v-model="settings.settings.disableSSL" type="checkbox" class="accent-[var(--color-primary)]" />
        </div>

        <div class="flex items-center justify-between">
          <label class="text-sm">外观模式</label>
          <select
            v-model="app.themeMode"
            class="px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm"
          >
            <option value="system">跟随系统</option>
            <option value="light">浅色</option>
            <option value="dark">深色</option>
          </select>
        </div>
      </div>

      <div class="flex justify-end gap-2 mt-6">
        <button
          class="px-4 py-1.5 rounded text-sm border border-[var(--color-border)] hover:bg-black/5 dark:hover:bg-white/10"
          @click="emit('close')"
        >
          取消
        </button>
        <button
          class="px-4 py-1.5 rounded text-sm text-white"
          style="background-color: var(--color-primary)"
          @click="settings.saveSettings(); emit('close')"
        >
          保存
        </button>
      </div>
    </div>
  </div>
</template>
