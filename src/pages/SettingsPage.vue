<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'
import { useAppStore } from '../stores/app'
import { useSettingsStore } from '../stores/settings'
import { useToast } from '../composables/useToast'
const router = useRouter()
const app = useAppStore()
const settings = useSettingsStore()
const toast = useToast()

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
    <header class="border-b border-[var(--color-border)] bg-[var(--color-surface)] px-6 py-3 flex items-center justify-between">
      <div class="flex items-center gap-4">
        <button
          @click="router.push('/')"
          class="flex items-center gap-1.5 text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-text)] transition-colors"
        >
          <ArrowLeft :size="18" />
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

    <main class="p-6 max-w-5xl mx-auto">

      <!-- ====== 编辑器 ====== -->
      <section class="mb-6">
        <h2 class="text-xs font-semibold text-[var(--color-text-secondary)] uppercase tracking-wider mb-3 px-1">编辑器</h2>
        <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-6">
          <div class="grid grid-cols-2 gap-6">
            <div class="flex items-center justify-between">
              <div>
                <div class="text-sm font-medium">字号</div>
                <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">编辑器文本显示大小</div>
              </div>
              <div class="flex items-center gap-2">
                <input v-model.number="settings.settings.fontSize" type="range" min="10" max="48" step="1" class="w-28 accent-[var(--color-primary)]" />
                <span class="text-sm w-8 text-center font-mono">{{ settings.settings.fontSize }}</span>
              </div>
            </div>

            <div class="flex items-center justify-between">
              <div>
                <div class="text-sm font-medium">索引排序</div>
                <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">故事索引下拉列表的显示顺序</div>
              </div>
              <select v-model="settings.settings.indexOrder" class="px-3 py-1.5 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm">
                <option value="desc">降序（最新的在底部）</option>
                <option value="asc">升序（最新的在顶部）</option>
              </select>
            </div>

            <label class="flex items-center justify-between cursor-pointer">
              <div>
                <div class="text-sm font-medium">切模式保留剧情</div>
                <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">切换翻/校/合时保留当前译文</div>
              </div>
              <input v-model="settings.settings.preserveStoryOnModeSwitch" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
            </label>
          </div>
        </div>
      </section>

      <!-- ====== 下载 ====== -->
      <section class="mb-6">
        <h2 class="text-xs font-semibold text-[var(--color-text-secondary)] uppercase tracking-wider mb-3 px-1">下载</h2>
        <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-6 max-w-lg">
          <h3 class="text-sm font-semibold mb-1">下载页默认目录</h3>
          <p class="text-xs text-[var(--color-text-secondary)] mb-3">专用下载页面 (/download) 的默认保存位置</p>
          <input
            v-model="settings.settings.jsonDownloadDir"
            type="text"
            placeholder="./downloads/json"
            class="w-full px-3 py-2 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm"
          />
        </div>
      </section>

      <!-- ====== 文件保存 ====== -->
      <section class="mb-6">
        <h2 class="text-xs font-semibold text-[var(--color-text-secondary)] uppercase tracking-wider mb-3 px-1">文件保存</h2>
        <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-6">
          <div class="grid grid-cols-2 gap-6">
            <label class="flex items-center justify-between cursor-pointer">
              <div>
                <div class="text-sm font-medium">保存 \\N 换行符</div>
                <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">翻译文件中保留 \\N 换行标记</div>
              </div>
              <input v-model="settings.settings.saveN" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
            </label>

            <label class="flex items-center justify-between cursor-pointer">
              <div>
                <div class="text-sm font-medium">保存语音文件</div>
                <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">下载并保存语音文件到本地</div>
              </div>
              <input v-model="settings.settings.saveVoice" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
            </label>

            <div class="col-span-2">
              <div class="text-sm font-medium mb-1">语音输出目录</div>
              <input
                v-model="settings.settings.voiceOutputDir"
                type="text"
                placeholder="留空使用默认目录"
                class="w-full px-3 py-2 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm"
              />
            </div>
          </div>
        </div>
      </section>

      <!-- ====== 外观 ====== -->
      <section class="mb-6">
        <h2 class="text-xs font-semibold text-[var(--color-text-secondary)] uppercase tracking-wider mb-3 px-1">外观</h2>
        <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-6">
          <div class="flex items-center justify-between max-w-sm">
            <div>
              <div class="text-sm font-medium">外观模式</div>
              <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">跟随系统或手动指定亮暗色</div>
            </div>
            <select
              v-model="app.themeMode"
              class="px-3 py-1.5 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm"
            >
              <option value="system">跟随系统</option>
              <option value="light">浅色</option>
              <option value="dark">深色</option>
            </select>
          </div>
        </div>
      </section>

      <!-- ====== 网络与调试 ====== -->
      <section class="mb-6">
        <h2 class="text-xs font-semibold text-[var(--color-text-secondary)] uppercase tracking-wider mb-3 px-1">网络与调试</h2>
        <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-6">
          <div class="grid grid-cols-2 gap-6">
            <div class="flex items-center justify-between">
              <div>
                <div class="text-sm font-medium">下载源</div>
                <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">故事 JSON 数据来源</div>
              </div>
              <span class="text-sm text-[var(--color-text-secondary)]">harukineo</span>
            </div>

            <label class="flex items-center justify-between cursor-pointer">
              <div>
                <div class="text-sm font-medium">SSL 验证</div>
                <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">禁用 SSL 证书验证（某些网络环境需要）</div>
              </div>
              <input v-model="settings.settings.disableSSL" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
            </label>

            <label class="flex items-center justify-between cursor-pointer">
              <div>
                <div class="text-sm font-medium">调试日志</div>
                <div class="text-xs text-[var(--color-text-secondary)] mt-0.5">在底部显示调试日志窗口</div>
              </div>
              <input v-model="settings.settings.debugEnabled" type="checkbox" class="accent-[var(--color-primary)] w-4 h-4" />
            </label>
          </div>
        </div>
      </section>

      <div class="flex justify-end gap-2 border-t border-[var(--color-border)] pt-4">
        <button @click="router.push('/')" class="px-4 py-1.5 rounded-lg text-sm border border-[var(--color-border)] hover:bg-black/5 dark:hover:bg-white/10 transition-colors">取消</button>
        <button @click="saveAndBack()" class="px-4 py-1.5 rounded-lg text-sm text-white transition-opacity hover:opacity-90" style="background-color: var(--color-primary)">保存设置</button>
      </div>
    </main>
  </div>
</template>
