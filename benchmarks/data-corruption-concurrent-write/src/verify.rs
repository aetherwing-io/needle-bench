use std::fs;
use std::io;

use crate::config::Config;

/// Verify that each segment in the file contains the expected fill pattern.
/// Returns Ok(true) if all segments are correct, Ok(false) if corruption found.
pub fn verify_file(path: &str) -> io::Result<bool> {
    let cfg = Config::default();
    let data = fs::read(path)?;

    if data.len() != cfg.total_size() {
        println!(
            "ERROR: File size {} != expected {}",
            data.len(),
            cfg.total_size()
        );
        return Ok(false);
    }

    let mut corrupt_segments = Vec::new();

    for seg in 0..cfg.num_segments {
        let offset = seg * cfg.segment_size;
        let end = offset + cfg.segment_size;
        let segment_data = &data[offset..end];

        let expected_fill = expected_fill_for_segment(seg);

        // Check every byte in the segment
        let mut mismatch_count = 0;
        let mut first_mismatch = None;
        for (i, &byte) in segment_data.iter().enumerate() {
            if byte != expected_fill {
                mismatch_count += 1;
                if first_mismatch.is_none() {
                    first_mismatch = Some((i, byte));
                }
            }
        }

        if mismatch_count > 0 {
            let (off, got) = first_mismatch.unwrap();
            println!(
                "CORRUPT: Segment {} — {} bytes wrong (first at offset {}: expected 0x{:02X}, got 0x{:02X})",
                seg, mismatch_count, off, expected_fill, got
            );
            corrupt_segments.push(seg);
        }
    }

    if corrupt_segments.is_empty() {
        Ok(true)
    } else {
        println!(
            "Found corruption in {} of {} segments: {:?}",
            corrupt_segments.len(),
            cfg.num_segments,
            corrupt_segments
        );
        Ok(false)
    }
}

/// Expected fill byte for a given segment.
fn expected_fill_for_segment(segment_id: usize) -> u8 {
    let base: u8 = 0xAA;
    base.wrapping_add(segment_id as u8)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_expected_fill() {
        assert_eq!(expected_fill_for_segment(0), 0xAA);
        assert_eq!(expected_fill_for_segment(1), 0xAB);
        assert_eq!(expected_fill_for_segment(5), 0xAF);
    }
}
