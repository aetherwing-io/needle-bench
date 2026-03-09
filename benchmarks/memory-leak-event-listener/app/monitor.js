class Monitor {
  constructor(eventBus) {
    this.eventBus = eventBus;
    this.events = {
      validated: 0,
      errors: 0,
      tasks: 0,
    };

    // These listeners are registered once — no leak here
    this.eventBus.on('data:validated', () => { this.events.validated++; });
    this.eventBus.on('data:error', () => { this.events.errors++; });
    this.eventBus.on('task:queued', () => { this.events.tasks++; });
  }

  getStats() {
    const mem = process.memoryUsage();
    return {
      memory: {
        rss_mb: Math.round(mem.rss / 1024 / 1024 * 100) / 100,
        heapUsed_mb: Math.round(mem.heapUsed / 1024 / 1024 * 100) / 100,
        heapTotal_mb: Math.round(mem.heapTotal / 1024 / 1024 * 100) / 100,
        external_mb: Math.round(mem.external / 1024 / 1024 * 100) / 100,
      },
      events: { ...this.events },
      listenerCounts: {
        'data:validated': this.eventBus.listenerCount('data:validated'),
        'data:error': this.eventBus.listenerCount('data:error'),
        'task:queued': this.eventBus.listenerCount('task:queued'),
      },
      uptime: process.uptime(),
    };
  }
}

module.exports = { Monitor };
