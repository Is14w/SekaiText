<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAppStore } from '../stores/app'
import { useSettingsStore } from '../stores/settings'
import { useToast } from '../composables/useToast'
import { useDebugLog } from '../composables/useDebugLog'

const router = useRouter()
const app = useAppStore()
const settings = useSettingsStore()
const toast = useToast()
const debug = useDebugLog()

function saveAndBack() {
  settings.saveSettings().then(() => {
    toast.show('设置已保存', 'success')
    router.push('/')
  }).catch(() => {
    toast.show('保存失败', 'error')
  })
}
</script>

<template>
  <div class="min-h-screen bg-[var(--color-bg)]">
    <!-- Header -->
    <header data-tauri-drag-region class="border-b border-[var(--color-border)] bg-[var(--color-surface)] px-6 py-3 flex items-center justify-between">
      <div class="flex items-center gap-4">
        <button
          @click="router.push('/')"
          class="flex items-center gap-1.5 text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-text)] transition-colors"
        >
          <span class="text-base leading-none">←</span>
          返回编辑器
        </button>
        <span class="text-sm font-medium">设置</span>
      </div>
      <button
        @click="saveAndBack()"
        class="px-4 py-1.5 rounded-lg text-sm text-white transition-opacity hover:opacity-90"
        style="background-color: var(--color-primary)"
      >
        保存并返回
      </button>
    </header>

    <!-- Content -->
    <main class="max-w-lg mx-auto p-6">
      <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-6 space-y-5">
        <!-- Font Size -->
        <div class="flex items-center justify-between">
          <div>
            <div class="text-sm font-medium">字号</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">编辑器文本显示大小</div>
          </div>
          <div class="flex items-center gap-2">
            <input
              v-model.number="settings.settings.fontSize"
              type="range" min="10" max="48" step="1"
              class="w-28 accent-[var(--color-primary)]"
            />
            <span class="text-sm w-8 text-center font-mono">{{ settings.settings.fontSize }}</span>
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
          <input v-model="settings.settings.saveN" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
        </label>

        <div class="border-t border-[var(--color-border)]" />

        <!-- Save Voice -->
        <label class="flex items-center justify-between cursor-pointer">
          <div>
            <div class="text-sm font-medium">保存语音文件</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">下载并保存语音文件到本地</div>
          </div>
          <input v-model="settings.settings.saveVoice" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
        </label>

        <div class="border-t border-[var(--color-border)]" />

        <!-- SSL Verification -->
        <label class="flex items-center justify-between cursor-pointer">
          <div>
            <div class="text-sm font-medium">SSL 验证</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">禁用 SSL 证书验证（某些网络环境需要）</div>
          </div>
          <input v-model="settings.settings.disableSSL" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
        </label>

        <div class="border-t border-[var(--color-border)]" />

        <!-- Index Order -->
        <div class="flex items-center justify-between">
          <div>
            <div class="text-sm font-medium">索引排序</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">故事索引下拉列表的显示顺序</div>
          </div>
          <select v-model="settings.settings.indexOrder" class="px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-surface)] text-sm">
            <option value="desc">降序（最新的在底部）</option>
            <option value="asc">升序（最新的在顶部）</option>
          </select>
        </div>

        <div class="border-t border-[var(--color-border)]" />

        <!-- Dark Mode -->
        <label class="flex items-center justify-between cursor-pointer">
          <div>
            <div class="text-sm font-medium">暗色主题</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">切换亮色/暗色显示</div>
          </div>
          <input
            :checked="app.isDark"
            type="checkbox"
            class="accent-[var(--color-primary)] w-4 h-4"
            @change="app.toggleTheme()"
          />
        </label>

        <div class="border-t border-[var(--color-border)]" />

        <!-- Debug -->
        <label class="flex items-center justify-between cursor-pointer">
          <div>
            <div class="text-sm font-medium">调试日志</div>
            <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">在底部显示调试日志窗口</div>
          </div>
          <input v-model="debug.enabled" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
        </label>
      </div>

      <div class="mt-4 flex justify-end gap-2">
        <button
          @click="router.push('/')"
          class="px-4 py-1.5 rounded-lg text-sm border border-[var(--color-border)] hover:bg-black/5 dark:hover:bg-white/10 transition-colors"
        >
          取消
        </button>
        <button
          @click="saveAndBack()"
          class="px-4 py-1.5 rounded-lg text-sm text-white transition-opacity hover:opacity-90"
          style="background-color: var(--color-primary)"
        >
          保存设置
        </button>
      </div>
    </main>
  </div>
</template>
