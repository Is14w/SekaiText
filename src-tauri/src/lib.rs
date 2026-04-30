use tauri::Manager;
use std::os::windows::process::CommandExt;
use std::process::{Child, Command as StdCommand};
use std::sync::Mutex;

const CREATE_NO_WINDOW: u32 = 0x08000000;

struct SidecarProcess(Mutex<Option<Child>>);

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
  tauri::Builder::default()
    .plugin(tauri_plugin_dialog::init())
    .setup(|app| {
      if cfg!(debug_assertions) {
        app.handle().plugin(
          tauri_plugin_log::Builder::default()
            .level(log::LevelFilter::Info)
            .build(),
        )?;
      }

      // Spawn Go backend in release mode (dev uses external server)
      if !cfg!(debug_assertions) {
        let resource_dir = app.path().resource_dir()
          .expect("failed to resolve resource directory");

        // Sidecar binary is next to app.exe in the install dir
        let exe_dir = std::env::current_exe()
          .expect("failed to get exe path")
          .parent()
          .expect("failed to get exe dir")
          .to_path_buf();

        let sidecar_path = exe_dir.join("sekaitext-backend.exe");

        let child = StdCommand::new(&sidecar_path)
          .args(["--port", "9800", "--dir", &resource_dir.to_string_lossy()])
          .creation_flags(CREATE_NO_WINDOW)
          .spawn()
          .unwrap_or_else(|e| panic!("failed to spawn {}: {}", sidecar_path.display(), e));

        app.manage(SidecarProcess(Mutex::new(Some(child))));
      }

      Ok(())
    })
    .on_window_event(|window, event| {
      if let tauri::WindowEvent::Destroyed = event {
        // Kill Go sidecar in release mode
        if !cfg!(debug_assertions) {
          if let Some(state) = window.app_handle().try_state::<SidecarProcess>() {
            if let Some(mut child) = state.0.lock().unwrap().take() {
              let _ = child.kill();
            }
          }
        }
        // Dev cleanup: kill the separate Go+Vite processes
        if cfg!(debug_assertions) {
          let _ = StdCommand::new("node")
            .args(["../scripts/cleanup.mjs", "9800", "5173"])
            .spawn();
        }
      }
    })
    .run(tauri::generate_context!())
    .expect("error while running tauri application");
}
