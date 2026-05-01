<script setup lang="ts">
import { ref } from "vue";
import { useAppStore } from "../stores/app";
import { useEditorStore } from "../stores/editor";
import { useSettingsStore } from "../stores/settings";
import { api } from "../api/client";
import { Pencil, Check, CircleDot, ChevronLeft, ChevronRight, Cog } from "lucide-vue-next";
import StoryNavigator from "./navigation/StoryNavigator.vue";
import EditorWorkspace from "./editor/EditorWorkspace.vue";
import SpeakerCountDialog from "./dialogs/SpeakerCountDialog.vue";
import SpeakerCheckDialog from "./dialogs/SpeakerCheckDialog.vue";

const app = useAppStore();
const editor = useEditorStore();
const settings = useSettingsStore();

const showSpeakerCount = ref(false);
const showSpeakerCheck = ref(false);
const sidebarOpen = ref(true);

function setMode(key: number) {
  app.setEditorMode(key as 0 | 1 | 2);
}

const modes = [
  { key: 0, label: "翻译" },
  { key: 1, label: "校对" },
  { key: 2, label: "合意" },
];

const modeIcons: Record<number, typeof Pencil> = {
  0: Pencil,
  1: Check,
  2: CircleDot,
};

async function handleOpen() {
  const path = window.prompt("输入翻译文件路径：");
  if (!path) return;
  try {
    const result = await api.translationLoad(path);
    editor.setTalks(result.talks, result.talks, []);
    editor.currentFilePath = path;
    editor.markSaved();
  } catch (e: any) {
    alert("打开失败：" + e.message);
  }
}

async function handleSave() {
  if (editor.talks.length === 0) return;
  let path: string | null = editor.currentFilePath;
  if (!path) {
    path = window.prompt("输入保存路径：");
    if (!path) return;
  }
  try {
    await api.translationSave(path, editor.dstTalks, app.saveN);
    editor.currentFilePath = path;
    editor.markSaved();
  } catch (e: any) {
    alert("保存失败：" + e.message);
  }
}

function handleClear() {
  if (editor.hasUnsavedChanges) {
    if (!confirm("有未保存的更改，确定清空吗？")) return;
  }
  editor.clearAll();
}

async function handleFullCheck() {
  if (editor.talks.length === 0) return;
  let hasIssues = false;
  const msgs: string[] = [];
  for (const talk of editor.talks) {
    if (!talk.checked && talk.save) {
      hasIssues = true;
      msgs.push(`行 ${talk.idx}: ${talk.text.split("\n")[0]}`);
    }
  }
  if (hasIssues) {
    alert("发现以下问题：\n" + msgs.join("\n"));
  } else {
    alert("全文检查通过");
  }
}
</script>

<template>
  <div class="min-h-screen bg-[var(--color-bg)] flex">
    <!-- Left Sidebar -->
    <aside
      class="flex flex-col border-r border-[var(--color-border)] bg-[var(--color-surface)] flex-shrink-0 transition-all duration-200 overflow-hidden"
      :class="sidebarOpen ? 'w-44' : 'w-12'"
    >
      <!-- Toggle -->
      <button
        @click="sidebarOpen = !sidebarOpen"
        class="flex items-center gap-2 h-10 px-3 text-[var(--color-text-secondary)] hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors flex-shrink-0"
        :title="sidebarOpen ? '收起' : '展开'"
      >
        <ChevronLeft v-if="sidebarOpen" :size="18" />
        <ChevronRight v-else :size="18" />
        <span v-if="sidebarOpen" class="text-xs font-medium whitespace-nowrap"
          >导航</span
        >
      </button>

      <div class="border-b border-[var(--color-border)]" />

      <!-- Mode Buttons -->
      <div class="flex flex-col gap-0.5 p-1.5">
        <button
          v-for="m in modes"
          :key="m.key"
          @click="setMode(m.key)"
          class="flex items-center gap-2.5 h-9 px-2 rounded-lg transition-colors text-sm flex-shrink-0 group relative"
          :class="
            app.editorMode === m.key
              ? 'bg-[var(--color-primary)]/10 text-[var(--color-primary)] font-medium'
              : 'text-[var(--color-text-secondary)] hover:bg-gray-100 dark:hover:bg-gray-800'
          "
        >
          <component :is="modeIcons[m.key]" :size="18" />
          <span v-if="sidebarOpen" class="whitespace-nowrap">{{
            m.label
          }}</span>
          <span
            v-if="!sidebarOpen"
            class="absolute left-full ml-2 px-2 py-0.5 rounded text-xs bg-gray-900 text-white whitespace-nowrap opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity z-50"
          >
            {{ m.label }}
          </span>
        </button>
      </div>

      <div class="border-t border-[var(--color-border)]" />

      <!-- Settings Section -->
      <div
        v-if="sidebarOpen"
        class="flex-1 overflow-y-auto p-3 space-y-3 text-sm"
      >
        <div
          class="text-xs font-medium text-[var(--color-text-secondary)] uppercase tracking-wider"
        >
          设置
        </div>

        <div class="flex items-center justify-between">
          <label class="text-[var(--color-text-secondary)]">字号</label>
          <input
            v-model.number="settings.settings.fontSize"
            type="number"
            min="10"
            max="48"
            class="w-16 px-1.5 py-0.5 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-xs text-center"
          />
        </div>

        <label class="flex items-center justify-between cursor-pointer">
          <span class="text-[var(--color-text-secondary)]">保存 \\N</span>
          <input
            v-model="settings.settings.saveN"
            type="checkbox"
            class="accent-[var(--color-primary)] w-3.5 h-3.5"
          />
        </label>

        <div class="flex items-center justify-between gap-2">
          <span class="text-[var(--color-text-secondary)]">外观</span>
          <select
            v-model="app.themeMode"
            class="w-20 px-1.5 py-0.5 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-xs"
          >
            <option value="system">系统</option>
            <option value="light">浅色</option>
            <option value="dark">深色</option>
          </select>
        </div>

        <label class="flex items-center justify-between cursor-pointer">
          <span class="text-[var(--color-text-secondary)]">SSL 验证</span>
          <input
            v-model="settings.settings.disableSSL"
            type="checkbox"
            class="accent-[var(--color-primary)] w-3.5 h-3.5"
          />
        </label>

        <label class="flex items-center justify-between cursor-pointer">
          <span class="text-[var(--color-text-secondary)]">切模式保留剧情</span>
          <input
            v-model="settings.settings.preserveStoryOnModeSwitch"
            type="checkbox"
            class="accent-[var(--color-primary)] w-3.5 h-3.5"
          />
        </label>
      </div>

      <!-- Collapsed: Settings button -->
      <div v-if="!sidebarOpen" class="flex-1" />

      <!-- Bottom section -->
      <div class="border-t border-[var(--color-border)] p-1.5">
        <button
          v-if="!sidebarOpen"
          @click="settings.saveSettings()"
          class="flex items-center gap-2.5 h-9 w-full px-2 rounded-lg transition-colors text-sm text-[var(--color-text-secondary)] hover:bg-gray-100 dark:hover:bg-gray-800 group relative"
          title="保存设置"
        >
          <Cog :size="18" />
          <span
            class="absolute left-full ml-2 px-2 py-0.5 rounded text-xs bg-gray-900 text-white whitespace-nowrap opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity z-50"
          >
            设置
          </span>
        </button>
        <button
          v-if="sidebarOpen"
          @click="settings.saveSettings()"
          class="w-full px-2 py-1.5 rounded-lg text-xs text-white transition-opacity hover:opacity-90"
          style="background-color: var(--color-primary)"
        >
          保存设置
        </button>
      </div>
    </aside>

    <!-- Main Area -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- Top Bar -->
      <header
        class="border-b border-[var(--color-border)] bg-[var(--color-surface)] px-4 py-2"
      >
        <StoryNavigator />
      </header>

      <!-- Toolbar -->
      <div
        class="border-b border-[var(--color-border)] bg-[var(--color-surface)] px-4 py-1.5"
      >
        <div class="flex items-center gap-2 flex-wrap text-sm">
          <button
            @click="handleOpen"
            class="px-2.5 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-800 text-[var(--color-text-secondary)] hover:text-[var(--color-text)] transition-colors"
          >
            打开
          </button>
          <button
            @click="handleSave"
            class="px-2.5 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-800 text-[var(--color-text-secondary)] hover:text-[var(--color-text)] transition-colors"
          >
            保存
          </button>
          <button
            @click="handleClear"
            class="px-2.5 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-800 text-[var(--color-text-secondary)] hover:text-[var(--color-text)] transition-colors"
          >
            清空
          </button>

          <div class="w-px h-4 bg-[var(--color-border)]" />

          <label
            class="flex items-center gap-1 cursor-pointer text-[var(--color-text-secondary)] hover:text-[var(--color-text)]"
          >
            <input
              v-model="app.showFlashback"
              type="checkbox"
              class="accent-[var(--color-primary)] w-3 h-3"
            />
            闪回
          </label>
          <label
            class="flex items-center gap-1 cursor-pointer text-[var(--color-text-secondary)] hover:text-[var(--color-text)]"
          >
            <input
              v-model="app.syncScroll"
              type="checkbox"
              class="accent-[var(--color-primary)] w-3 h-3"
            />
            同步
          </label>
          <label
            class="flex items-center gap-1 cursor-pointer text-[var(--color-text-secondary)] hover:text-[var(--color-text)]"
          >
            <input
              v-model="app.showDiff"
              type="checkbox"
              class="accent-[var(--color-primary)] w-3 h-3"
            />
            差异
          </label>

          <div class="w-px h-4 bg-[var(--color-border)]" />

          <button
            @click="showSpeakerCheck = true"
            class="px-2.5 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-800 text-[var(--color-text-secondary)] hover:text-[var(--color-text)] transition-colors"
          >
            说话人
          </button>
          <button
            @click="handleFullCheck"
            class="px-2.5 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-800 text-[var(--color-text-secondary)] hover:text-[var(--color-text)] transition-colors"
          >
            检查
          </button>
          <button
            @click="showSpeakerCount = true"
            class="px-2.5 py-1 rounded hover:bg-gray-100 dark:hover:bg-gray-800 text-[var(--color-text-secondary)] hover:text-[var(--color-text)] transition-colors"
          >
            统计
          </button>
        </div>
      </div>

      <!-- Editor -->
      <main class="flex-1 p-4 min-h-0">
        <EditorWorkspace />
      </main>
    </div>

    <!-- Dialogs -->
    <SpeakerCountDialog
      v-if="showSpeakerCount"
      @close="showSpeakerCount = false"
    />
    <SpeakerCheckDialog
      v-if="showSpeakerCheck"
      @close="showSpeakerCheck = false"
    />
  </div>
</template>
