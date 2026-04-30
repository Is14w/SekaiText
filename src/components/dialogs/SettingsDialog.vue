<script setup lang="ts">
import { reactive } from 'vue'
import { useAppStore } from '../../stores/app'
import { useSettingsStore } from '../../stores/settings'
import { useToast } from '../../composables/useToast'
import { useDebugLog } from '../../composables/useDebugLog'
import type { Settings } from '../../types/api'

const app = useAppStore()
const settings = useSettingsStore()
const toast = useToast()
const debug = useDebugLog()

const emit = defineEmits<{
  close: []
}>()

const local = reactive<Settings>(JSON.parse(JSON.stringify(settings.settings)))

async function save() {
  Object.assign(settings.settings, local)
  debug.enabled.value = local.debugEnabled
  try {
    await settings.saveSettings()
    toast.show('设置已保存', 'success')
    emit('close')
  } catch {
    toast.show('保存失败', 'error')
  }
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" @click.self="emit('close')">
    <div class="bg-[var(--color-surface)] rounded-xl shadow-xl border border-[var(--color-border)] w-full max-w-md max-h-[85vh] overflow-y-auto p-6">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-lg font-semibold">设置</h2>
        <button
          @click="emit('close')"
          class="w-8 h-8 flex items-center justify-center rounded-lg text-[var(--color-text-secondary)] hover:bg-black/5 dark:hover:bg-white/10 transition-colors text-lg leading-none"
        >✕</button>
      </div>

      <div class="space-y-5">
        <!-- Font Size -->
        <div class="flex items-center justify-between">
          <div>
            <div class="text-sm font-medium">字号</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">编辑器文本显示大小</div>
          </div>
          <div class="flex items-center gap-2">
            <input
              v-model.number="local.fontSize"
              type="range" min="10" max="48" step="1"
              class="w-28 accent-[var(--color-primary)]"
            />
            <span class="text-sm w-8 text-center font-mono">{{ local.fontSize }}</span>
          </div>
        </div>

        <div class="border-t border-[var(--color-border)]" />

        <!-- Download Source -->
        <div class="flex items-center justify-between">
          <div>
            <div class="text-sm font-medium">下载源</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">故事 JSON 数据来源</div>
          </div>
          <span class="text-sm text-[var(--color-text-secondary)]">harukineo</span>
        </div>

        <div class="border-t border-[var(--color-border)]" />

        <!-- Save \N -->
        <label class="flex items-center justify-between cursor-pointer">
          <div>
            <div class="text-sm font-medium">保存 \\N 换行符</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">翻译文件中保留 \\N 换行标记</div>
          </div>
          <input v-model="local.saveN" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
        </label>

        <div class="border-t border-[var(--color-border)]" />

        <!-- Save Voice -->
        <label class="flex items-center justify-between cursor-pointer">
          <div>
            <div class="text-sm font-medium">保存语音文件</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">下载并保存语音文件到本地</div>
          </div>
          <input v-model="local.saveVoice" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
        </label>

        <div class="border-t border-[var(--color-border)]" />

        <!-- SSL Verification -->
        <label class="flex items-center justify-between cursor-pointer">
          <div>
            <div class="text-sm font-medium">SSL 验证</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">禁用 SSL 证书验证（某些网络环境需要）</div>
          </div>
          <input v-model="local.disableSSL" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
        </label>

        <div class="border-t border-[var(--color-border)]" />

        <!-- Index Order -->
        <div class="flex items-center justify-between">
          <div>
            <div class="text-sm font-medium">索引排序</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">故事索引下拉列表的显示顺序</div>
          </div>
          <select v-model="local.indexOrder" class="px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-surface)] text-sm">
            <option value="desc">降序（最新的在底部）</option>
            <option value="asc">升序（最新的在顶部）</option>
          </select>
        </div>

        <div class="border-t border-[var(--color-border)]" />

        <!-- Appearance -->
        <div class="flex items-center justify-between">
          <div>
            <div class="text-sm font-medium">外观模式</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">跟随系统或手动指定亮暗色</div>
          </div>
          <select
            v-model="app.themeMode"
            class="px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-surface)] text-sm"
          >
            <option value="system">跟随系统</option>
            <option value="light">浅色</option>
            <option value="dark">深色</option>
          </select>
        </div>

        <div class="border-t border-[var(--color-border)]" />

        <!-- Debug -->
        <label class="flex items-center justify-between cursor-pointer">
          <div>
            <div class="text-sm font-medium">调试日志</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">在底部显示调试日志窗口</div>
          </div>
          <input v-model="local.debugEnabled" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
        </label>
      </div>

      <div class="mt-6 flex justify-end gap-2 border-t border-[var(--color-border)] pt-4">
        <button
          @click="emit('close')"
          class="px-4 py-1.5 rounded-lg text-sm border border-[var(--color-border)] hover:bg-black/5 dark:hover:bg-white/10 transition-colors"
        >
          取消
        </button>
        <button
          @click="save()"
          class="px-4 py-1.5 rounded-lg text-sm text-white transition-opacity hover:opacity-90"
          style="background-color: var(--color-primary)"
        >
          保存设置
        </button>
      </div>
    </div>
  </div>
</template>
