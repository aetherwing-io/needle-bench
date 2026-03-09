use std::env;
use std::process;

mod writer;
mod verify;
mod config;

fn main() {
    let args: Vec<String> = env::args().collect();

    match args.get(1).map(|s| s.as_str()) {
        Some("write") => {
            let output = args.get(2).unwrap_or_else(|| {
                eprintln!("Usage: concurrent-writer write <output-file>");
                process::exit(1);
            });
            let cfg = config::Config::default();
            if let Err(e) = writer::concurrent_write(output, &cfg) {
                eprintln!("Error: {}", e);
                process::exit(1);
            }
        }
        Some("verify") => {
            let path = args.get(2).unwrap_or_else(|| {
                eprintln!("Usage: concurrent-writer verify <file>");
                process::exit(1);
            });
            match verify::verify_file(path) {
                Ok(true) => {
                    println!("VERIFY: File integrity check passed");
                }
                Ok(false) => {
                    println!("VERIFY: File integrity check FAILED — data corruption detected");
                    process::exit(1);
                }
                Err(e) => {
                    eprintln!("Error verifying file: {}", e);
                    process::exit(1);
                }
            }
        }
        _ => {
            eprintln!("Usage: concurrent-writer <write|verify> <file>");
            eprintln!();
            eprintln!("  write   - Write data to file using concurrent workers");
            eprintln!("  verify  - Verify file integrity");
            process::exit(1);
        }
    }
}
