use std::fs::{File, OpenOptions};
use std::io::{self, Seek, SeekFrom, Write};
use std::sync::Arc;
use std::thread;

use crate::config::Config;

/// Fill a byte pattern for a given segment.
/// Each segment has a deterministic fill value based on segment index.
fn segment_fill(segment_id: usize) -> u8 {
    // Deterministic: segment 0 = 0xAA, segment 1 = 0xBB, etc.
    let base: u8 = 0xAA;
    base.wrapping_add(segment_id as u8)
}

/// Write segments to the file. Each worker opens its own file descriptor
/// and seeks to the appropriate offset before writing.
fn write_segments(
    path: &str,
    worker_id: usize,
    segments: Vec<usize>,
    segment_size: usize,
    passes: usize,
) -> io::Result<()> {
    for _pass in 0..passes {
        // Open file each pass — simulates real-world pattern where workers
        // reopen files (e.g., log rotation, batch processing)
        let mut file = OpenOptions::new().write(true).open(path)?;

        for &seg in &segments {
            let offset = (seg * segment_size) as u64;
            let fill = segment_fill(seg);
            let data = vec![fill; segment_size];

            file.seek(SeekFrom::Start(offset))?;
            file.write_all(&data)?;
            // NOTE: no fsync between segments — we rely on OS buffering
        }
    }

    eprintln!("Worker {} completed: segments {:?}", worker_id, segments);
    Ok(())
}

/// Perform a concurrent write to the output file using multiple threads.
/// Each thread is assigned a range of segments to write.
pub fn concurrent_write(output_path: &str, cfg: &Config) -> io::Result<()> {
    let total_size = cfg.total_size();

    // Pre-allocate file to full size
    {
        let file = File::create(output_path)?;
        file.set_len(total_size as u64)?;
    }

    eprintln!(
        "Writing {} bytes with {} workers ({} segments of {} bytes)",
        total_size, cfg.num_workers, cfg.num_segments, cfg.segment_size
    );

    let path = Arc::new(output_path.to_string());
    let mut handles = Vec::new();

    for worker_id in 0..cfg.num_workers {
        let path = Arc::clone(&path);
        let segments = cfg.segments_for_worker(worker_id);
        let segment_size = cfg.segment_size;
        let passes = cfg.passes_per_worker;

        let handle = thread::spawn(move || {
            write_segments(&path, worker_id, segments, segment_size, passes)
        });
        handles.push(handle);
    }

    // Wait for all workers
    for (i, handle) in handles.into_iter().enumerate() {
        match handle.join() {
            Ok(Ok(())) => {}
            Ok(Err(e)) => {
                return Err(io::Error::new(
                    io::ErrorKind::Other,
                    format!("Worker {} failed: {}", i, e),
                ));
            }
            Err(_) => {
                return Err(io::Error::new(
                    io::ErrorKind::Other,
                    format!("Worker {} panicked", i),
                ));
            }
        }
    }

    eprintln!("All workers completed");
    Ok(())
}
