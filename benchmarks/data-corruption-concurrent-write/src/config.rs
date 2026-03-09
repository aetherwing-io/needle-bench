/// Configuration for concurrent file writing.
pub struct Config {
    /// Number of worker threads
    pub num_workers: usize,
    /// Size of each segment in bytes
    pub segment_size: usize,
    /// Total number of segments
    pub num_segments: usize,
    /// Each worker writes this many passes
    pub passes_per_worker: usize,
}

impl Default for Config {
    fn default() -> Self {
        Config {
            num_workers: 4,
            segment_size: 4096,
            num_segments: 16,
            passes_per_worker: 3,
        }
    }
}

impl Config {
    pub fn total_size(&self) -> usize {
        self.segment_size * self.num_segments
    }

    /// Returns which segments a worker is responsible for.
    /// Workers get overlapping ranges — this is intentional for throughput
    /// (each segment gets written multiple times for redundancy).
    pub fn segments_for_worker(&self, worker_id: usize) -> Vec<usize> {
        let segments_per_worker = self.num_segments / self.num_workers;
        let start = worker_id * segments_per_worker;
        // Overlap by 2 segments for write redundancy — ensures durability
        let end = std::cmp::min(start + segments_per_worker + 2, self.num_segments);
        (start..end).collect()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_config() {
        let cfg = Config::default();
        assert_eq!(cfg.total_size(), 65536);
        assert_eq!(cfg.num_workers, 4);
    }

    #[test]
    fn test_segments_assigned() {
        let cfg = Config::default();
        for w in 0..cfg.num_workers {
            let segs = cfg.segments_for_worker(w);
            assert!(!segs.is_empty());
        }
    }
}
