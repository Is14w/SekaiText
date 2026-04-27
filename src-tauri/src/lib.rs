#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
  tauri::Builder::default()
    .setup(|app| {
      if cfg!(debug_assertions) {
        app.handle().plugin(
          tauri_plugin_log::Builder::default()
            .level(log::LevelFilter::Info)
            .build(),
        )?;
      }
      Ok(())
    })
    .on_window_event(|_window, event| {
      if let tauri::WindowEvent::Destroyed = event {
        if cfg!(debug_assertions) {
          let _ = std::process::Command::new("node")
            .args(["../scripts/cleanup.mjs", "9800", "5173"])
            .spawn();
        }
      }
    })
    .run(tauri::generate_context!())
    .expect("error while running tauri application");
}
